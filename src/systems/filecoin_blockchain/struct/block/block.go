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

// will return tipset from closest prior (or equal) epoch with a tipset
// return epoch should be checked accordingly
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

// TODO: add SHA256 to filcrypto
// TODO: import SHA256 from filcrypto
var SHA256 = func([]byte) []byte { return nil }

func (chain *Chain_I) TicketOutputAtEpoch(epoch ChainEpoch) Bytes {
	ts := chain.TipsetAtEpoch(epoch)
	if ts.Epoch() != epoch {
		// there was no tipset at wanted epoch
		// craft ticket from prior valid ticket
		return SHA256(append(ts.MinTicket(), epoch))
	}
	return ts.MinTicket().Output()
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
