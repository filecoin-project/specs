package block

import (
	util "github.com/filecoin-project/specs/util"
)

func SmallerBytes(a, b util.Bytes) util.Bytes {
	if util.CompareBytesStrict(a, b) > 0 {
		return b
	}
	return a
}

func (chain *Chain_I) TipsetAtEpoch(epoch ChainEpoch) Tipset {
	panic("")

	// dist := chain.HeadEpoch() - epoch
	// current := chain.HeadTipset()
	// parents := current.Parents()
	// for i := 0; i < dist; i++ {
	// 	current = parents
	// 	parents = current.Parents
	// }

	// return current
}

func (chain *Chain_I) TicketAtEpoch(epoch ChainEpoch) Ticket {
	ts := chain.TipsetAtEpoch(epoch)
	return ts.MinTicket()
}

func (chain *Chain_I) FinalizedEpoch() ChainEpoch {
	panic("")
	// ep := chain.HeadEpoch()
	// return ep - GetFinality()
}

func (chain *Chain_I) HeadEpoch() ChainEpoch {
	panic("")
}

func (chain *Chain_I) HeadTipset() Tipset {
	panic("")
}

func (bl *Block_I) TipsetAtEpoch(epoch ChainEpoch) Tipset_I {
	panic("")
	// dist := bl.Epoch - epoch - 1
	// current := bl.ParentTipset
	// parents := current.Parents
	// for i := range dist {
	// 	current = parents
	// 	parent = current.Parents
	// }
	// return current
}

func (bl *Block_I) TicketAtEpoch(epoch ChainEpoch) Ticket {
	ts := bl.TipsetAtEpoch(epoch)
	return ts.MinTicket()
}
