package entitys

import (
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type BulletBuilder struct {
	Entity *Entity
}

func NewBulletBuilder() BulletBuilder {
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
	return BulletBuilder{&e}
}

func (bb *BulletBuilder) SetDrawable(drawable engoCommon.Drawable) {
	bb.Entity.SetDrawable(drawable)
}

func (bb *BulletBuilder) SetEntitySize(width, height float32) {
	bb.Entity.SetEntitySize(width, height)
}

func (bb *BulletBuilder) SetZIndex(index float32) {
	bb.Entity.SetZIndex(index)
}

func (bb *BulletBuilder) SetVirtualPosition(point engo.Point) {
	bb.Entity.SetVirtualPosition(point)
}

func (bb *BulletBuilder) SetCollisionDetectionRelatevePoint(point engo.Point) {
	bb.Entity.CollisionDetectionRelativePoint.Set(point.X, point.Y)
}

func (bb *BulletBuilder) SetCollisionDetectionSize(size float32) {
	bb.Entity.CollisionDetectionSize = size
}

func (bb *BulletBuilder) SetMoveFunc(movefunc EntityMoveFunc) {
	bb.Entity.Move = movefunc
}

func (bb *BulletBuilder) SetSpeed(speed float32) {
	bb.Entity.Speed = speed
}

func (bb *BulletBuilder) SetAngle(angle float32) {
	bb.Entity.Angle = angle
}

func (bb *BulletBuilder) SetSpeedRate(speedrate float32) {
	bb.Entity.SpeedRate = speedrate
}

func (bb *BulletBuilder) SetAngleRate(anglerate float32) {
	bb.Entity.AngleRate = anglerate
}

func (bb *BulletBuilder) Build() Entity {
	e := *bb.Entity
	e.basicEntity = ecs.NewBasic()

	moveFunc := func(entity *Entity, vx, vy float32) {
		rad := float64((entity.Angle - 90) / float32(180) * math.Pi)
		vx = float32(math.Cos(rad))
		vy = float32(math.Sin(rad))
		entity.Angle += entity.AngleRate
		entity.SpeedRate += entity.SpeedRate
		//
		if vx == 0 && vy == 0 {
			return
		}
		x := entity.virtualPosition.X
		y := entity.virtualPosition.Y

		speed := float32(entity.Speed) / float32(math.Sqrt(float64(vx*vx+vy*vy)))
		x += speed * vx
		y += speed * vy
		entity.SetVirtualPosition(engo.Point{X: x, Y: y})
	}
	e.Move = moveFunc

	return e
}
