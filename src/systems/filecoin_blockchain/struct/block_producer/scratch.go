package block_producer

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
