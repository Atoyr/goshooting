package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
	"github.com/jinzhu/copier"
	"math"
)

type Bullet struct {
	*EntityModel

	Speed        float32
	SpeedRate    float32
	RotationRate float32
}

func (b *Bullet) Move() {
	rad := float64((b.Rotation() - 90) / float32(180) * math.Pi)
	vx := float32(math.Cos(rad))
	vy := float32(math.Sin(rad))

	if vx == 0 && vy == 0 {
		return
	}

	vector := engo.Point{X: vx, Y: vy}
	speed := float32(b.Speed) / float32(math.Sqrt(float64(vx*vx+vy*vy)))
	vector.MultiplyScalar(speed)

	b.AddPosition(vector)
	b.AddRotation(b.RotationRate)
	b.Speed += b.SpeedRate
}

type BulletBuilder struct {
	*Bullet
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
		hitPoint:        0,
	}
	model.renderComponent.Scale.MultiplyScalar(model.scale)

	model.SetPosition(engo.Point{X: 0, Y: 0})

	bullet := new(Bullet)

	bullet.EntityModel = &model
	bullet.Speed = 0
	bullet.SpeedRate = 0
	bullet.RotationRate = 0
	return BulletBuilder{bullet}
}

func (bb *BulletBuilder) Build() Modeler {
	entityModel := new(EntityModel)
	bullet := new(Bullet)
	copier.Copy(&entityModel, bb.EntityModel)
	copier.Copy(&bullet, bb.Bullet)
	bullet.EntityModel = entityModel
	bullet.basicEntity = ecs.NewBasic()

	return *bullet
}

func (bb *BulletBuilder) Clone() Builder {
	builder := new(BulletBuilder)
	entityModel := new(EntityModel)
	bullet := new(Bullet)
	copier.Copy(&entityModel, bb.EntityModel)
	copier.Copy(&bullet, bb.Bullet)
	bullet.EntityModel = entityModel
	builder.Bullet = bullet

	return builder
}
