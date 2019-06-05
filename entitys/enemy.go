package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type EnemyBuilder struct {
	Entity *Entity
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

	move := EntityMove{}
	attack := EntityAttack{}
	e := Entity{EntityModel: &model, EntityMove: &move, EntityAttack: &attack}
	e.SetVirtualPosition(engo.Point{X: 0, Y: 0})

	return EnemyBuilder{&e}
}

func (eb *EnemyBuilder) SetDrawable(drawable engoCommon.Drawable) {
	eb.Entity.SetDrawable(drawable)
}

func (eb *EnemyBuilder) SetEntitySize(width, height float32) {
	eb.Entity.SetEntitySize(width, height)
}

func (eb *EnemyBuilder) SetZIndex(index float32) {
	eb.Entity.SetZIndex(index)
}

func (eb *EnemyBuilder) SetVirtualPosition(point engo.Point) {
	eb.Entity.SetVirtualPosition(point)
}

func (eb *EnemyBuilder) SetCollisionDetectionRelatevePoint(point engo.Point) {
	eb.Entity.CollisionDetectionRelativePoint.Set(point.X, point.Y)
}

func (eb *EnemyBuilder) SetCollisionDetectionSize(size float32) {
	eb.Entity.CollisionDetectionSize = size
}

func (eb *EnemyBuilder) SetSpeed(speed float32) {
	eb.Entity.Speed = speed
}

func (eb *EnemyBuilder) SetAngle(angle float32) {
	eb.Entity.Angle = angle
}

func (eb *EnemyBuilder) SetSpeedRate(speedrate float32) {
	eb.Entity.SpeedRate = speedrate
}

func (eb *EnemyBuilder) SetAngleRate(anglerate float32) {
	eb.Entity.AngleRate = anglerate
}

func (eb *EnemyBuilder) Build() Entity {
	e := eb.Entity.Clone()
	e.basicEntity = ecs.NewBasic()
	return *e
}
