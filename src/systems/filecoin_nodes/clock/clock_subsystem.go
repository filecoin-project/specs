package clock

import "time"

func (_ *Clock_I) UTCTime() Time {
	return Time(time.Now().Unix())
}

func (_ *Clock_I) UnixNano() Time {
	return Time(time.Now().UnixNano())
}

func (self *Clock_I) CurrentEpoch() ChainEpoch {
	return self.currentEpoch_
}

func (self *Clock_I) CurrentEpochState() ChainEpoch {
	panic("")
	// if self.UnixNano()-self.currentEpochStart > self.cutoffAt_ {
	// 	return ChainEpochState.PastCutoff
	// }
	// return ChainEpochState.Active
}

func (self *Clock_I) ResetEpoch(genesisTime Time) struct{} {
	panic("")
	// self.currentEpoch = (self.UnixNano() - genesisTime) / epochDuration
}
