package notifier

import (
	"time"
)

const (
	DEFAULT_RATE_LIMITER_DURATION = 60
)

// rateLimiterHook is a structure for rateLimiterHook configuration
type rateLimiterHook struct {
	// The value which controls the duration of this rate limiter, default 60 min
	Duration time.Duration

	// The value which controls the maximum number events allowed in Duration
	MaxEvents int

	// Stores the timestamps of recorded events
	eventTimestamps []time.Time

	// index to the last timestamp
	lastEventIndex int
}

func NewRateLimiterHook(maxEvents int) *rateLimiterHook {
	r := rateLimiterHook{
		MaxEvents:       maxEvents,
		eventTimestamps: make([]time.Time, maxEvents),
		lastEventIndex:  -1,
		Duration:        time.Minute * DEFAULT_RATE_LIMITER_DURATION,
	}

	return &r
}

func (r *rateLimiterHook) RecordEvent() {
	if r.lastEventIndex == r.MaxEvents-1 {
		r.lastEventIndex = 0
	} else {
		r.lastEventIndex = r.lastEventIndex + 1
	}

	r.eventTimestamps[r.lastEventIndex] = time.Now()
}

func (r *rateLimiterHook) IsLimitCrossed() bool {

	limitCrossed := false

	defer func() {
		if r := recover(); r != nil {
			// Handling panic. A false will be returned in this situation
		}
	}()

	if r.lastEventIndex == -1 {
		return limitCrossed
	}

	lastEventTimeStamp := r.getLastEventTimestamp()
	oldestEventTimeStamp, completeRoundFlag := r.getOldestEventTimestamp()

	if lastEventTimeStamp.Equal(oldestEventTimeStamp) {
		return limitCrossed
	}

	diffDuration := lastEventTimeStamp.Sub(oldestEventTimeStamp).Seconds()

	if diffDuration <= r.Duration.Seconds() && completeRoundFlag {
		limitCrossed = true
	}

	return limitCrossed
}

func (r *rateLimiterHook) getLastEventTimestamp() time.Time {
	return r.eventTimestamps[r.lastEventIndex]
}

func (r *rateLimiterHook) getOldestEventTimestamp() (time.Time, bool) {

	completeRoundFlag := false

	if r.lastEventIndex == r.MaxEvents-1 {
		completeRoundFlag = true
		return r.eventTimestamps[0], completeRoundFlag
	} else if r.lastEventIndex == 0 && r.eventTimestamps[r.lastEventIndex+1].Equal(time.Time{}) {
		return r.eventTimestamps[0], completeRoundFlag
	} else if r.eventTimestamps[r.lastEventIndex+1].Equal(time.Time{}) {
		return r.eventTimestamps[0], completeRoundFlag
	} else {
		completeRoundFlag = true
		return r.eventTimestamps[r.lastEventIndex+1], completeRoundFlag
	}
}
