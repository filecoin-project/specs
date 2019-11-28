package block_producer

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"

func (bp *BlockProducer_I) GenerateBlock(ep block.ElectionProof, T0 block.Ticket, ts block.Tipset, minerAddr addr.Address) {
	panic("TODO")
}

func (bp *BlockProducer_I) AssembleBlock(ep block.ElectionProof, T0 block.Tipset, ts block.Tipset, minerAddr addr.Address, messages []msg.SignedMessage) block.Block {
	panic("TODO")
}
