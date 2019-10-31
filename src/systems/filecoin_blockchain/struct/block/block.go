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

	dist := chain.HeadEpoch() - epoch
	current := chain.HeadTipset()
	parents := current.Parents()
	for i := 0; i < int(dist); i++ {
		current = parents
		parents = current.Parents()
	}

	return current
}

func (chain *Chain_I) TicketAtEpoch(epoch ChainEpoch) Ticket {
	ts := chain.TipsetAtEpoch(epoch)
	return ts.MinTicket()
}

func (chain *Chain_I) HeadEpoch() ChainEpoch {
	panic("")
}

func (chain *Chain_I) HeadTipset() Tipset {
	panic("")
}

// should return the tipset from the nearest epoch to epoch containing a Tipset
// that is from the closest epoch less than or equal to epoch
func (bl *Block_I) TipsetAtEpoch(epoch ChainEpoch) Tipset {

	current := bl.Header_.Parents()
	parents := current.Parents()
	for current.Epoch() > epoch {
		current = parents
		parents = current.Parents()
	}
	return current
}

// should return the ticket from the Tipset generated at the nearest height leq to epoch
func (bl *Block_I) TicketAtEpoch(epoch ChainEpoch) Ticket {
	ts := bl.TipsetAtEpoch(epoch)
	return ts.MinTicket()
}
