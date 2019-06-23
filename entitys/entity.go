package entitys

import (
	"fmt"
	"image/color"
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
	"github.com/jinzhu/copier"
)

// Builder is Entity Build Interface
type Builder interface {
	Build() Entity
}

// Mover is Entity Move Interface
type Mover interface {
	Move(vx, vy, speed float32)
	MoveInfo(frame float32) (vx, vy float32)
}

// Attacker is Entity Attacking interface
type Attacker interface {
	Attack(entity *Entity, frame float32)
}

// EntityMove is Moving on entity
type EntityMove struct {
	Speed     float32
	Angle     float32
	SpeedRate float32
	AngleRate float32
}

// EntityAttack is Attacking entity
type EntityAttack struct {
	Attack           EntityAttackFunc
	AttackStartFrame float32
	AttackFrame      float32
}

// EntityCollision is Judge Collision on Entity
type EntityCollision struct {
	isRenderCollision               bool
	collisionBasicEntity            ecs.BasicEntity
	collisionRenderComponent        engoCommon.RenderComponent
	collisionSpaceComponent         engoCommon.SpaceComponent
	collisionDetectionRelativePoint engo.Point // Collision Detection Position from relative virtualPosition
	collisionDetectionSize          float32    // Collision Detection Size circle
}

// EntityAttackFunc is called entity.Attack()
type EntityAttackFunc func(entity *Entity, frame float32)

// Entity is GameAreaEntityObject
type Entity struct {
	*EntityModel
	*EntityMove
	*EntityAttack
	*EntityCollision
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

// SetHidden is Entity hiddened
func (e *Entity) SetHidden(b bool) {
	e.renderComponent.Hidden = b
	if e.isRenderCollision {
		e.collisionRenderComponent.Hidden = true
	}
}

// SetHitPoint Set hitpoint
func (e *Entity) SetHitPoint(hp int32) {
	e.hitPoint = hp
}

// AddHitPoint Add hitpoint
func (e *Entity) AddHitPoint(hp int32) {
	e.hitPoint += hp
}

// AddHitPoint Add hitpoint
func (e *Entity) HitPoint() int32 {
	return e.hitPoint
}

func (e *Entity) Mergin() engo.Point {
	return engo.Point{X: e.renderComponent.Drawable.Width() * e.renderComponent.Scale.X, Y: e.renderComponent.Drawable.Height() * e.renderComponent.Scale.Y}
}

func (e *Entity) Hidden() bool {
	return e.renderComponent.Hidden
}

func (e *Entity) SetCollisionDetectionRelativePoint(point engo.Point) {
	e.collisionDetectionRelativePoint = point
}

func (e *Entity) updateCollisionSpaceComponentCenterPosition() {
	s := common.NewSetting()
	p := new(engo.Point)
	p.Add(e.virtualPosition)
	p.Add(e.collisionDetectionRelativePoint)
	e.collisionSpaceComponent.SetCenter(s.ConvertVirtualPositionToRenderPosition(*p))
}

func (e *Entity) SetCollisionBasicEntity(basic ecs.BasicEntity) {
	e.collisionBasicEntity = basic
}

func (e *Entity) IsRenderCollision() bool {
	return e.isRenderCollision
}

func (e *Entity) AddedRenderSystemToCollisionComponent(rs *engoCommon.RenderSystem) {
	rs.Add(&e.collisionBasicEntity, &e.collisionRenderComponent, &e.collisionSpaceComponent)
	e.isRenderCollision = true
}

func (e *Entity) RemovedRenderSystemToCollisionComponent(rs *engoCommon.RenderSystem) uint64 {
	i := e.collisionBasicEntity.ID()
	rs.Remove(e.collisionBasicEntity)
	e.isRenderCollision = false
	return i
}

func (e *Entity) MoveInfo(frame float32) (vx, vy float32) {
	rad := float64((e.Angle - 90) / float32(180) * math.Pi)
	vx = float32(math.Cos(rad))
	vy = float32(math.Sin(rad))
	return vx, vy
}

func (e *Entity) Move(vx, vy, speed float32) {
	if vx == 0 && vy == 0 {
		return
	}
	x := e.virtualPosition.X
	y := e.virtualPosition.Y

	speed = float32(speed) / float32(math.Sqrt(float64(vx*vx+vy*vy)))
	x += speed * vx
	y += speed * vy

	s := common.NewSetting()
	gameArea := s.GameAreaSize()
	min := gameArea
	min.MultiplyScalar(-0.5)
	max := gameArea
	max.MultiplyScalar(0.5)
	mergin := engo.Point{X: e.renderComponent.Drawable.Width(), Y: e.renderComponent.Drawable.Height()}
	mergin.Multiply(s.Scale())
	mergin.MultiplyScalar(e.scale * 0.5)

	if minX := min.X + mergin.X; x < minX {
		x = minX
	} else if maxX := max.X - mergin.X; x > maxX {
		x = maxX
	}
	if minY := min.Y + mergin.Y; y < minY {
		y = minY
	} else if maxY := max.Y - mergin.Y; y > maxY {
		y = maxY
	}

	e.SetPosition(engo.Point{X: x, Y: y})
	e.Angle += e.AngleRate
	e.SpeedRate += e.SpeedRate
}

func (e *Entity) Update(frame float32) {
	vx, vy := e.MoveInfo(frame)
	e.Move(vx, vy, e.Speed)

	if e.Attack != nil {
		e.Attack(e, frame)
	}
}

// Clone is Cloned Entity
func (e *Entity) Clone() *Entity {
	entityModel := new(EntityModel)
	entityMove := new(EntityMove)
	entityAttack := new(EntityAttack)
	entityCollision := new(EntityCollision)
	copier.Copy(&entityModel, e.EntityModel)
	copier.Copy(&entityMove, e.EntityMove)
	copier.Copy(&entityAttack, e.EntityAttack)
	copier.Copy(&entityCollision, e.EntityCollision)

	entity := new(Entity)
	entity.EntityModel = entityModel
	entity.EntityMove = entityMove
	entity.EntityAttack = entityAttack
	entity.EntityCollision = entityCollision
	return entity
}

func (e *Entity) String() string {
	return fmt.Sprintf("%#v %#v", e.EntityModel, e.EntityMove)
}

func (e *Entity) RenderCollisionDetection(b bool) {
	// s := common.NewSetting()
	if b {
		bgcolor := color.RGBA{200, 200, 200, 255}
		borderColor := color.RGBA{0, 0, 0, 255}
		rect := engoCommon.Circle{BorderWidth: 1, BorderColor: borderColor}
		e.collisionRenderComponent = engoCommon.RenderComponent{
			Drawable: rect,
		}
		e.collisionRenderComponent.SetZIndex(999)
		e.collisionRenderComponent.Color = bgcolor
		sc := engoCommon.SpaceComponent{}
		point := engo.Point{X: 0, Y: 0}
		point.Add(e.virtualPosition)
		point.Add(e.collisionDetectionRelativePoint)
		// sc.SetCenter(s.ConvertPositionToRenderPosition(point))
		sc.Width = e.collisionDetectionSize * 2
		sc.Height = e.collisionDetectionSize * 2
		e.collisionSpaceComponent = sc
	}
}
