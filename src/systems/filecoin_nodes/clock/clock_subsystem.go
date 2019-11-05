package clock

import "time"

func (_ *Clock_I) UTCTime() Time {
	return Time(time.Now().Unix())
}

func (_ *Clock_I) UnixNano() Time {
	return Time(time.Now().UnixNano())
}

func (self *Clock_I) CurrentEpoch() ChainEpoch {
	// return self.currentEpoch
}

func (self *Clock_I) CurrentEpochState() ChainEpoch {
	// if self.UnixNano()-self.currentEpochStart > self.cutoffAt {
	// 	return ChainEpochState.PastCutoff
	// }
	// return ChainEpochState.Active
}

func (self *Clock_I) ResetEpoch(genesisTime Time) struct{} {
	// self.currentEpoch = (self.UnixNano() - genesisTime) / epochDuration
}
