package storage_mining

// SectorState is an enum of
// SectorCommitted, SectorRecovering, SectorActive, SectorFaulted, SectorCleared
// FaultCount is only relevant for SectorRecovering and SectorFaulted.
// FaultCount is relevant in accounting for the number
// of consecutive proving periods that a sector is faulted.
// Sectors that faulted more than CONSECUTIVE_FAULT_COUNT_LIMIT
// in a row will result in Sectors getting cleared and miners penalized.
// The enum is written in this awkward way because of golang limitation.

type SectorState struct {
	StateNumber int
	FaultCount int
}

func SectorCommitted() SectorState {
	return SectorState {
		StateNumber: 0,
		FaultCount: -1,
	}
}

func SectorRecovering(count int) SectorState {
	return SectorState {
		StateNumber: 1,
		FaultCount: count,
	}
}

func SectorActive() SectorState {
	return SectorState {
		StateNumber: 2,
		FaultCount: -1,
	}
}

func SectorFaulted(count int) SectorState {
	return SectorState {
		StateNumber: 3,
		FaultCount: count,
	}
}

func SectorCleared() SectorState {
	return SectorState {
		StateNumber: 4,
		FaultCount: -1,
	}
}
