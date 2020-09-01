package helpers

import "time"

func AvailableAt(delay time.Duration) int64 {
    return parseDateInterval(delay).Unix()
}

func parseDateInterval(delay time.Duration) time.Time {
    return time.Now().Add(delay)
}