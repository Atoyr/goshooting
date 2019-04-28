package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type EntityBuilder interface {
	VirtualPosition(engo.Point) EntityBuilder
	Size(float32) EntityBuilder
	Mergin(float32) EntityBuilder
	Build() EntityModel
}

type EntityInterface interface {
	Move(vx, vy, speed, angle float32) engo.Point
	Attack(playervx, playervy, speed, angle float32)
	AddedRenderSystem(rs *common.RenderSystem)
	GetId() uint64
	GetVPoint() engo.Point
	GetPoint() engo.Point
}

type EntityMoveFunc func(vx, vy, speed, angle float32) engo.Point
type EntityAttackFunc func(playervx, playervy, speed, angle float32)

type EntityModel struct {
	BasicEntity     ecs.BasicEntity
	RenderComponent common.RenderComponent
	SpaceComponent  common.SpaceComponent
	VirtualPosition engo.Point
	Size            float32
	Mergin          engo.Point
	Frame           uint64
	MoveFunc        EntityMoveFunc
	AttackFunc      EntityAttackFunc
}

func (e *EntityModel) convertPosition() engo.Point {
	return engo.Point{X: e.VirtualPosition.X - e.Mergin.X, Y: e.VirtualPosition.Y - e.Mergin.Y}
}

func (e *EntityModel) Move(vx, vy, speed, angle float32) engo.Point {
	return e.MoveFunc(vx, vy, speed, angle)
}

func (e *EntityModel) Attack(playervx, playervy, speed, angle float32) {
	e.AttackFunc(playervx, playervy, speed, angle)
}

func (e *EntityModel) AddedRenderSystem(rs *common.RenderSystem) {
	rs.Add(&e.BasicEntity, &e.RenderComponent, &e.SpaceComponent)
}
func (e *EntityModel) GetId() uint64 {
	return e.BasicEntity.ID()
}
func (e *EntityModel) GetPoint() engo.Point {
	return e.SpaceComponent.Position
}
func (e *EntityModel) GetVPoint() engo.Point {
	return e.VirtualPosition
}
