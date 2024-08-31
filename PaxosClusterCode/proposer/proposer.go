package proposer

type ProposerRole struct {
    roleId uint64
    log *replicatedlog.Log
    peers *clusterpeers.Cluster
    proposals *manager.ProposalManager
    client chan ClientRequest
    heartbeat chan uint64
    terminator chan bool
}

func Construct(roleId uint64, proposals *manager.ProposalManager, log *replicatedlog.Log, 
	peers *clusterpeers.Cluster) *ProposerRole {
		newProposerRole := ProposerRole {
			roleId: roleId,
			log: log,
			peers: peers,
			proposals: proposals,    
			client: make(chan ClientRequest),
			heartbeat: make(chan uint64),
			terminator: make(chan bool),
		}
		return &newProposerRole
	}

// Starts proposer role state machine
func Run(this *ProposerRole) {
    isLeaderStateChannel := make(chan bool)
    isNotLeaderStateChannel := make(chan bool)
    go this.isNotLeaderState(isLeaderStateChannel, isNotLeaderStateChannel)
    go this.isLeaderState(isNotLeaderStateChannel, isLeaderStateChannel)
}

// Role is not leader; will reject client requests
// <-chan readonly, chan<- write only
func (this *ProposerRole) isNotLeaderState(trans chan<- bool, self <-chan bool) {
    electionNotify := make(chan bool)
    startElection := make(chan bool)
    go this.electLeader(startElection, electionNotify)

    for {
        select {
        case <- electionNotify:
            trans <- true
            <- self
            startElection <- true
        case request := <- this.client:
            request.reply <- fmt.Errorf("[ PROPOSER %d ] Failure: not cluster leader", this.roleId)
        case <- this.terminator:
            return
        }
    }
}

// Elects self leader if not receiving heartbeat signal from role with higher ID
func (this *ProposerRole) electLeader(startElection <-chan bool, electionNotify chan<- bool) {
    for {
        select {
        case <- this.heartbeat:
            continue
        case <- time.After(2*time.Second):
            electionNotify <- true
            <- startElection
        }
    }
}

// Role is leader; will furnish client requests
// trans = isNotLeaderStateChannel; self = isLeaderStateChannel
func (this *ProposerRole) isLeaderState(trans chan<- bool, self <-chan bool) {
    <- self

    for {
        select {
        case <- this.heartbeat:
            trans <- true
            <- self
        case request := <- this.client:
            fmt.Println("[ PROPOSER", this.roleId, "] Initiating paxos for client request", request.value)
            go func () { request.reply <- this.paxos(request.value) }()
        case <- this.terminator:
            return
        }
    }
}

// Executes single round of Paxos protocol
func (this *ProposerRole) paxos(value string) error {

	for !chosen {
        index := this.log.GetFirstUnchosenIndex()
		proposalId := this.proposals.GetCurrentProposalId()
		usingValue := value

        // Prepare phase
        fmt.Println("[ PROPOSER", roleId, "] Executing prepare phase of protocol for value", usingValue)
        request := acceptor.PrepareReq {
            ProposalId: proposalId, 
            Index: index,
        }
		peerCount, endpoint := this.peers.BroadcastPrepareRequest(request)
		success, changed, changedValue, err := this.recvPromises(peerCount, endpoint)

		if success {
            if changed {
                usingValue = changedValue
                fmt.Println("[ PROPOSER", roleId, "] Value", value, "superseded by value", changedValue)
            }

            // Proposal phase
            fmt.Println("[ PROPOSER", roleId, "] Executing proposal phase of protocol for value", usingValue)
            request := acceptor.ProposalReq {
                ProposalId: proposalId, 
                Index: index, 
                Value: usingValue, 
                FirstUnchosenIndex: this.log.GetFirstUnchosenIndex(),
            }
            peerCount, endpoint := this.peers.BroadcastProposalRequest(request, nil)
            success, err = this.recvAccepts(request, peerCount, endpoint)
            if err != nil { return err }

            if success {
                fmt.Println("[ PROPOSER", roleId, "] Success; chose", usingValue, "for log entry", index)
                this.log.SetEntryAt(index, usingValue, proposal.Chosen())
                chosen = !changed
            } else {
                this.proposals.GenerateNextProposalId()
            }
        } else {
            fmt.Println("[ PROPOSER", roleId, "] Retrying prepare phase for", usingValue)
            this.proposals.GenerateNextProposalId()
        }
	}
}

// Receves replies to prepare requests
func (this *ProposerRole) recvPromises(peerCount uint64, endpoint <-chan clusterpeers.Response) (bool, bool, string, error) {
    success := false
    changed := false
    value := ""
    majority := this.peers.GetPeerCount()/2+1
    replyCount := uint64(0)
    promiseCount := this.peers.GetSkipPromiseCount()
    highestAccepted := proposal.Default()

    for promiseCount < majority && replyCount < peerCount {
        var promise acceptor.PrepareResp
        select {
        case reply := <- endpoint:
            promise = *reply.Data.(*acceptor.PrepareResp)
            replyCount++
        case <- time.After(time.Second):
            return success, changed, value, nil
        }

        if promise.PromiseAccepted {
            promiseCount++

            if promise.AcceptedProposalId.IsGreaterThan(highestAccepted) {
                highestAccepted = promise.AcceptedProposalId
                changed = true
                value = promise.AcceptedValue
            } else {
                this.peers.SetPromiseRequirement(promise.RoleId, !promise.NoMoreAccepted)
            }
        }
    }

    fmt.Println("[ PROPOSER", this.roleId, "] Processed", replyCount, "replies,", promiseCount, "promises.")
    success = promiseCount >= majority
    return success, changed, value, nil
}

// Receves replies to proposal
func (this *ProposerRole) recvAccepts(request acceptor.ProposalReq, peerCount uint64, endpoint <-chan clusterpeers.Response) (bool, error) {
    majority := peerCount/2+1
    acceptCount := uint64(0)
    received := make(map[uint64]bool)

    for acceptCount < majority {
        var response acceptor.ProposalResp
        select {
            case reply := <- endpoint:
                response = *reply.Data.(*acceptor.ProposalResp)
                received[response.RoleId] = true
            case <- time.After(time.Second):
                return false, nil
        }

        if request.ProposalId.IsGreaterThan(response.AcceptedId) ||
            request.ProposalId == response.AcceptedId {
            acceptCount++
        } else {
            this.peers.SetPromiseRequirement(response.RoleId, true)
            return false, nil
        }

        if request.FirstUnchosenIndex > response.FirstUnchosenIndex {
            go this.notifyOfSuccess(response.RoleId, request.FirstUnchosenIndex, response.FirstUnchosenIndex)
        }
    }

    go this.processAllAccepts(request, peerCount, received, endpoint)

    return true, nil
}
