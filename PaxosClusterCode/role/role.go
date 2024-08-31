package role

import (
    "os"
    "fmt"
    "net/rpc"
	"github/paxoscluster/clusterpeers"
	"github/paxoscluster/replicatedlog"
	"github/paxoscluster/proposal/manager"
)

func LaunchNode(roleId, )  {
	cluster, roleId, address, err := clusterpeers.ConstructCluster(assignedId, disk)

	log, err := replicatedlog.ConstructLog(roleId, disk)

	proposals, err := manager.ConstructProposalManager(roleId, disk)

	acceptorRole := acceptor.Construct(roleId, log)
	proposerRole := proposer.Construct(roleId, proposals, log, cluster)
	
	handler := rpc.NewServer()
    err = handler.Register(acceptorRole)
    if err != nil { return address, err }
    err = handler.Register(proposerRole)
    if err != nil { return address, err }
    err = cluster.Listen(handler)
    if err != nil { return address, err }

	 // Connects to peers
	 go cluster.Connect()

	  // Dispatches heartbeat signal
	  go func() {
        for {
            go cluster.BroadcastHeartbeat(roleId)
            time.Sleep(time.Second)
        }
	}()
	
	// Begins leader election
    go proposer.Run(proposerRole)
}