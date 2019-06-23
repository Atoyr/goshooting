package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

// Modeler is EntityModel Interface
type Modeler interface {
	ID() uint64
	AddedRenderSystem(rs *engoCommon.RenderSystem)
	RemovedRenderSystem(rs *engoCommon.RenderSystem) uint64
	IsCollision(target Modeler) bool
	Position() engo.Point
	SetPosition(point engo.Point)
	AddPosition(point engo.Point)
	Size() float32
}

// EntityModel is Entity Base
type EntityModel struct {
	basicEntity     ecs.BasicEntity
	renderComponent engoCommon.RenderComponent
	spaceComponent  engoCommon.SpaceComponent
	virtualPosition engo.Point // center of entity
	size            float32
	scale           float32
	hitPoint        int32
}

// ID is return BasicEntity.ID()
func (em *EntityModel) ID() uint64 {
	return em.basicEntity.ID()
}

// AddedRenderSystem is added entitymodel at rendersystem
func (em *EntityModel) AddedRenderSystem(rs *engoCommon.RenderSystem) {
	rs.Add(&em.basicEntity, &em.renderComponent, &em.spaceComponent)
}

func (em *EntityModel) RemovedRenderSystem(rs *engoCommon.RenderSystem) uint64 {
	i := em.basicEntity.ID()
	rs.Remove(em.basicEntity)
	return i
}

func (em *EntityModel) IsCollision(target Modeler) bool {
	Position := em.Position()
	targetPosition := target.Position()
	collisionDetectionSize := em.Size() + target.Size()
	return Position.PointDistanceSquared(targetPosition) <= collisionDetectionSize*collisionDetectionSize
}

func (em *EntityModel) Position() engo.Point {
	return em.virtualPosition
}

// SetPosition is Set virtual Position and update spaceComponentPosition
func (em *EntityModel) SetPosition(point engo.Point) {
	em.virtualPosition = point
	em.updateSpaceComponentCenterPosition()
}

// SetPosition is Add virtual Position and update spaceComponentPosition
func (em *EntityModel) AddPosition(point engo.Point) {
	em.virtualPosition.Add(point)
	em.updateSpaceComponentCenterPosition()
}

func (em *EntityModel) Size() float32 {
	return em.size
}

func (em *EntityModel) BasicEntity() ecs.BasicEntity {
	return em.basicEntity
}

func (em *EntityModel) updateSpaceComponentCenterPosition() {
	s := common.NewSetting()
	em.spaceComponent.SetCenter(s.ConvertVirtualPositionToRenderPosition(em.virtualPosition))
}
