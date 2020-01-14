package block

import (
	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

func (tix *Ticket_I) Verify(randomness util.Bytes, pk filcrypto.VRFPublicKey, minerActorAddr addr.Address) bool {

	input := filcrypto.DeriveRandWithMinerAddr(filcrypto.DomainSeparationTag_TicketProduction, randomness, minerActorAddr)
	return tix.VRFResult_.Verify(input, pk)
}

func (tix *Ticket_I) DrawRandomness(epoch abi.ChainEpoch) util.Bytes {
	return filcrypto.DeriveRandWithEpoch(filcrypto.DomainSeparationTag_TicketDrawing, tix.Output(), int(epoch))
}
