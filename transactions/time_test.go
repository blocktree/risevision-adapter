package transactions

import (
	"fmt"
	"testing"
	"time"

)

var (
	defaultTime      = time.Date(2016, 4, 24, 17, 0, 0, 0, time.UTC)
	defaultTimeMs    = defaultTime.UnixNano() / int64(time.Millisecond)
	defaultEpochTime = uint32(20)
)

func TestGetTimeFromBlockchainEpoch(t *testing.T) {
	if val := getTimeFromBlockchainEpoch(defaultTimeMs); val != defaultEpochTime {
		t.Errorf("GetTimeFromBlockchainEpoch(%v)=%v; want %v", defaultTime, val, defaultEpochTime)
	}
}

func TestGetTimeWithOffset(t *testing.T) {
	if val := getTimeWithOffset(defaultTimeMs, int64(3)); val != defaultEpochTime+3 {
		t.Errorf("GetTimeWithOffset(%v)=%v; want %v", defaultTime, val, defaultEpochTime)
	}

	if val := getTimeWithOffset(defaultTimeMs, int64(-3)); val != defaultEpochTime-3 {
		t.Errorf("GetTimeWithOffset(%v)=%v; want %v", defaultTime, val, defaultEpochTime)
	}
}
func TestGetCurrentTimeWithOffset(t *testing.T) {
	fmt.Println(GetCurrentTimeWithOffset(0))
}


