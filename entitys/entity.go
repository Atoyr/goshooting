package entitys

import (
	"fmt"
	"image/color"
	"math"

	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
	"github.com/jinzhu/copier"
)

// Builder is Entity Build Interface
type Builder interface {
	Build() Modeler
}

// Mover is Entity Move Interface
type Mover interface {
	Move(vx, vy, speed float32)
}

// Attacker is Entity Attacking interface
type Attacker interface {
	Attack(entity *Entity, frame float32)
}

// EntityAttackFunc is called entity.Attack()
type EntityAttackFunc func(entity *Entity, frame float32)

// Entity is GameAreaEntityObject
type Entity struct {
	*EntityModel

	Attack           EntityAttackFunc
	AttackStartFrame float32
	AttackFrame      float32
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
