package entitys

import (
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type PlayerBuilder struct {
	Entity *Entity
}

func NewPlayerBuilder() PlayerBuilder {
	s := common.NewSetting()
	sc := engoCommon.SpaceComponent{Position: engo.Point{X: 0, Y: 0}}
	rc := engoCommon.RenderComponent{Scale: s.Scale()}
	model := EntityModel{
		spaceComponent:  &sc,
		renderComponent: &rc,
		virtualPosition: &engo.Point{X: 0, Y: 0},
		scale:           0.5,
	}
	model.renderComponent.Scale.MultiplyScalar(model.scale)

	move := EntityMove{}
	e := Entity{EntityModel: &model, EntityMove: &move}
	e.SetVirtualPosition(engo.Point{X: 0, Y: 0})
	return PlayerBuilder{&e}
}

func (pb *PlayerBuilder) SetDrawable(drawable engoCommon.Drawable) {
	pb.Entity.SetDrawable(drawable)
}

func (pb *PlayerBuilder) SetEntitySize(width, height float32) {
	pb.Entity.SetEntitySize(width, height)
}

func (pb *PlayerBuilder) SetZIndex(index float32) {
	pb.Entity.SetZIndex(index)
}

func (pb *PlayerBuilder) SetVirtualPosition(point engo.Point) {
	pb.Entity.SetVirtualPosition(point)
}

func (pb *PlayerBuilder) SetCollisionDetectionRelatevePoint(point engo.Point) {
	pb.Entity.CollisionDetectionRelativePoint.Set(point.X, point.Y)
}

func (pb *PlayerBuilder) SetCollisionDetectionSize(size float32) {
	pb.Entity.CollisionDetectionSize = size
}

func (pb *PlayerBuilder) SetMoveFunc(movefunc EntityMoveFunc) {
	pb.Entity.Move = movefunc
}

func (pb *PlayerBuilder) SetSpeed(speed float32) {
	pb.Entity.Speed = speed
}

func (pb *PlayerBuilder) SetAngle(angle float32) {
	pb.Entity.Angle = angle
}

func (pb *PlayerBuilder) SetSpeedRate(speedrate float32) {
	pb.Entity.SpeedRate = speedrate
}

func (pb *PlayerBuilder) SetAngleRate(anglerate float32) {
	pb.Entity.AngleRate = anglerate
}

func (pb *PlayerBuilder) Build() Entity {
	e := *pb.Entity
	e.basicEntity = ecs.NewBasic()

	moveFunc := func(entity *Entity, vx, vy float32) {
		if vx == 0 && vy == 0 {
			return
		}
		s := common.NewSetting()
		x := entity.virtualPosition.X
		y := entity.virtualPosition.Y

		speed := float32(entity.Speed) / float32(math.Sqrt(float64(vx*vx+vy*vy)))
		x += speed * vx
		y += speed * vy

		gameArea := s.GetGameAreaSize()
		min := gameArea
		min.MultiplyScalar(-0.5)
		max := gameArea
		max.MultiplyScalar(0.5)
		mergin := engo.Point{X: entity.renderComponent.Drawable.Width(), Y: entity.renderComponent.Drawable.Height()}
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

		entity.SetVirtualPosition(engo.Point{X: x, Y: y})
	}
	e.Move = moveFunc

	return e
}
