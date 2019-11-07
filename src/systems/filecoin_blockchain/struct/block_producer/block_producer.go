package block_producer

<<<<<<< HEAD
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"

func (bp *BlockProducer_I) GenerateBlock(ep block.ElectionProof, T0 block.Ticket, ts block.Tipset, minerAddr addr.Address) {
	panic("TODO")
}

func (bp *BlockProducer_I) AssembleBlock(ep block.ElectionProof, T0 block.Tipset, ts block.Tipset, minerAddr addr.Address, messages []msg.Message) block.Block {
	panic("TODO")
}
=======
/*
func (miner *BlockProducer_I) MineBlock(messages []Message) Block {
	var chain = miner.Blockchain()
    T := chain.CurrentMiningEpoch()
    K := Param_ElectionLookback

    if T - K < 0 {
        panic("");  // TODO: handle genesis block
    }

	var parentTipset    = chain.TipsetAtEpoch(T)
	var lookbackTipset  = chain.TipsetAtEpoch(T - K)

	var currSeed = append([]byte("ELECTION"),lookbackTipset.ElectionSeed()...)
    var currTicket Ticket
    var tickets []Ticket

    for {
        currTicket = Ticket.Generate(nil, miner.VRFKeyPair(), currSeed)
        tickets = append(tickets, currTicket)

        if currTicket.IsWinning(miner.PowerFraction()) {
			ret := &BlockI {
				minerAddress: miner.Address(),
				tickets:      tickets,
				parentTipset: parentTipset,
				weight:       ComputeBlockWeight(parentTipset, tickets),
				height:       T + 1,
				stateTree:    parentTipset.StateTree(),
				messages:     messages,
				timestamp:    CurrentTime(),
				blockSig:     nil,
			};
			msg := append([]byte("BLOCK"),ret.ComputeUnsignedFingerprint()...)
			ret.blockSig = miner.SigKeyPair().Sign(msg)
			return ret
        }
    }
}
*/
>>>>>>> made changes to define domain sep tag
