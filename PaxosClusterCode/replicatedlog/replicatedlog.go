package replicatedlog


type LogEntry struct {
    Index int
    Value string
    AcceptedProposalId proposal.Id
}


func ConstructLog(roleId, disk)  {
	values, acceptedProposals, err := disk.RecoverLog(roleId)

    if err != nil { return nil, err }

	minProposalId, err := disk.RecoverMinProposalId(roleId)

    if err != nil { return nil, err }

    newLog := Log {
        roleId: roleId,
        values: values,
        acceptedProposals: acceptedProposals,
        minProposalId: minProposalId,
        firstUnchosenIndex: 0,
        disk: disk,
    }

    newLog.updateFirstUnchosenIndex()
    return &newLog, nil
}

// Updates the location of the first unchosen index; exclude MUST be locked before calling
func (this *Log) updateFirstUnchosenIndex() {
    limit := len(this.acceptedProposals)

	/**
	Check all accepted proposals, if there is a accept proposals is unchosen, set firstUnchosenIndex = idx of accept proposals
	**/
    for idx := this.firstUnchosenIndex; idx < limit; idx++ {
        if !this.acceptedProposals[idx].IsChosen() {
            this.firstUnchosenIndex = idx
            return
        } else {
            this.emit(idx) // only print on scream
        }
    }

	// If there is no unchosen in the list, set firstUnchosenIndex = len() 
    this.firstUnchosenIndex = len(this.acceptedProposals)
}

// Returns details of the log entry at the specified index
func (this *Log) GetEntryAt(index int) LogEntry {
	entry := LogEntry {
        Index: index,
        Value: "",
        AcceptedProposalId: proposal.Default(),
    }

    if index < len(this.values) && index < len(this.acceptedProposals) {
        entry.Value = this.values[index]
        entry.AcceptedProposalId = this.acceptedProposals[index]
    } 

    return entry
}


// Certifies that no proposals have been accepted past the specified index
func (this *Log) NoMoreAcceptedPast(index int) bool {
    if index+1 >= len(this.acceptedProposals) {
        return true
    }
    for _, proposalId := range this.acceptedProposals[index+1:] {
        if proposalId != proposal.Default() {
            return false
        }
    }
    return true
}

// Updates minProposalId to the greater of itself and the provided proposalId
func (this *Log) UpdateMinProposalId(proposalId proposal.Id) proposal.Id {
    this.exclude.Lock()
    defer this.exclude.Unlock()

    if proposalId.IsGreaterThan(this.minProposalId) {
        this.minProposalId = proposalId
        err := this.disk.UpdateMinProposalId(this.roleId, this.minProposalId)
        if err != nil {
            fmt.Println("[ LOG", this.roleId, "] Failed to write minProposalId update to disk")
        }
    }

    return this.minProposalId
}

// Sets the value of the log entry at the specified index
func (this *Log) SetEntryAt(index int, value string, proposalId proposal.Id) {
    this.exclude.Lock()
    defer this.exclude.Unlock()
    
    // Extends log as necessary
    if index >= len(this.values) || index >= len(this.acceptedProposals) {
        valuesDiff := index-len(this.values)+1
        proposalsDiff := index-len(this.acceptedProposals)+1
        this.values = append(this.values, make([]string, valuesDiff)...)
        this.acceptedProposals = append(this.acceptedProposals, make([]proposal.Id, proposalsDiff)...)
    } 

    if !this.acceptedProposals[index].IsChosen() &&
        (proposalId.IsGreaterThan(this.minProposalId) ||
        proposalId == this.minProposalId) {
        this.values[index] = value 
        this.acceptedProposals[index] = proposalId
        fmt.Println("[ LOG", this.roleId, "] Values:", this.values)
        fmt.Println("[ LOG", this.roleId, "] Proposals:", this.acceptedProposals)
        err := this.disk.UpdateLogRecord(this.roleId, index, value, proposalId)
        if err != nil {
            fmt.Println("[ LOG", this.roleId, "] Failed to write", proposalId, index, value, "to disk")
        }
    }

    // Updates firstUnchosenIndex if value is being chosen there
    if proposalId.IsChosen() && this.firstUnchosenIndex == index {
        this.updateFirstUnchosenIndex()
    }
}