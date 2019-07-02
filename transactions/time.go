package transactions

import (
	"math"
	"time"
)

// GetCurrentTimeWithOffset returns the current blockchain time with an offset
func GetCurrentTimeWithOffset(offset int64) uint32 {
	return uint32(time.Now().Unix() - 1464109200)
}

func getTimeWithOffset(timestamp, offset int64) uint32 {
	timeWithOffset := timestamp + offset*1000
	return getTimeFromBlockchainEpoch(timeWithOffset)
}

func getTimeFromBlockchainEpoch(timestamp int64) uint32 {
	return uint32(math.Floor(float64(timestamp-epochTimeMs) / 1000))
}
