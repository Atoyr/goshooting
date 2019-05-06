package entitys

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type EntityGroup struct {
	entity []Entity

	VirtualPosition engo.Point
	Size            float32
	Mergin          engo.Point
}

func (eg *EntityGroup) Move(vx, vy, speed float32) engo.Point {
	eg.VirtualPosition.X += vx
	eg.VirtualPosition.Y += vy
	for i := range eg.entity {
		eg.entity[i].Move(vx, vy, speed)
	}
	return eg.VirtualPosition
}

// func (eg *EntityGroup) Attack(playervx, playervy, speed, angle float32) {
// 	for i := range eg.entityModels {
// 		eg.entityModels[i].AttackFunc(playervx, playervy, speed, angle)
// 	}
// }
func (eg *EntityGroup) AddedRenderSystem(rs *common.RenderSystem) {
	for i := range eg.entity {
		eg.entity[i].AddedRenderSystem(rs)
	}
}
func (eg *EntityGroup) GetId() uint64 {
	return 0
}
func (eg *EntityGroup) GetVPoint() engo.Point {
	return eg.VirtualPosition
}
func (eg *EntityGroup) GetPoint() engo.Point {
	p := engo.Point{X: 0, Y: 0}
	if len(eg.entity) > 0 {
		p = eg.entity[0].spaceComponent.Position
	}
	return p
}
