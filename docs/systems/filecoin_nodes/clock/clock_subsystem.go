package clock

import "time"

// UTCMaxDrift is how large the allowable drift is in Filecoin's use of UTC time.
var UTCMaxDrift = time.Second

// UTCSyncPeriod notes how often to sync the UTC clock with an authoritative
// source, such as NTP, or a very precise hardware clock.
var UTCSyncPeriod = time.Hour

// EpochDuration is a constant that represents the UTC time duration
// of a blockchain epoch.
var EpochDuration = time.Second * 15

// ISOFormat is the ISO timestamp format we use, in Go time package notation.
var ISOFormat = "2006-01-02T15:04:05.999999999Z"

func (_ *UTCClock_I) NowUTC() Time {
	return Time(time.Now().Format(ISOFormat))
}

func (_ *UTCClock_I) NowUTCUnix() UnixTime {
	return UnixTime(time.Now().Unix())
}

func (_ *UTCClock_I) NowUTCUnixNano() UnixTimeNano {
	return UnixTimeNano(time.Now().UnixNano())
}

// EpochAtTime returns the ChainEpoch corresponding to t.
// It first subtracts GenesisTime, then divides by EpochDuration
// and returns the resulting number of epochs.
func (c *ChainEpochClock_I) EpochAtTime(t Time) ChainEpoch {
	g1 := c.GenesisTime()
	g2, err := time.Parse(ISOFormat, string(g1))
	if err != nil {
		// an implementation should probably not panic here
		// this is for simplicity of the spec
		panic(err)
	}

	t2, err := time.Parse(ISOFormat, string(t))
	if err != nil {
		panic(err)
	}

	difference := t2.Sub(g2)
	epochs := difference / EpochDuration
	return ChainEpoch(epochs)
}
