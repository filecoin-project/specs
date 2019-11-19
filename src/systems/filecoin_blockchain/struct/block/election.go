package block

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

// TODO: add SHA256 to filcrypto
// TODO: import SHA256 from filcrypto
var SHA256 = func([]byte) []byte { return nil }

func (tix *Ticket_I) Output() bytes {
	return SHA256(tix.VRFResult_.Output_)
}

func (tix *Ticket_I) Verify(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return tix.VRFResult_.Verify(input, pk) && tix.Output_ == SHA256(tix.VRFResult_.Output_)
}

func (ep *ElectionProof_I) ValidateSyntax() bool {
	return ep.VRFResult_.ValidateSyntax()
}

func (ep *ElectionProof_I) Verify(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return ep.VRFResult_.Verify(input, pk)
}
