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
func SHA256(input util.Bytes) util.Bytes {
	ret := make([]byte, 0)
	return ret
}

func sliceEqual(a util.Bytes, b util.Bytes) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func toBytes(e ChainEpoch) util.Bytes {
	ret := make([]byte, 0)
	return ret
}

func (chain *Chain_I) TicketOutputAtEpoch(epoch ChainEpoch) util.Bytes {
	ts := chain.TipsetAtEpoch(epoch)
	if ts.Epoch() != epoch {
		// there was no tipset at wanted epoch
		// craft ticket from prior valid ticket
		shaInput := ts.MinTicket().Output()
		for _, b := range toBytes(epoch) {
			shaInput = append(shaInput, b)
		}

		return SHA256(shaInput)
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
