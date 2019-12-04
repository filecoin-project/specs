package block_producer

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
)

func (bp *BlockProducer_I) GenerateBlock(ePoStInfo sector.OnChainPoStVerifyInfo, T0 block.Ticket, ts block.Tipset, minerAddr addr.Address) {
	panic("TODO")
}

func (bp *BlockProducer_I) AssembleBlock(ePoStInfo sector.OnChainPoStVerifyInfo, T0 block.Tipset, ts block.Tipset, minerAddr addr.Address, messages []msg.SignedMessage) block.Block {
	panic("TODO")
}
