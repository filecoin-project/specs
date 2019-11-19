package block

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

func (tix *Ticket_I) Output() util.Bytes {
	return SHA256(tix.VRFResult().Output())
}

func (tix *Ticket_I) Verify(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return tix.VRFResult_.Verify(input, pk) && sliceEqual(tix.Output(), SHA256(tix.VRFResult().Output()))
}

func (ep *ElectionProof_I) ValidateSyntax() bool {
	return ep.VRFResult_.ValidateSyntax()
}

func (ep *ElectionProof_I) Verify(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return ep.VRFResult_.Verify(input, pk)
}
