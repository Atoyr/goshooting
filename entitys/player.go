package entitys

import (
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
		spaceComponent:  sc,
		renderComponent: rc,
		virtualPosition: engo.Point{X: 0, Y: 0},
		scale:           0.5,
	}
	model.renderComponent.Scale.MultiplyScalar(model.scale)

	move := EntityMove{}
	attack := EntityAttack{}
	e := Entity{EntityModel: &model, EntityMove: &move, EntityAttack: &attack}
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
	e := pb.Entity.Clone()
	e.basicEntity = ecs.NewBasic()
	return *e
}
