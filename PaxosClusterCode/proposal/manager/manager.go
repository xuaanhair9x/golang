package manager

type ProposalManager struct {
    roleId uint64
    proposalCount uint64
    currentId proposal.Id
    disk *recovery.Manager
    exclude sync.Mutex
}

func ConstructProposalManager(roleId, disk)  {
	id, err := disk.RecoverCurrentProposalId(roleId)

    newManager := ProposalManager {
        roleId: roleId,
        proposalCount: id.Sequence,
        currentId: id,
        disk: disk,
    }
	return &newManager, nil
}

// Generates & returns a new proposal ID
func (this *ProposalManager) GenerateNextProposalId() (proposal.Id, error) {
    this.exclude.Lock()
    defer this.exclude.Unlock()

    this.proposalCount++
    this.currentId = proposal.ConstructProposalId(this.roleId, this.proposalCount)

    err := this.disk.UpdateCurrentProposalId(this.roleId, this.currentId)
    if err != nil { return proposal.Default(), err }

    return this.currentId, nil
}