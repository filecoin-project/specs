package block

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

func (tix *Ticket_I) Verify(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return tix.VRFResult_.Verify(input, pk)
}

func (tix *Ticket_I) DrawRandomness(minerAddr addr.Address, epoch ChainEpoch) util.Bytes {
	input := tix.Output()
	input = append(input, addrToLittleEndianBytes(minerAddr)...)
	input = append(input, epochToLittleEndianBytes(epoch)...)
	return input
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

func (ep *ElectionProof_I) Verify(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return ep.VRFResult_.Verify(append([]byte(filcrypto.ElectionTag), input...), pk)
}
