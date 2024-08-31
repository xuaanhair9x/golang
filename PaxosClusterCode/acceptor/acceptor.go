package acceptor

/*
 * Acceptor Role
 */
 type AcceptorRole struct {
    roleId uint64
    log *replicatedlog.Log
}


// Request sent out by proposer during prepare phase
type PrepareReq struct {
    ProposalId proposal.Id
    Index int
}

// Response sent by acceptors during prepare phase
type PrepareResp struct {
    PromiseAccepted bool
    AcceptedProposalId proposal.Id
    AcceptedValue string
    NoMoreAccepted bool
    RoleId uint64
}

func Construct(roleId uint64, log *replicatedlog.Log)  {
	this := AcceptorRole{roleId, log}
    return &this
}

func (this *AcceptorRole) Prepare(req *PrepareReq, reply *PrepareResp) error {
    minProposalId := this.log.GetMinProposalId()
    fmt.Println("[ ACCEPTOR", this.roleId, "] Prepare: considering proposal", req.ProposalId, 
                "vs", minProposalId, "for index", req.Index)
    logEntry := this.log.GetEntryAt(req.Index)
    reply.PromiseAccepted = req.ProposalId.IsGreaterThan(minProposalId)
    reply.AcceptedProposalId = logEntry.AcceptedProposalId
    reply.AcceptedValue = logEntry.Value
    reply.NoMoreAccepted = this.log.NoMoreAcceptedPast(req.Index)
    reply.RoleId = this.roleId
    this.log.UpdateMinProposalId(req.ProposalId)
    return nil
}

// Request sent out by proposer during proposal phase
type ProposalReq struct {
    ProposalId proposal.Id
    Index int
    Value string
    FirstUnchosenIndex int
}

// Response sent by acceptors during proposal phase
type ProposalResp struct {
    AcceptedId proposal.Id
    RoleId uint64
    FirstUnchosenIndex int
}

func (this *AcceptorRole) Accept(proposal *ProposalReq, reply *ProposalResp) error {
    fmt.Println("[ ACCEPTOR", this.roleId, "] Proposal: considering proposal", proposal.ProposalId,
                "of", proposal.Value, "for index", proposal.Index)
    this.log.MarkAsChosen(proposal.ProposalId, proposal.FirstUnchosenIndex)
    minProposalId := this.log.GetMinProposalId()
    if proposal.ProposalId.IsGreaterThan(minProposalId) || proposal.ProposalId == minProposalId {
        this.log.SetEntryAt(proposal.Index, proposal.Value, proposal.ProposalId)
    }
    reply.AcceptedId = minProposalId
    reply.RoleId = this.roleId
    reply.FirstUnchosenIndex = this.log.GetFirstUnchosenIndex()
    return nil
}