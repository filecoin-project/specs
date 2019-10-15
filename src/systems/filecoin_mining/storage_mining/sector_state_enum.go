package storage_mining

// SectorState is an enum of
// SectorCommitted, SectorRecovering, SectorActive, SectorFaulted, SectorCleared
// FaultCount is only relevant for SectorRecovering and SectorFaulted.
// FaultCount is relevant in accounting for the number
// of consecutive proving periods that a sector is faulted.
// Sectors that faulted more than CONSECUTIVE_FAULT_COUNT_LIMIT
// in a row will result in Sectors getting cleared and miners penalized.
// The enum is written in this awkward way because of golang limitation.

const (
	SectorCommittedStateNo  uint8 = 0
	SectorRecoveringStateNo uint8 = 1
	SectorActiveStateNo     uint8 = 2
	SectorFaultedStateNo    uint8 = 3
	SectorClearedStateNo    uint8 = 4
)

type SectorState struct {
	StateNumber uint8
	FaultCount  uint8
}

func SectorCommitted() SectorState {
	return SectorState{
		StateNumber: SectorCommittedStateNo,
		FaultCount:  0,
	}
}

func SectorRecovering(count uint8) SectorState {
	return SectorState{
		StateNumber: SectorRecoveringStateNo,
		FaultCount:  count,
	}
}

func SectorActive() SectorState {
	return SectorState{
		StateNumber: SectorActiveStateNo,
		FaultCount:  0,
	}
}

func SectorFaulted(count uint8) SectorState {
	return SectorState{
		StateNumber: SectorFaultedStateNo,
		FaultCount:  count,
	}
}

// SectorCleared is not directly represented in the spec
// when a Sector is cleared, it is deleted from sm
func SectorCleared() SectorState {
	return SectorState{
		StateNumber: SectorClearedStateNo,
		FaultCount:  0,
	}
}
