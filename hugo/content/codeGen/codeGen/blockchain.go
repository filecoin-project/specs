package codeGen

type Blockchain interface {
	TipsetAtEpoch(T Word)      Tipset
	CurrentConfirmedEpoch()    Word
	CurrentMiningEpoch()       Word
}

type BlockMiner interface {
	Address()             Address
	Blockchain()          Blockchain
	VRFKeyPair()          VRFKeyPair
	SigKeyPair()          SigKeyPair
	PowerFraction()       Fraction
	MineBlock([]Message)  Block
}

func (miner *BlockMinerI) MineBlock(messages []Message) Block {
	var chain = miner.Blockchain()
    T := chain.CurrentMiningEpoch()
    K := Param_ElectionLookback

    if T - K < 0 {
        panic("");  // TODO: handle genesis block
    }

	var parentTipset    = chain.TipsetAtEpoch(T)
	var lookbackTipset  = chain.TipsetAtEpoch(T - K)

	var currSeed = lookbackTipset.ElectionSeed()
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
			ret.blockSig = miner.SigKeyPair().Sign(ret.ComputeUnsignedFingerprint())
			return ret
        }
    }
}

func ComputeBlockWeight(parentTipset Tipset, tickets []Ticket) BlockWeight {
	panic("TODO")
}



////////////////////
// Implementation //
////////////////////

type BlockMinerI struct {
	// TODO
}

func (m *BlockMinerI) Address() Address {
	panic("TODO")
}

func (m *BlockMinerI) Blockchain() Blockchain {
	panic("TODO")
}

func (m *BlockMinerI) VRFKeyPair() VRFKeyPair {
	panic("TODO")
}

func (m *BlockMinerI) SigKeyPair() SigKeyPair {
	panic("TODO")
}

func (m *BlockMinerI) PowerFraction() Fraction {
	panic("TODO")
}
