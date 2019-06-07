package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type EnemyBuilder struct {
	*Entity
}

func NewEnemyBuilder() EnemyBuilder {
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

	return EnemyBuilder{&e}
}

func (eb *EnemyBuilder) SetCollisionDetectionRelatevePoint(point engo.Point) {
	eb.collisionDetectionRelativePoint.Set(point.X, point.Y)
}

func (eb *EnemyBuilder) SetCollisionDetectionSize(size float32) {
	eb.collisionDetectionSize = size
}

func (eb *EnemyBuilder) SetSpeed(speed float32) {
	eb.Speed = speed
}

func (eb *EnemyBuilder) SetAngle(angle float32) {
	eb.Angle = angle
}

func (eb *EnemyBuilder) SetSpeedRate(speedrate float32) {
	eb.SpeedRate = speedrate
}

func (eb *EnemyBuilder) SetAngleRate(anglerate float32) {
	eb.AngleRate = anglerate
}

func (eb *EnemyBuilder) Build() Entity {
	e := eb.Entity.Clone()
	e.basicEntity = ecs.NewBasic()
	return *e
}
