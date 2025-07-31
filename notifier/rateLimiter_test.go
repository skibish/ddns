package notifier

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestRateLimiter(t *testing.T) {

	testCases := []struct {
		tname        string
		maxThreshold int
		eventRun     int
		duration     time.Duration
		isErr        bool
	}{
		{
			tname:        "event run under limit | rate limit not crossed",
			maxThreshold: 5,
			eventRun:     3,
			duration:     time.Minute * 1,
			isErr:        false,
		},
		{
			tname:        "event run over limit | rate limit crossed",
			maxThreshold: 5,
			eventRun:     6,
			duration:     time.Minute * 1,
			isErr:        true,
		},
		{
			tname:        "no event run | rate limit not crossed",
			maxThreshold: 5,
			eventRun:     0,
			duration:     time.Minute * 1,
			isErr:        false,
		},
		{
			tname:        "event run equal to limit | rate limit crossed",
			maxThreshold: 5,
			eventRun:     5,
			duration:     time.Minute * 1,
			isErr:        true,
		},
		{
			tname:        "event run over limit over time | rate limit not crossed",
			maxThreshold: 5,
			eventRun:     10,
			duration:     time.Second * 3,
			isErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tname, func(t *testing.T) {
			is := is.New(t)

			r := NewRateLimiterHook(tc.maxThreshold)
			r.Duration = tc.duration

			for i := 0; i < tc.eventRun; i++ {
				r.RecordEvent()
				time.Sleep(time.Second * 1)
			}

			limitCrossed := r.IsLimitCrossed()

			is.True(tc.isErr == limitCrossed)
		})
	}

}
