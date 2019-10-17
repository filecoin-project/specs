package storage_mining

// SectorState is an enum of
// SectorCommitted, SectorRecovering, SectorActive, SectorFailing, SectorCleared
// FaultCount is only relevant for SectorRecovering and SectorFailing.
// FaultCount is relevant in accounting for the number
// of consecutive proving periods that a sector is Failing.
// Sectors that are in Failing for more than CONSECUTIVE_FAULT_COUNT_LIMIT
// in a row will result in Sectors getting cleared and miners penalized.
// The enum is written in this awkward way because of golang limitation.

type FaultCount uint8

const MAX_CONSECUTIVE_FAULTS = FaultCount(3)

const (
	SectorClearedSN    uint8 = 0
	SectorCommittedSN  uint8 = 1
	SectorActiveSN     uint8 = 2
	SectorRecoveringSN uint8 = 3
	SectorFailingSN    uint8 = 4
)

type SectorState struct {
	StateNumber uint8
	FaultCount  FaultCount
}

func SectorCommitted() SectorState {
	return SectorState{
		StateNumber: SectorCommittedSN,
		FaultCount:  0, // always zero for SectorCommitted
	}
}

func SectorRecovering(count FaultCount) SectorState {
	return SectorState{
		StateNumber: SectorRecoveringSN,
		FaultCount:  count,
	}
}

func SectorActive() SectorState {
	return SectorState{
		StateNumber: SectorActiveSN,
		FaultCount:  0, // always zero for SectorActive
	}
}

func SectorFailing(count FaultCount) SectorState {
	return SectorState{
		StateNumber: SectorFailingSN,
		FaultCount:  count,
	}
}

func SectorCleared() SectorState {
	return SectorState{
		StateNumber: SectorClearedSN,
		FaultCount:  0,
	}
}
