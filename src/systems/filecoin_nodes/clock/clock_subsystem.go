package clock

import "time"

// UTCSyncPeriod notes how often to sync the UTC clock with an authoritative
// source, such as NTP, or a very precise hardware clock.
var UTCSyncPeriod = time.Hour

// EpochDuration is a constant that represents the duration in seconds
// of a blockchain epoch.
var EpochDuration = UnixTime(15)

func (_ *UTCClock_I) NowUTCUnix() UnixTime {
	return UnixTime(time.Now().Unix())
}

// EpochAtTime returns the ChainEpoch corresponding to time `t`.
// It first subtracts GenesisTime, then divides by EpochDuration
// and returns the resulting number of epochs.
func (c *ChainEpochClock_I) EpochAtTime(t UnixTime) ChainEpoch {
	difference := t - c.GenesisTime()
	epochs := difference / EpochDuration
	return ChainEpoch(epochs)
}
