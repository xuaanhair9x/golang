package recovery

type Manager struct {
    sigint chan os.Signal
    exclude sync.Mutex
}

func ConstructManager()  {
	newManager := Manager {
        sigint: make(chan os.Signal, 1),
	}
	return &newManager, nil
}

func (this *Manager) RecoverLog(roleId uint64) ([]string, []proposal.Id, error) {
	/**
		- Read data from "log.csv"
		- Values = first column
		- acceptedProposals = remain columns as type Id
		type Id struct {
			RoleId uint64
			Sequence uint64
			Chosen bool
		} 
	**/
	return values, proposals, nil
}

func (this *Manager) RecoverCurrentProposalId()  {
	return recoverProposalId(roleId, "currentproposalid.csv")
}

func (this *Manager) RecoverMinProposalId()  {
	return recoverProposalId(roleId, "minproposalid.csv")
}
func recoverProposalId(roleId, filename string)  {
	/**
		Convert value in file as data
		type Id struct {
			RoleId uint64
			Sequence uint64
			Chosen bool
		}
	**/
}


// Updates a record in the replicated log on disk.
func (this *Manager) UpdateLogRecord(roleId uint64, index int, value string, id proposal.Id) error {
    this.exclude.Lock()
    defer this.exclude.Unlock()

    // Ensures existence of directory, creating it if necessary
    err := idempotentCreateRecoveryDir(roleId)
    if err != nil { return err }

    // If log file exists, read from it
    var records [][]string = nil
    exists, err := recoveryFileExists(roleId, LOG_FILENAME)
    if err != nil { return err }
    if exists {
        logFile, err := openRecoveryFile(roleId, LOG_FILENAME)
        if err != nil { return err }
        logFileReader := csv.NewReader(logFile)
        records, err = logFileReader.ReadAll()
        logFile.Close()
        if err != nil { return err }
    }

    // Modifies record
    blank := append([]string{""}, proposal.SerializeToCSV(proposal.Default())...)
    record := append([]string{value}, proposal.SerializeToCSV(id)...)
    for recordCount := len(records); recordCount <= index; recordCount++ {
        records = append(records, blank)
    }
    records[index] = record

    // Writes out log contents
    modifiedLogFile, err := createRecoveryFile(roleId, LOG_FILENAME)
    if err != nil { return err }
    defer modifiedLogFile.Close()
    logFileWriter := csv.NewWriter(modifiedLogFile)
    return logFileWriter.WriteAll(records)
}

