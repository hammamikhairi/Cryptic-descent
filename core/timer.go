package core

import (
	"math/rand"
)

// TeleportTimer manages the countdown until the next teleportation
type TeleportTimer struct {
	maxTime           float32 // Maximum time before teleportation
	currentTime       float32 // Current countdown time
	tickFrequency     float32 // Time interval between ticks
	tickAcceleration  float32 // Speed at which ticks become more frequent
	lastTickTime      float32 // Time since the last tick
	teleportTriggered bool    // Indicates if teleportation has been triggered
}

// NewTeleportTimer creates and initializes a new teleport timer
func NewTeleportTimer(minTime, maxTime, tickFreq, tickAccel float32) *TeleportTimer {
	initialTime := minTime + rand.Float32()*(maxTime-minTime)
	return &TeleportTimer{
		maxTime:           initialTime,
		currentTime:       initialTime,
		tickFrequency:     tickFreq,
		tickAcceleration:  tickAccel,
		lastTickTime:      0,
		teleportTriggered: false,
	}
}

// Update updates the teleport timer and checks if it's time to teleport
func (t *TeleportTimer) Update(deltaTime float32) {
	if t.teleportTriggered {
		return // Skip update if teleportation already happened
	}

	t.currentTime -= deltaTime // Decrease the timer by elapsed time

	// Decrease the interval between ticks as we approach zero
	if t.currentTime < t.maxTime/3 {
		t.tickFrequency = max(0.1, t.tickFrequency-t.tickAcceleration*deltaTime)
	}

	// If time for the next tick has passed, reset the tick timer
	t.lastTickTime += deltaTime
	if t.lastTickTime >= t.tickFrequency {
		t.lastTickTime = 0 // Reset time since last tick
		t.TriggerTick()    // Handle tick logic (e.g., play tick sound)
	}

	// Check if timer has reached zero
	if t.currentTime <= 0 {
		t.teleportTriggered = true // Mark teleportation as triggered
	}
}

// Reset resets the teleport timer with a new random time
func (t *TeleportTimer) Reset(minTime, maxTime float32) {
	t.currentTime = minTime + rand.Float32()*(maxTime-minTime)
	t.tickFrequency = 1.0 // Reset tick frequency
	t.lastTickTime = 0    // Reset tick timer
	t.teleportTriggered = false
}

// TriggerTick handles the tick event (e.g., play a tick sound)
func (t *TeleportTimer) TriggerTick() {
	// Placeholder: Should be integrated with the SoundManager to play tick sound
	println("Tick") // Replace with sound logic
}

// TeleportTriggered returns true if teleportation has been triggered
func (t *TeleportTimer) TeleportTriggered() bool {
	return t.teleportTriggered
}

// max utility function returns the maximum of two float32 values
func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
