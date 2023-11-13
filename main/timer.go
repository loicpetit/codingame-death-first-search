package main

import (
	"fmt"
	"time"
)

type Timer struct {
	initStart      time.Time
	initDuration   time.Duration
	roundStart     time.Time
	roundDuration  time.Duration
	roundMin       time.Duration
	roundMax       time.Duration
	roundDurations []time.Duration
}

func (timer *Timer) startInit() {
	if timer == nil {
		return
	}
	timer.initStart = time.Now()
	timer.initDuration = 0
}

func (timer *Timer) endInit() {
	if timer == nil {
		return
	}
	timer.initDuration = time.Since(timer.initStart)
}

func (timer *Timer) startRound() {
	if timer == nil {
		return
	}
	timer.roundStart = time.Now()
	timer.roundDuration = 0
}

func (timer *Timer) endRound() {
	if timer == nil {
		return
	}
	timer.roundDuration = time.Since(timer.initStart)
	timer.roundDurations = append(timer.roundDurations, timer.roundDuration)
	if timer.roundMin == 0 || timer.roundDuration < timer.roundMin {
		timer.roundMin = timer.roundDuration
	}
	if timer.roundDuration > timer.roundMax {
		timer.roundMax = timer.roundDuration
	}
}

func (timer *Timer) roundAverage() time.Duration {
	if timer == nil || len(timer.roundDurations) == 0 {
		return 0
	}
	durations := time.Duration(0)
	for _, timerDuration := range timer.roundDurations {
		durations += timerDuration
	}
	return durations / time.Duration(len(timer.roundDurations))
}

func (timer *Timer) String() string {
	if timer == nil || timer.initDuration == 0 {
		return ""
	}
	if timer.roundDuration == 0 {
		return fmt.Sprintf("{init: %v}", timer.initDuration)
	}
	return fmt.Sprintf("{init: %v, round: %v, min: %v, max: %v, average: %v}", timer.initDuration, timer.roundDuration, timer.roundMin, timer.roundMax, timer.roundAverage())
}
