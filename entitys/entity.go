package entitys

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type EntityBuilder interface {
	Build() Entity
}

type EntityModeler interface {
	ID() uint64
	BasicEntity() ecs.BasicEntity
	Point() engo.Point
	SetDrawable(drawable engoCommon.Drawable)
	SetEntitySize(width, height float32)
	SetZIndex(index float32)
	SetVirtualPosition(point engo.Point)
	SetHidden(b bool)
	AddVirtualPosition(point engo.Point)
	VertualPosition() engo.Point
	IsCollision(target Entity) bool
	RenderCollisionDetection(b bool)
	Mergin() engo.Point
	Hidden() bool

	AddedRenderSystem(rs *engoCommon.RenderSystem)
	RemovedRenderSystem(rs *engoCommon.RenderSystem) uint64
}

type EntityModel struct {
	basicEntity     ecs.BasicEntity
	renderComponent *engoCommon.RenderComponent
	spaceComponent  *engoCommon.SpaceComponent
	virtualPosition *engo.Point // center of entity
	scale           float32

	CreateFrame float32

	CollisionDetectionRelativePoint *engo.Point // Collision Detection Position from relative virtualPosition
	CollisionDetectionSize          float32     // Collision Detection Size circle
}

type EntityMove struct {
	Move      EntityMoveFunc
	MoveInfo  EntityMoveInfoFunc
	Speed     float32
	Angle     float32
	SpeedRate float32
	AngleRate float32
}

type EntityAttack struct {
	Attack           EntityAttackFunc
	AttackStartFrame float32
	AttackFrame      float32
}

type EntityAttacker interface {
	Attack(entity *Entity, frame float32)
}

// EntityMoveFunc is called entity.Move()
type EntityMoveFunc func(entity *Entity, vx, vy float32)
type EntityMoveInfoFunc func(entity *Entity, frame float32) (vx, vy float32)

// EntityAttackFunc is called entity.Attack()
type EntityAttackFunc func(entity *Entity, frame float32)

type EntityAddedFunc func(rs *engoCommon.RenderSystem)
type EntityRemovedFunc func(rs *engoCommon.RenderSystem) uint64

// Entity is GameAreaEntityObject
type Entity struct {
	*EntityModel
	*EntityMove
	*EntityAttack
}

// GetID is return BasicEntity.ID()
func (e *Entity) ID() uint64 {
	return e.basicEntity.ID()
}

func (e *Entity) BasicEntity() ecs.BasicEntity {
	return e.basicEntity
}

func (e *Entity) Point() engo.Point {
	return e.spaceComponent.Center()
}

func (e *Entity) SetDrawable(drawable engoCommon.Drawable) {
	e.renderComponent.Drawable = drawable
	e.spaceComponent.Width = drawable.Width() * e.renderComponent.Scale.X
	e.spaceComponent.Height = drawable.Height() * e.renderComponent.Scale.X
	e.spaceComponent.Rotation = 0
}

func (e *Entity) SetScale(scale float32) {
	e.scale = scale
	e.UpdateScale()
}

func (e *Entity) UpdateScale() {
	s := common.NewSetting()
	baseScale := s.Scale()
	baseScale.MultiplyScalar(e.scale)
	e.renderComponent.Scale = baseScale
}

func (e *Entity) SetEntitySize(width, height float32) {
	e.spaceComponent.Width = width
	e.spaceComponent.Height = height
}

func (e *Entity) SetZIndex(index float32) {
	e.renderComponent.SetZIndex(index)
}

// RenderCollisionDetection
func (e *Entity) RenderCollisionDetection(b bool) {

}

// AddedRenderSystem is added entitymodel at rendersystem
func (e *Entity) AddedRenderSystem(rs *engoCommon.RenderSystem) {
	rs.Add(&e.basicEntity, e.renderComponent, e.spaceComponent)
}

func (e *Entity) RemovedRenderSystem(rs *engoCommon.RenderSystem) uint64 {
	i := e.basicEntity.ID()
	rs.Remove(e.basicEntity)
	return i
}

func (e *Entity) SetVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	e.EntityModel.virtualPosition = &engo.Point{X: point.X, Y: point.Y}
	e.spaceComponent.SetCenter(s.ConvertVirtualPositionToPhysicsPosition(*e.virtualPosition))
}

// SetHidden is Entity hiddened
func (e *Entity) SetHidden(b bool) {
	e.renderComponent.Hidden = b
}

func (e *Entity) AddVirtualPosition(point engo.Point) {
	s := common.NewSetting()
	e.virtualPosition.Add(point)
	p := engo.Point{X: e.virtualPosition.X, Y: e.virtualPosition.Y}
	e.spaceComponent.SetCenter(s.ConvertVirtualPositionToPhysicsPosition(p))
	//	if e.isRenderCollisionDetection {
	//		cpoint := engo.Point{X: 0, Y: 0}
	//		cpoint.Add(*e.virtualPosition)
	//		cpoint.Add(*e.collisionDetectionRelativePoint)
	//		e.collisionSpaceComponent.SetCenter(s.ConvertVirtualPositionToPhysicsPosition(cpoint))
	//
	//	}
}

func (e *Entity) VirtualPosition() engo.Point {
	return *e.virtualPosition
}

func (e *Entity) Mergin() engo.Point {
	return engo.Point{X: e.renderComponent.Drawable.Width() * e.renderComponent.Scale.X, Y: e.renderComponent.Drawable.Height() * e.renderComponent.Scale.Y}
}

