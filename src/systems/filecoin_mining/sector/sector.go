package sector

import (
	actors "github.com/filecoin-project/specs/actors"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	util "github.com/filecoin-project/specs/util"
)

var IMPL_FINISH = util.IMPL_FINISH

type Serialization = util.Serialization

// If a sector PreCommit appear at epoch T, then the corresponding ProveCommit
// must appear between epochs
//   (T + MIN_PROVE_COMMIT_SECTOR_EPOCH, T + MAX_PROVE_COMMIT_SECTOR_EPOCH)
// inclusive.
// TODO: placeholder epoch values -- will be set later
const MIN_PROVE_COMMIT_SECTOR_EPOCH = actors.ChainEpoch(5)
const MAX_PROVE_COMMIT_SECTOR_EPOCH = actors.ChainEpoch(10)

const (
	DeclaredFault StorageFaultType = 1 + iota
	DetectedFault
	TerminatedFault
)

func (amt *DealExpirationAMT_I) Size() int {
	return 0
}

func (amt *DealExpirationAMT_I) Add(key actors.ChainEpoch, value DealExpirationValue) {
	// helper function to add entry into the AMT
}

func (amt *DealExpirationAMT_I) ActiveDealIDs() []deal.DealID {
	ret := make([]deal.DealID, 0)
	return ret
}

// return last item in the expiration amt
func (q *DealExpirationAMT_I) LastDealExpiration() actors.ChainEpoch {
	ret := actors.ChainEpoch(0)
	return ret
}

// return deal IDs expiring in epoch range
func (q *DealExpirationAMT_I) ExpiredDealsInRange(start actors.ChainEpoch, end actors.ChainEpoch) []DealExpirationValue {
	ret := make([]DealExpirationValue, 0)
	return ret
}
