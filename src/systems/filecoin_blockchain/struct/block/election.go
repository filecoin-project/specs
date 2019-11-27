package block

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

func (tix *Ticket_I) Verify(priorTixRandomness util.Bytes, pk filcrypto.VRFPublicKey, minerAddr addr.Address) bool {
	var input []byte
	input = append(input, byte(filcrypto.DomainSeparationTag_Case_Ticket))
	input = append(input, priorTixRandomness...)
	input = append(input, byte(filcrypto.InputDelimeter_Case_Bytes))
	input = append(input, minerAddr.AddrToBytes()...)
	return tix.VRFResult_.Verify(input, pk)
}

func (tix *Ticket_I) DrawRandomness(epoch ChainEpoch) util.Bytes {
	var input []byte
	input = append(input, tix.Output()...)
	input = append(input, byte(filcrypto.InputDelimeter_Case_Bytes))
	input = append(input, byte(epoch))
	return filproofs.HashBytes_SHA256Hash(input)
}

func (tix *Ticket_I) DrawRandomness(minerAddr addr.Address, epoch ChainEpoch) util.Bytes {
	input := tix.Output()
	input = append(input, addrToLittleEndianBytes(minerAddr)...)
	input = append(input, epochToLittleEndianBytes(epoch)...)
	return input
}

func (ep *ElectionProof_I) ValidateSyntax() bool {
	return ep.VRFResult_.ValidateSyntax()
}

func (ep *ElectionProof_I) Verify(partialTicket util.Bytes, pk filcrypto.VRFPublicKey) bool {
	var input []byte
	input = append(input, byte(filcrypto.DomainSeparationTag_Case_PoSt))
	input = append(input, partialTicket...)
	return ep.VRFResult_.Verify(input, pk)
}
