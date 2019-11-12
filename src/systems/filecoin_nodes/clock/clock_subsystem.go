package clock

import "time"

func (_ *Clock_I) NowUTC() Time {
	return Time(time.Now().Unix())
}

func (_ *Clock_I) NowUTCNano() Time {
	return Time(time.Now().UnixNano())
}

// given some starting epoch and time since, what epoch should this be?
func (self *Clock_I) epochAfterTime(start ChainEpoch, t Time) ChainEpoch {
	panic("")
	// return start + t/epochDuration
}

func (self *Clock_I) EpochAfterGenesis() ChainEpoch {
	panic("")
	// return self.epochAfterTime(0 + self.NowUTCNano() - self.genesisTime)
}

func (self *Clock_I) currentEpochStart() Time {
	panic("")
	// return EpochAfterGenesis() * self.epochDuration + genesisTime
}

func (self *Clock_I) CurrentEpochState() ChainEpoch {
	panic("")
	// if self.NowUTCNano()-self.currentEpochStart() > self.cutoffAt_ {
	// 	return ChainEpochState.PastCutoff
	// }
	// return ChainEpochState.Active
}
