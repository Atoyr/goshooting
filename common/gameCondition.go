package common

import (
	"sync"

	"github.com/EngoEngine/engo"
)

// GameCondition is game status and mode
type GameCondition struct {
	Frame                 uint64
	PlayerVirtualPosition engo.Point
}

var (
	gameCondition     *GameCondition
	gameConditionOnce sync.Once
)

// NewGameCondition get singleton GameCondition
func NewGameCondition() *GameCondition {
	gameConditionOnce.Do(func() {
		gameCondition = &GameCondition{
			Frame:                 0,
			PlayerVirtualPosition: engo.Point{X: 0, Y: 0},
		}
	})
	return gameCondition
}
