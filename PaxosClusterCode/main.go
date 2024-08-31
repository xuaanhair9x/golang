package main


import (
    "os"
    "fmt"
    "net/rpc"
	"github/paxoscluster/role"
	"github/paxoscluster/recovery"
)


func main()  {

	// Inital functions get data from files
	disk, err := recovery.ConstructManager()

	// Get Role Id from command line
	address, err := role.LaunchNode(roleId, disk)

	// Use this address to Receive value from client 
	
	cxn, err := rpc.Dial("tcp", nodeAddress)

	for {
        var input string
        fmt.Scanln(&input)
        var output string
        err = cxn.Call("ProposerRole.Replicate", &input, &output)
        if err != nil { fmt.Println(err) }
    }
}