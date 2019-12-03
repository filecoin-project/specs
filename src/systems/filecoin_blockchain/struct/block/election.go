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

	// var input []byte
	// input = append(input, byte(filcrypto.DomainSeparationTag_Case_Ticket))
	// input = append(input, randomness1...)
	// input = append(input, byte(filcrypto.InputDelimeter_Case_Bytes))
	// input = append(input, minerAddr.ToBytes()...)
	inputRand := Serialize_TicketProductionInput(&TicketProductionInput_I{
		PastTicket_: randomness,
		MinerAddr_:  minerActorAddr,
	})
	input := filcrypto.DomainSeparationTag_TicketProduction.DeriveRand(inputRand)

	return tix.VRFResult_.Verify(input, pk)
}

func (tix *Ticket_I) DrawRandomness(epoch ChainEpoch) util.Bytes {
	ser := util.SerializeBytes(tix.Output())
	return filcrypto.DomainSeparationTag_TicketDrawing.DeriveRandWithIndex(ser, int(epoch))
	// var input []byte
	// input = append(input, tix.Output()...)
	// input = append(input, byte(epoch))
	// return SHA256(input)
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
