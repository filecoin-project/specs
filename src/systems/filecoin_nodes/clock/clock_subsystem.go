package clock

import "time"

func (_ *WallClock_I) NowUTC() Time {
  return Time(time.Now().Unix())
}

func (_ *WallClock_I) NowUTCNano() Time {
  return Time(time.Now().UnixNano())
}
