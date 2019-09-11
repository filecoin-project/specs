package codeGen

type Ticket interface {
	VRFResult() VRFResult
	VDFResult() VDFResult

	Generate(k VRFKeyPair, seed ElectionSeed) Ticket
	IsWinning(power Fraction) bool
}

////////////////////
// Implementation //
////////////////////

func (__ *TicketI) Generate(k VRFKeyPair, seed ElectionSeed) Ticket {
	Assert(__ == nil)

	var vrfResult VRFResult = VRFEval(k, seed.ToBytesInternal())

	var vdfInput []byte = Hash(
		HashRole_TicketVDFInputFromVRFOutput,
		vrfResult.ToBytesInternal(),
	)
	var vdfResult VDFResult = VDFEval(vdfInput)

	return &TicketI{
		vrfResult,
		vdfResult,
	}
}

func (__ *TicketI) IsWinning(power Fraction) bool {
	panic("TODO")
}

type TicketI struct {
	vrfResult VRFResult
	vdfResult VDFResult
}

func (t *TicketI) VRFResult() VRFResult {
	return t.vrfResult
}

func (t *TicketI) VDFResult() VDFResult {
	return t.vdfResult
}
