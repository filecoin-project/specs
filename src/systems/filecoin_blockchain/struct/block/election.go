package block

import (
	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs-actors/actors/abi"
	crypto "github.com/filecoin-project/specs/algorithms/randomness"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

func (tix *Ticket_I) Verify(proof util.Bytes, digest util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return tix.VRFResult_.Verify(proof, pk) && digest == blake2b.Sum256(proof)
}
