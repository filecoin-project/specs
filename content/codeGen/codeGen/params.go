package codeGen

type HashRole int

const (
	HashRole_TicketVDFInputFromVRFOutput HashRole = 1
	HashRole_ElectionSeedFromVRFOutput   HashRole = 2
	HashRole_BlockSig                    HashRole = 3
)

const (
	Param_ElectionLookback Word = 50
)
