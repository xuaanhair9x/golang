package clusterpeers

type Cluster struct {
    roleId uint64
    nodes map[uint64]Peer
    registerBadConnection chan uint64
    skipPromiseCount uint64
    disk *recovery.Manager
    exclude sync.Mutex
}
/**

skipPromiseCount ??
requirePromise ?? 

**/
type Peer struct {
    roleId uint64
    address string
    comm *rpc.Client
    requirePromise bool
}
func ConstructCluster(roleId uint64, disk *recovery.Manager) (*Cluster, uint64, string, error) {

	 // Builds peers map
	 peers := make(map[uint64]Peer)
	 for id, address := range addresses {
		 newPeer := Peer {
			 roleId: id,
			 address: address,
			 comm: nil,
			 requirePromise: true,
		 }
		 peers[id] = newPeer
	 }

	newCluster := Cluster {
        roleId: roleId,
        nodes: peers,
        registerBadConnection: make(chan uint64, 16),
        skipPromiseCount: 0,
        disk: disk,
	}
	
	go newCluster.connectionManager()

    return &newCluster, newCluster.roleId, address, nil
}

// Sets server to listen on this node's port
func (this *Cluster) Listen(handler *rpc.Server) error {

    // Listens on specified address
    ln, err := net.Listen("tcp", this.nodes[this.roleId].address)

    // Dispatches connection processing loop
    go func() {
        for {
            connection, err := ln.Accept()
            if err != nil { continue }
            go handler.ServeConn(connection)
        }
    }()

    return nil
}

// Initializes connections to cluster peers
func (this *Cluster) Connect() {

    for roleId, peer := range this.nodes {
        connection, err := rpc.Dial("tcp", peer.address)
        if err != nil {
            this.registerBadConnection <- roleId
        } else {
            peer.comm = connection
            this.nodes[roleId] = peer
        }
    }
}


// Triages connection complaints, organizes repair attempts
func (this *Cluster) connectionManager() {
    establishing := make(map[uint64]bool)
    connectionEstablished := make(chan uint64)
    for {
        select {
        case roleId := <- this.registerBadConnection:
            if !establishing[roleId] {
                fmt.Println("[ NETWORK", this.roleId, "] Attempting to establish connection to", roleId)
                establishing[roleId] = true
                go this.establishConnection(roleId, connectionEstablished)
            }
        case roleId := <- connectionEstablished:
            establishing[roleId] = false
            fmt.Println("[ NETWORK", this.roleId, "] Connection to", roleId, "has been established")
        }
    }
}

// Attempts to re-connect to the specified role
func (this *Cluster) establishConnection(roleId uint64, connectionEstablished chan<- uint64) {
    this.exclude.Lock()
    peer := this.nodes[roleId]
    this.exclude.Unlock()

    for {
        connection, err := rpc.Dial("tcp", peer.address)
        if err != nil {
            time.Sleep(time.Second)
            continue
        }

        this.exclude.Lock()
        peer = this.nodes[roleId] 
        peer.comm = connection
        this.nodes[roleId] = peer
        connectionEstablished <- roleId
        this.exclude.Unlock()
        return
    }
}


// Sends pulse to all nodes in the cluster
func (this *Cluster) BroadcastHeartbeat(roleId uint64) {

	// Send role of node to other nodes
	peerCount := len(this.nodes)
    endpoint := make(chan *rpc.Call, peerCount)
    for _, peer := range this.nodes {
        if peer.comm != nil {
            var reply uint64
            peer.comm.Go("ProposerRole.Heartbeat", &roleId, &reply, endpoint)
        }
	}
	
	// Records nodes which return the heartbeat signal
    received := make(map[uint64]bool)
    failures := false
    replyCount := 0
    for replyCount < peerCount {
        select {
        case reply := <- endpoint:
            if reply.Error == nil {
                id := *reply.Reply.(*uint64)
                received[id] = true 
            } else {
                failures = true
            }
            replyCount++
        case <- time.After(time.Second/2):
            failures = true
            replyCount = peerCount
        }
	}
	
	// Registers bad connections if reply was not received
    if failures {
        for roleId := range this.nodes {
            if !received[roleId] {
                peer := this.nodes[roleId]
                if !peer.requirePromise {
                    this.skipPromiseCount--
                }
                peer.requirePromise = true
                this.nodes[roleId] = peer
                this.registerBadConnection <- roleId // Attempts to re-connect to the specified role
            }
        }
	}
	
	// Broadcasts a prepare phase request to the cluster
func (this *Cluster) BroadcastPrepareRequest(request acceptor.PrepareReq) (uint64, <-chan Response) {
    this.exclude.Lock()
    defer this.exclude.Unlock()

    peerCount := uint64(0)
    nodeCount := uint64(len(this.nodes))
    endpoint := make(chan *rpc.Call, nodeCount)

    if this.skipPromiseCount < nodeCount/2+1 {
        for _, peer := range this.nodes {
            if peer.requirePromise && peer.comm != nil {
                var response acceptor.PrepareResp
                peer.comm.Go("AcceptorRole.Prepare", &request, &response, endpoint)
                peerCount++
            } 
        }
    } else {
        fmt.Println("[ NETWORK", this.roleId, "] Skipping prepare phase: know state of majority")
    }


    responses := make(chan Response, peerCount)
    go this.wrapReply(peerCount, endpoint, responses)
    return peerCount, responses 
}

// Wraps RPC return data to remove direct dependency of caller on net/rpc and improve testability
func (this *Cluster) wrapReply(peerCount uint64, endpoint <-chan *rpc.Call, forward chan<- Response) {
    replyCount := uint64(0)
    for replyCount < peerCount {
        select {
        case reply := <- endpoint:
            if reply.Error == nil {
                forward <- Response{reply.Reply}
            }
            replyCount++
        case <- time.After(2*time.Second):
            return
        }
    }
}


// Returns number of peers in cluster
func (this *Cluster) GetPeerCount() uint64 {
    this.exclude.Lock()
    defer this.exclude.Unlock()

    return uint64(len(this.nodes))
}

// Returns number of peers from which no promise is required
func (this *Cluster) GetSkipPromiseCount() uint64 {
    this.exclude.Lock()
    defer this.exclude.Unlock()

    return this.skipPromiseCount
}

// Broadcasts a proposal phase request to the cluster
func (this *Cluster) BroadcastProposalRequest(request acceptor.ProposalReq, filter map[uint64]bool) (uint64, <-chan Response) {
    this.exclude.Lock()
    defer this.exclude.Unlock()

    peerCount := uint64(0)
    endpoint := make(chan *rpc.Call, len(this.nodes)) 
    for roleId, peer := range this.nodes {
        if !filter[roleId] && peer.comm != nil {
            var response acceptor.ProposalResp
            peer.comm.Go("AcceptorRole.Accept", &request, &response, endpoint)
            peerCount++
        }
    }

    responses := make(chan Response, peerCount)
    go this.wrapReply(peerCount, endpoint, responses)
    return peerCount, responses 
}