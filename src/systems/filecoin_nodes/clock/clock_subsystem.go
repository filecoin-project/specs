package clock

import "time"

// UTCMaxDrift is how large the allowable drift is in Filecoin's use of UTC time.
var UTCMaxDrift = time.Second

// UTCSyncPeriod notes how often to sync the UTC clock with an authoritative
// source, such as NTP, or a very precise hardware clock.
var UTCSyncPeriod = time.Hour

// EpochDuration is a constant that represents the time in seconds
// of a blockchain epoch.
var EpochDuration = UnixTime(15)

func (_ *UTCClock_I) NowUTCUnix() UnixTime {
	return UnixTime(time.Now().Unix())
}

// EpochAtTime returns the ChainEpoch corresponding to t.
// It first subtracts GenesisTime, then divides by EpochDuration
// and returns the resulting number of epochs. If t is before
// GenesisTime zero is returned.
func (c *ChainEpochClock_I) EpochAtTime(t UnixTime) ChainEpoch {
	if t <= c.GenesisTime() {
		return ChainEpoch(0)
	}
	difference := t - c.GenesisTime()
	epochs := difference / EpochDuration
	return ChainEpoch(epochs)
}
