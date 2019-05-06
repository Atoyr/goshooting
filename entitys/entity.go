package entitys

import (
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

// EntityBuilder is Build Entity Interface
type EntityBuilder struct {
	*Entity
}

func NewEntityBuilder(rc *engoCommon.RenderComponent, sc *engoCommon.SpaceComponent) EntityBuilder {
	e := Entity{
		basicEntity:     ecs.NewBasic(),
		renderComponent: rc,
		spaceComponent:  sc,

		virtualPosition:                 &engo.Point{X: 0, Y: 0},
		entitySize:                      &engo.Point{X: rc.Drawable.Width() * rc.Scale.X, Y: rc.Drawable.Height() * rc.Scale.Y},
		collisionDetectionRelativePoint: &engo.Point{X: 0, Y: 0},
		collisionDetectionSize:          0,
		mergin:                          &engo.Point{X: rc.Drawable.Width() * rc.Scale.X / 2, Y: rc.Drawable.Height() * rc.Scale.Y / 2},
		speed:                           0,
		angle:                           0,
		speedRate:                       0,
		angleRate:                       0,
	}
	e.Move = e.EntityMove
	e.Attack = func(playervx, playervy, speed, angle float32) {}
	e.AddedRenderSystem = e.addedRenderSystem
	e.RemovedRenderSystem = e.removedRenderSystem

	return EntityBuilder{&e}
}

func (eb *EntityBuilder) BuildVirtualPosition(p engo.Point) *EntityBuilder {
	eb.SetVirtualPosition(p)
	return eb
}

func (eb *EntityBuilder) BuildSpeed(s float32) *EntityBuilder {
	eb.SetSpeed(s)
	return eb
}

func (eb *EntityBuilder) BuildAngle(a float32) *EntityBuilder {
	eb.SetAngle(a)
	return eb
}

func (eb *EntityBuilder) BuildSpeedRate(s float32) *EntityBuilder {
	eb.SetSpeedRate(s)
	return eb
}

func (eb *EntityBuilder) BuildAngleRate(a float32) *EntityBuilder {
	eb.SetAngleRate(a)
	return eb
}

func (eb *EntityBuilder) BuildZIndex(index float32) *EntityBuilder {
	eb.SetZIndex(index)
	return eb
}

// Build is Building Entity
func (eb *EntityBuilder) Build() Entity {
	return *eb.Entity
}

// EntityInterface is entity func interface
type EntityInterface interface {
	Move(vx, vy, speed float32) engo.Point
	AddedRenderSystem(rs *engoCommon.RenderSystem)
	RemovedRenderSystem(rs *engoCommon.RenderSystem)

	GetId() uint64
	GetVertialPosition() engo.Point
	SetVirtualPosition()
	GetSpeed() float32
	SetSpeed()
	GetAngle() float32
	SetAngle()
	GetSpeedRate() float32
	SetSpeedRate()
	GetAngleRate() float32
	SetAngleRate()
	GetPoint() engo.Point
}

// EntityMoveFunc is called entity.Move()
type EntityMoveFunc func(vx, vy, speed float32)

// EntityAttackFunc is called entity.Attack()
type EntityAttackFunc func(playervx, playervy, speed, angle float32)

type EntityAddedFunc func(rs *engoCommon.RenderSystem)
type EntityRemovedFunc func(rs *engoCommon.RenderSystem) uint64

// Entity is GameAreaEntityObject
type Entity struct {
	basicEntity     ecs.BasicEntity
	renderComponent *engoCommon.RenderComponent
	spaceComponent  *engoCommon.SpaceComponent

	virtualPosition                 *engo.Point // center of entity
	entitySize                      *engo.Point // entity Size
	collisionDetectionRelativePoint *engo.Point // Collision Detection Position from relative virtualPosition
	collisionDetectionSize          float32     // Collision Detection Size circle
	mergin                          *engo.Point // entity mergin
	speed                           float32
	angle                           float32
	speedRate                       float32
	angleRate                       float32

	Move                EntityMoveFunc
	Attack              EntityAttackFunc
	AddedRenderSystem   EntityAddedFunc
	RemovedRenderSystem EntityRemovedFunc
}

// AddedRenderSystem is added entitymodel at rendersystem
func (e *Entity) addedRenderSystem(rs *engoCommon.RenderSystem) {
	rs.Add(&e.basicEntity, e.renderComponent, e.spaceComponent)
}

func (e *Entity) removedRenderSystem(rs *engoCommon.RenderSystem) uint64 {
	i := e.basicEntity.ID()
	rs.Remove(e.basicEntity)
	return i
}

// GetID is return BasicEntity.ID()
func (e *Entity) GetID() uint64 {
	return e.basicEntity.ID()
}

func (e *Entity) GetBasicEntity() ecs.BasicEntity {
	return e.basicEntity
}

func (e *Entity) GetPoint() engo.Point {
	return e.spaceComponent.Position
}
func (e *Entity) GetVirtualPosition() engo.Point {
	return *e.virtualPosition
}

func (e *Entity) SetVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	e.virtualPosition = &engo.Point{X: point.X, Y: point.Y}
	e.spaceComponent.Position = s.ConvertRenderPosition(*point.Subtract(*e.mergin))
}

