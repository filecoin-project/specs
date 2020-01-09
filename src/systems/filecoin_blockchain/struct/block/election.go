package block

import (
	actors "github.com/filecoin-project/specs/actors"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

func (tix *Ticket_I) Verify(randomness util.Bytes, pk filcrypto.VRFPublicKey, minerActorAddr addr.Address) bool {

	inputRand := Serialize_TicketProductionSeedInput(&TicketProductionSeedInput_I{
		PastTicket_: randomness,
		MinerAddr_:  minerActorAddr,
	})
	input := filcrypto.DeriveRand(filcrypto.DomainSeparationTag_TicketProduction, inputRand)

	return tix.VRFResult_.Verify(input, pk)
}

func (tix *Ticket_I) DrawRandomness(epoch actors.ChainEpoch) util.Bytes {
	input := Serialize_TicketDrawingSeedInput(&TicketDrawingSeedInput_I{
		PastTicket_: tix.Output(),
		Epoch_:      epoch,
	})
	return filcrypto.DeriveRand(filcrypto.DomainSeparationTag_TicketDrawing, input)
}
