package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type PlayerBuilder struct {
	*Entity
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

	move := new(EntityMove)
	attack := new(EntityAttack)
	collision := new(EntityCollision)
	e := Entity{EntityModel: &model, EntityMove: move, EntityAttack: attack, EntityCollision: collision}
	e.SetVirtualPosition(engo.Point{X: 0, Y: 0})
	return PlayerBuilder{&e}
}

func (pb *PlayerBuilder) SetCollisionDetectionRelatevePoint(point engo.Point) {
	pb.collisionDetectionRelativePoint.Set(point.X, point.Y)
}

func (pb *PlayerBuilder) SetCollisionDetectionSize(size float32) {
	pb.collisionDetectionSize = size
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