func (e *Entity) AddVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	e.virtualPosition.Add(point)
	p := engo.Point{X: e.virtualPosition.X, Y: e.virtualPosition.Y}
	e.spaceComponent.Position = s.ConvertRenderPosition(*p.Subtract(*e.mergin))
}

func (e *Entity) GetMergin() engo.Point {
	return *e.mergin
}

func (e *Entity) GetSpeed() float32 {
	return e.speed
}

func (e *Entity) SetSpeed(s float32) {
	e.speed = s
}

func (e *Entity) GetAngle() float32 {
	return e.angle
}

func (e *Entity) SetAngle(a float32) {
	e.angle = a
}

func (e *Entity) GetSpeedRate() float32 {
	return e.speedRate
}

func (e *Entity) SetSpeedRate(s float32) {
	e.speedRate = s
}

func (e *Entity) GetAngleRate() float32 {
	return e.angleRate
}

func (e *Entity) SetAngleRate(a float32) {
	e.angleRate = a
}

func (e *Entity) SetZIndex(index float32) {
	e.renderComponent.SetZIndex(index)
}

func (e *Entity) EntityMove(vx, vy, speed float32) {
	p := engo.Point{X: vx * speed, Y: vy * speed}
	e.AddVirtualPosition(p)
}

func (e *Entity) EntityMoveForPlayer(vx, vy, speed float32) {
	s := common.NewSetting()
	x := e.virtualPosition.X
	y := e.virtualPosition.Y

	if vx != 0 && vy != 0 {
		speed = float32(speed) / 1.414
	}
	x += speed * vx
	y += speed * vy
	max := s.GetGameAreaSize()
	if x < e.mergin.X {
		x = e.mergin.X
	} else if x > max.X-e.mergin.X {
		x = max.X - e.mergin.X
	}
	if y < e.mergin.Y {
		y = e.mergin.Y
	} else if y > max.Y-e.mergin.Y {
		y = max.Y - e.mergin.Y
	}
	e.SetVirtualPosition(engo.Point{X: x, Y: y})
}

func (e *Entity) GetMoveInfo() (vx, vy, speed float32) {
	rad := float64((e.angle - 90) / float32(180) * math.Pi)
	vx = float32(math.Cos(rad))
	vy = float32(math.Sin(rad))
	speed = e.speed
	e.angle += e.angleRate
	e.speed += e.speedRate

	return vx, vy, speed
}

func (e *Entity) GetPlayerMoveInfo(isleft, isright, isup, isdown, islowspeed bool) (vx, vy, speed float32) {
	vx = 0
	vy = 0
	if islowspeed {
		speed = float32(e.speed / 2)
	} else {
		speed = float32(e.speed)
	}
	if isleft && !isright {
		vx = -1
	} else if !isleft && isright {
		vx = 1
	}
	if isup && !isdown {
		vy = -1
	} else if !isup && isdown {
		vy = 1
	}
	return vx, vy, speed
}
