package storage_power_consensus

type ConsensusFaultType int

const (
	UncommittedPowerFault     ConsensusFaultType = 0
	DoubleForkMiningFault     ConsensusFaultType = 1
	ParentGrindingFault       ConsensusFaultType = 2
	SameForkDoubleMiningFault ConsensusFaultType = 3
)