func (e *Entity) Hidden() bool {
	return e.renderComponent.Hidden
}

func (e *Entity) IsCollision(target Entity) bool {
	point := engo.Point{X: e.VirtualPosition().X, Y: e.VirtualPosition().Y}
	point.Add(*e.CollisionDetectionRelativePoint)
	targetpoint := engo.Point{X: target.VirtualPosition().X, Y: target.VirtualPosition().Y}
	targetpoint.Add(*target.CollisionDetectionRelativePoint)
	collisionDetectionSize := e.CollisionDetectionSize + target.CollisionDetectionSize
	return point.PointDistanceSquared(targetpoint) <= collisionDetectionSize*collisionDetectionSize
}

func (e *Entity) Update(frame float32) {
	var vx, vy float32
	if e.Move != nil {
		if e.MoveInfo != nil {
			vx, vy = e.MoveInfo(e, frame)
		}
		e.Move(e, vx, vy)
	}

	if e.Attack != nil {
		e.Attack(e, frame)
	}

}

func (e *Entity) Clone() *Entity {
	entityModel := *e.EntityModel
	entityMove := *e.EntityMove
	entityAttack := *e.EntityAttack
	entity := new(Entity)
	entity.EntityModel = &entityModel
	entity.EntityMove = &entityMove
	entity.EntityAttack = &entityAttack
	return entity
}

func (e *Entity) String() string {
	return fmt.Sprintf("%#v %#v", e.EntityModel, e.EntityMove)
}

// func (eb *entityModelBuilder) BuildZIndex(index float32) *entityModelBuilder {
// 	eb.renderComponent.SetZIndex(index)
// 	return eb
// }
//
// func (eb *entityModelBuilder) buildEntityModel() EntityModel {
// 	b := ecs.NewBasic()
// 	rc := &eb.renderComponent
// 	sc := &eb.spaceComponent
// 	e := EntityModel{
// 		basicEntity:     b,
// 		renderComponent: *rc,
// 		spaceComponent:  *sc,
// 	}
// 	return e
// }

// func (e *Entity) RenderCollisionDetection(b bool) {
// 	s := common.NewSetting()
// 	e.isRenderCollisionDetection = b
// 	if b {
// 		bgcolor := color.RGBA{200, 200, 200, 255}
// 		borderColor := color.RGBA{0, 0, 0, 255}
// 		rect := engoCommon.Circle{BorderWidth: 1, BorderColor: borderColor}
// 		e.collisionRenderComponent = &engoCommon.RenderComponent{
// 			Drawable: rect,
// 		}
// 		e.collisionRenderComponent.SetZIndex(999)
// 		e.collisionRenderComponent.Color = bgcolor
// 		sc := engoCommon.SpaceComponent{}
// 		point := engo.Point{X: 0, Y: 0}
// 		point.Add(*e.virtualPosition)
// 		point.Add(*e.collisionDetectionRelativePoint)
// 		sc.SetCenter(s.ConvertVirtualPositionToPhysicsPosition(point))
// 		sc.Width = e.collisionDetectionSize
// 		sc.Height = e.collisionDetectionSize
// 		e.collisionSpaceComponent = &sc
// 		e.AddedRenderSystem = e.addedRenderSystem
// 	}
// }
//
// func (e *Entity) EntityMove(vx, vy, speed float32) {
// 	p := engo.Point{X: vx * speed, Y: vy * speed}
// 	e.AddVirtualPosition(p)
// }
//
// func (e *Entity) EntityMoveForPlayer(vx, vy, speed float32) {
// 	s := common.NewSetting()
// 	x := e.virtualPosition.X
// 	y := e.virtualPosition.Y
//
// 	if vx != 0 && vy != 0 {
// 		speed = float32(speed) / 1.414
// 	}
// 	x += speed * vx
// 	y += speed * vy
// 	max := s.GetGameAreaSize()
// 	if x < e.mergin.X {
// 		x = e.mergin.X
// 	} else if x > max.X-e.mergin.X {
// 		x = max.X - e.mergin.X
// 	}
// 	if y < e.mergin.Y {
// 		y = e.mergin.Y
// 	} else if y > max.Y-e.mergin.Y {
// 		y = max.Y - e.mergin.Y
// 	}
// 	e.SetVirtualPosition(engo.Point{X: x, Y: y})
// }
//
// func (e *Entity) GetMoveInfo() (vx, vy, speed float32) {
// 	rad := float64((e.angle - 90) / float32(180) * math.Pi)
// 	vx = float32(math.Cos(rad))
// 	vy = float32(math.Sin(rad))
// 	speed = e.speed
// 	e.angle += e.angleRate
// 	e.speed += e.speedRate
//
// 	return vx, vy, speed
// }
//
// func (e *Entity) GetPlayerMoveInfo(isleft, isright, isup, isdown, islowspeed bool) (vx, vy, speed float32) {
// 	vx = 0
// 	vy = 0
// 	if islowspeed {
// 		speed = float32(e.speed / 2)
// 	} else {
// 		speed = float32(e.speed)
// 	}
// 	if isleft && !isright {
// 		vx = -1
// 	} else if !isleft && isright {
// 		vx = 1
// 	}
// 	if isup && !isdown {
// 		vy = -1
// 	} else if !isup && isdown {
// 		vy = 1
// 	}
// 	return vx, vy, speed
// }
