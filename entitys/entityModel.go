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
	AppendChild(child *ecs.BasicEntity)

	AddedRenderSystem(rs *engoCommon.RenderSystem)
	RemovedRenderSystem(rs *engoCommon.RenderSystem) uint64

	Hidden() bool
	SetHidden(hidden bool)

	IsCollision(target Modeler) bool

	Position() engo.Point
	SetPosition(point engo.Point)
	AddPosition(point engo.Point)

	HitPoint() int
	SetHitPoint(hp int)
	AddHitPoint(hp int)

	Width() float32
	Height() float32
}

// EntityModel is Entity Base
type EntityModel struct {
	basicEntity     ecs.BasicEntity
	renderComponent engoCommon.RenderComponent
	spaceComponent  engoCommon.SpaceComponent
	virtualPosition engo.Point // center of entity
	scale           float32
	hitPoint        int
}

// ID is return BasicEntity.ID()
func (em *EntityModel) ID() uint64 {
	return em.basicEntity.ID()
}

func (em *EntityModel) AppendChild(child *ecs.BasicEntity) {
	em.basicEntity.AppendChild(child)
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

func (em *EntityModel) Hidden() bool {
	return em.renderComponent.Hidden
}

// SetHidden is Entity hiddened
func (em *EntityModel) SetHidden(hidden bool) {
	em.renderComponent.Hidden = hidden
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

// AddHitPoint Add hitpoint
func (em *EntityModel) HitPoint() int {
	return em.hitPoint
}

// SetHitPoint Set hitpoint
func (em *EntityModel) SetHitPoint(hp int) {
	em.hitPoint = hp
}

// AddHitPoint Add hitpoint
func (em *EntityModel) AddHitPoint(hp int) {
	em.hitPoint += hp
}

func (em *EntityModel) Width() float32 {
	return em.spaceComponent.Width
}

func (em *EntityModel) Height() float32 {
	return em.spaceComponent.Height
}

func (em *EntityModel) BasicEntity() ecs.BasicEntity {
	return em.basicEntity
}

func (em *EntityModel) SetDrawable(drawable engoCommon.Drawable) {
	em.renderComponent.Drawable = drawable
	em.spaceComponent.Width = drawable.Width() * em.renderComponent.Scale.X
	em.spaceComponent.Height = drawable.Height() * em.renderComponent.Scale.X
	em.spaceComponent.Rotation = 0
}

func (em *EntityModel) SetScale(scale float32) {
	em.scale = scale
	em.UpdateScale()
}

func (em *EntityModel) UpdateScale() {
	s := common.NewSetting()
	baseScale := s.Scale()
	baseScale.MultiplyScalar(em.scale)
	em.renderComponent.Scale = baseScale
}

func (em *EntityModel) SetEntitySize(width, height float32) {
	em.spaceComponent.Width = width
	em.spaceComponent.Height = height
}

func (em *EntityModel) SetZIndex(index float32) {
	em.renderComponent.SetZIndex(index)
}

func (em *EntityModel) updateSpaceComponentCenterPosition() {
	s := common.NewSetting()
	em.spaceComponent.SetCenter(s.ConvertVirtualPositionToRenderPosition(em.virtualPosition))
}
