package block

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

func (tix *Ticket_I) Verify(randomness util.Bytes, pk filcrypto.VRFPublicKey, minerActorAddr addr.Address) bool {

	inputRand := Serialize_TicketProductionInput(&TicketProductionInput_I{
		PastTicket_: randomness,
		MinerAddr_:  minerActorAddr,
	})
	input := filcrypto.DomainSeparationTag_TicketProduction.DeriveRand(inputRand)

	return tix.VRFResult_.Verify(input, pk)
}

func (tix *Ticket_I) DrawRandomness(epoch ChainEpoch) util.Bytes {
	input := Serialize_TicketDrawingInput(&TicketDrawingInput_I{
		PastTicket_: tix.Output(),
		Epoch_:      epoch,
	})
	return filcrypto.DomainSeparationTag_TicketDrawing.DeriveRand(input)
}

func (ep *ElectionProof_I) ValidateSyntax() bool {
	panic("TODO")
	return ep.VRFResult_.ValidateSyntax()
}

func (ep *ElectionProof_I) Verify(partialTicket util.Bytes, pk filcrypto.VRFPublicKey) bool {
	panic("TODO")
	var input []byte
	input = append(input, partialTicket...)
	return ep.VRFResult_.Verify(input, pk)
}
