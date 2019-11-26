package block

import (
	"bytes"

	clock "github.com/filecoin-project/specs/systems/filecoin_nodes/clock"
)

func (ts *Tipset_I) MinTicket() Ticket {
	var ret Ticket

	for _, currBlock := range ts.Blocks() {
		tix := currBlock.Ticket()
		if ret == nil {
			ret = tix
		} else {
			smaller := SmallerBytes(tix.Output(), ret.Output())
			if bytes.Equal(smaller, tix.Output()) {
				ret = tix
			}
		}
	}

	return ret
}

func (ts *Tipset_I) LatestTimestamp() clock.Time {
	var latest clock.Time
	for _, blk := range ts.Blocks_ {
		if blk.Timestamp() > latest || latest == clock.Time(0) {
			latest = blk.Timestamp()
		}
	}
	return latest
}
