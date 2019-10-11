package block

import (
	filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {

	return tix.vrfResult_.ValidateSyntax() &&
		tix.vdfResult_.ValidateSyntax()

}

func (tix *Ticket_I) Validate(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return tix.vrfResult_.Verify(input, pk) &&
		tix.vdfResult_.Verify(tix.vdfResult_.Output())
}

func (ep *ElectionProof_I) ValidateSyntax() bool {
	return ep.vrfResult_.ValidateSyntax()
}

func (ep *ElectionProof_I) Validate(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return ep.vrfResult_.Verify(input, pk)
}

func (ep *ElectionProof_I) IsWinning(power PowerFraction) bool {
	panic("TODO")
}
