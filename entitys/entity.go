package entitys

import (
	"fmt"
	"image/color"
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

func NewEntityBuilder(rc *engoCommon.RenderComponent) EntityBuilder {
	e := Entity{
		basicEntity:     ecs.NewBasic(),
		renderComponent: rc,

		virtualPosition:                 &engo.Point{X: 0, Y: 0},
		collisionDetectionRelativePoint: &engo.Point{X: 0, Y: 0},
		collisionDetectionSize:          0,
		mergin:                          &engo.Point{X: rc.Drawable.Width() * rc.Scale.X / 2, Y: rc.Drawable.Height() * rc.Scale.Y / 2},
		speed:                           0,
		angle:                           0,
		speedRate:                       0,
		angleRate:                       0,

		isRenderCollisionDetection: false,
	}

	e.spaceComponent = &engoCommon.SpaceComponent{
		Width:    rc.Drawable.Width() * rc.Scale.X,
		Height:   rc.Drawable.Height() * rc.Scale.Y,
		Rotation: 0,
	}
	e.SetVirtualPosition(e.GetVirtualPosition())

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

func (eb *EntityBuilder) BuildEntitySize(width, height float32) *EntityBuilder {
	eb.SetEntitySize(width, height)
	return eb
}

func (eb *EntityBuilder) BuildRotate(r float32) *EntityBuilder {
	eb.SetRotation(r)
	return eb
}

// Build is Building Entity
func (eb *EntityBuilder) Build() Entity {
	return *eb.Entity
}

// EntityInterface is entity func interface
type EntityInterface interface {
	Move(vx, vy, speed float32) engo.Point
	Update(frame int64)
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

	isRenderCollisionDetection bool
	collisionBasicEntity       ecs.BasicEntity
	collisionRenderComponent   *engoCommon.RenderComponent
	collisionSpaceComponent    *engoCommon.SpaceComponent

	Move                EntityMoveFunc
	Attack              EntityAttackFunc
	AddedRenderSystem   EntityAddedFunc
	RemovedRenderSystem EntityRemovedFunc
}

// AddedRenderSystem is added entitymodel at rendersystem
func (e *Entity) addedRenderSystem(rs *engoCommon.RenderSystem) {
	rs.Add(&e.basicEntity, e.renderComponent, e.spaceComponent)
	if e.isRenderCollisionDetection {
		fmt.Println("ADD IT")
		e.collisionBasicEntity = ecs.NewBasic()
		e.collisionRenderComponent.SetZIndex(1000)
		rs.Add(&e.collisionBasicEntity, e.collisionRenderComponent, e.collisionSpaceComponent)
	}
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
	return e.spaceComponent.Center()
}
func (e *Entity) GetVirtualPosition() engo.Point {
	return *e.virtualPosition
}

func (e *Entity) SetVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	e.virtualPosition = &engo.Point{X: point.X, Y: point.Y}
	e.spaceComponent.SetCenter(s.ConvertRenderPosition(point))
	if e.isRenderCollisionDetection {
		cpoint := engo.Point{X: 0, Y: 0}
		cpoint.Add(*e.virtualPosition)
		cpoint.Add(*e.collisionDetectionRelativePoint)
		e.collisionSpaceComponent.SetCenter(s.ConvertRenderPosition(cpoint))

	}
}

func (e *Entity) AddVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	e.virtualPosition.Add(point)
	p := engo.Point{X: e.virtualPosition.X, Y: e.virtualPosition.Y}
	e.spaceComponent.SetCenter(s.ConvertRenderPosition(p))
	if e.isRenderCollisionDetection {
		cpoint := engo.Point{X: 0, Y: 0}
		cpoint.Add(*e.virtualPosition)
		cpoint.Add(*e.collisionDetectionRelativePoint)
		e.collisionSpaceComponent.SetCenter(s.ConvertRenderPosition(cpoint))

	}
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

func (e *Entity) SetCollisionDetectionRelativePoint(p engo.Point) {
	e.collisionDetectionRelativePoint = &engo.Point{X: p.X, Y: p.Y}
}

func (e *Entity) SetCollisionDetectionSize(s float32) {
	e.collisionDetectionSize = s
}

func (e *Entity) SetZIndex(index float32) {
	e.renderComponent.SetZIndex(index)
}

func (e *Entity) SetRotation(r float32) {
	center := e.spaceComponent.Center()
	fmt.Println(e.spaceComponent.Center())
	e.spaceComponent.Rotation = r
	fmt.Println(e.spaceComponent.Center())
	e.spaceComponent.SetCenter(center)
}

func (e *Entity) SetEntitySize(width, height float32) {
	e.spaceComponent.Width = width
	e.spaceComponent.Height = height
}

func (e *Entity) IsCollision(target *Entity) bool {
	point := engo.Point{X: e.GetVirtualPosition().X, Y: e.GetVirtualPosition().Y}
	point.Add(*e.collisionDetectionRelativePoint)
	targetPoint := engo.Point{X: target.GetVirtualPosition().X, Y: target.GetVirtualPosition().Y}
	targetPoint.Add(*target.collisionDetectionRelativePoint)
	collisionDetectionSize := e.collisionDetectionSize + target.collisionDetectionSize
	size2 := (point.X-targetPoint.X)*(point.X-targetPoint.X) + (point.Y-targetPoint.Y)*(point.Y-targetPoint.Y)
	return size2 <= collisionDetectionSize*collisionDetectionSize
}

func (e *Entity) RenderCollisionDetection(b bool) {
	s := common.NewSetting()
	e.isRenderCollisionDetection = b
	if b {
		bgcolor := color.RGBA{200, 200, 200, 255}
		borderColor := color.RGBA{0, 0, 0, 255}
		rect := engoCommon.Circle{BorderWidth: 1, BorderColor: borderColor}
		e.collisionRenderComponent = &engoCommon.RenderComponent{
			Drawable: rect,
		}
		e.collisionRenderComponent.SetZIndex(999)
		e.collisionRenderComponent.Color = bgcolor
		sc := engoCommon.SpaceComponent{}
		point := engo.Point{X: 0, Y: 0}
		point.Add(*e.virtualPosition)
		point.Add(*e.collisionDetectionRelativePoint)
		sc.SetCenter(s.ConvertRenderPosition(point))
		sc.Width = e.collisionDetectionSize
		sc.Height = e.collisionDetectionSize
		e.collisionSpaceComponent = &sc
		e.AddedRenderSystem = e.addedRenderSystem
	}
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
