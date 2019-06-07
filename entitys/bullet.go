package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type BulletBuilder struct {
	*Entity
}

func NewBulletBuilder() BulletBuilder {
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
	return BulletBuilder{&e}
}

func (bb *BulletBuilder) SetCollisionDetectionRelatevePoint(point engo.Point) {
	bb.Entity.collisionDetectionRelativePoint.Set(point.X, point.Y)
}

func (bb *BulletBuilder) SetCollisionDetectionSize(size float32) {
	bb.Entity.collisionDetectionSize = size
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
	e := bb.Entity.Clone()
	e.basicEntity = ecs.NewBasic()

	return *e
}
