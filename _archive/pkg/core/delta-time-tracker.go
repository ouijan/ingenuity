package core

import "time"

type deltaTimeManager struct {
	lastStep time.Time
}

func (sm *deltaTimeManager) Step() float32 {
	now := time.Now()
	dt := now.Sub(sm.lastStep).Seconds()
	sm.lastStep = now
	return float32(dt)
}

func NewDeltaTimeTracker() *deltaTimeManager {
	return &deltaTimeManager{
		lastStep: time.Now(),
	}
}
