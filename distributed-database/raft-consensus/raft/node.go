package raft

type Node interface {
	Init()
	Recover()
	TriggerElection()
	Vote(request VoteRequest)
	ReceiveVote(response VoteResponse)
	Broadcast(msg string)
	SendPeriodicHeartBeat()
	ReplicateLog(leader, follower string)
	HandleLogRequest(request LogRequest)
	AppendEntries(LogLength int, commitLength int, entries []log)
	HandleLogResponse(response LogResponse)
	CommitLogEntries()
}

type LogRequest struct {
	leaderId string
	currentTerm int
	logLength int
	prevLogTerm int
	commitLength int
	entries []Log
}

type LogResponse struct {
	nodeId int
	currentTerm int 
	ack int
	result bool
}

type Log struct {
	Message interface {}
	Term int
}	

type VoteRequest struct {
	nodeId string
	currentTerm int
	logLength int
	lastTerm int
}

type VoteResponse struct {
	nodeId string
	currentTerm int
	result bool
}

type Node struct {
	currentTerm int
	votedFor map[string]bool
	log []log
	commitLength int
	currentRole string
	currentLeader interface {}
	votesReceived []string
	sentLength map[string]int
	ackedLength map[string]int
	nodeId string
}

func NewNode(nodeId string) Node {
	n := node {
		currentRole: "follower",
		nodeId: nodeId,
	}
	go n.SendPeriodicHeartBeat()
	return n
}

func (n node) Init() {
	n.currentTerm = 0
	n.votedFor = nil
	n.Log = []Log{}
	n.currentRole = "follower"
	n.currentLeader = nil
	n.votesReceived = []string{}
	n.sentLength = nil
	n.ackedLength = nil
}

func (n node) Recover() {
	n.currentRole = "follower"
	n.currentLeader = nil
	n.votesReceived = []string {}
	n.sentLength = nil
	n.ackedLength = nil
}
// TriggerElection is for triggering election if:
// 1. There is a problem with leader or
// 2. Current election is time out
func (n node) TriggerElection()  {
	nodes := []string{""} // nodes should be maintain by implimentation owner.	
	n.currentTerm += n.currentTerm + 1
	n.currentRole = "candidate"
	n.votedFor[n.nodeId] = true
	n.votesReceived = []string{n.nodeId}

	lastTerm := 0
	if len(n.log) > 0 {
		lastTerm = n.log[len(n.log) - 1].Term
	}
	message := VoteRequest {
		nodeId: n.nodeId,
		currentTerm: n.currentTerm,
		logLength: len(n.log),
		lastTerm: lastTerm,
	}

	for _, node := range nodes {
		sendVoteRequest(message, node)
	}
	startTimer()
}

func (n node) Vote(request VoteRequest)  {
	currentTerm := n.Log[len(n.log)-1].Term

	isLogOk := request.lastTerm > currentTerm || 
		(request.lastTerm == currentTerm && request.logLength >= len(n.log))
	
	isTermOk := request.currentTerm > n.currentTerm ||
		(request.lastTerm == currentTerm && n.votedFor[request.nodeId] == false)
	
	if isLogOk && isTermOk {
		n.currentTerm = request.lastTerm
		n.currentRole = "follower"
		n.votedFor[request.nodeId] = true
		sendVoteResponse(VoteResponse {
			nodeId: n.nodeId,
			currentTerm: n.currentTerm,
			result: true,
		}, request.nodeId)
	} else {
		sendVoteResponse(VoteResponse {
			nodeId: n.nodeId,
			currentTerm: n.currentTerm,
			result: false,
		}, request.nodeId)
	}
}
// ReceiveVote is for handling a vote result for specific node.
func (n node) ReceiveVote(response VoteResponse) {
	nodes := []string{""}
	if n.currentRole == "candidate" && response.currentTerm == n.currentTerm && response.result {
		n.votesReceived = append(n.votesReceived, response.nodeId)
		if len(n.votesReceived) >= (len(nodes) + 1) / 2 {
			n.currentRole = "leader"
			n.currentLeader = n.nodeId
			cancelTimer()
			for _, follower := range nodes {
				if follower == n.nodeId {
					continue
				}
				n.sentLength[follower] = len(n.log)
				n.ackedLength[follower] = 0

				n.ReplicateLog(n.nodeId, follower)
			}
		}
	} else if response.currentTerm > n.currentTerm {
		n.currentTerm = response.currentTerm
		n.currentRole = "follower"
		n.votedFor = nil
		cancelTimer()
	}
}

// Broadcast is for sending message to followers.
// for example, client puts data into map, leader broadcast binary data to followers.
func (n node) Broadcast(msg string) {
	nodes := []string{""}
	if n.currentRole == "leader" {
		n.log = append(n.log, Log {
			Message: msg,
			Term: n.currentTerm,
		})
		n.ackedLength[n.nodeId] = len(n.log)
		for _, follower := range nodes {
			if follower == n.nodeId {
				continue
			}
			n.ReplicateLog(n.nodeId, follower)
		}
	} else {
		forward(msg, n.currentLeader)
	}
}

// SendPeriodicHeartBeat is for periodically send heartbeat to followers to say,
// "Hey, I'm the leader and I'm alive" 
func (n node) SendPeriodicHeartBeat() {
	nodes := []string{""}
	if n.currentRole == "leader" {
		for _, follower := range nodes {
			if follower == n.nodeId {
				continue
			}
			n.ReplicateLog(n.nodeId, follower)
		}
	}
}

fun (n node) ReplicateLog(leader, follower string) {
	i := n.sentLength[follower]
	entriesToBeSent := n.log[i:]
	prevLogTerm := 0
	if i > 0 {
		prevLogTerm = n.log[i-1].Term
	}
	replicate(LogRequest {
		leaderId: leader,
		currentTerm: n.currentTerm,
		logLength: i,
		prevLogTerm: prevLogTerm,
		commitLength, n.commitLength,
		entries: entriesToBeSent,
	}, follower)
}

func (n node) HandleLogRequest(request LogRequest)  {
	if request.currentTerm > n.currentTerm {
		n.currentTerm = request.currentTerm
		n.votedFor = nil
	}

	isLogOk := len(n.log) >= request.logLength

	if isLogOk && request.logLength > 0 {
		isLogOk = request.prevLogTerm = n.log[request.logLength - 1].Term
	}
}

func (n node) HandleLogResponse(response LogResponse)  {

}

func (n node) CommitLogEntries()  {
	
}

func (n node) AppendEntries(logLength int, leaderCommitLength int, entries []log)  {
	
}

// deliver is for persisting log into real application.
func deliver(log log)  {}

func forward(msg string, leader interface{})  {}
func sendVoteRequest(msg VoteRequest, node string)  {}
func sendVoteResponse(msg VoteResponse, node string)  {
	
}
func sendLogResponse(response VoteResponse, nodeId string)  {
	
}
func startTimer()  {
	
}
func cancelTimer()  {
	
}
func replicate(req LogRequest, follower string)  {
	
}

func ready(logs []Log, length int) []Log {
	return []Log{} // find ready log to be commited.
}