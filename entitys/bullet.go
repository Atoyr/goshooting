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

	Speed     float32
	Angle     float32
	SpeedRate float32
	AngleRate float32
}

func (b *Bullet) SetSpeed(speed float32) {
	b.Speed = speed
}

func (b *Bullet) SetAngle(angle float32) {
	b.Angle = angle
}

func (b *Bullet) SetSpeedRate(speedrate float32) {
	b.SpeedRate = speedrate
}

func (b *Bullet) SetAngleRate(anglerate float32) {
	b.AngleRate = anglerate
}

func (b *Bullet) Move() {
	rad := float64((b.Angle - 90) / float32(180) * math.Pi)
	vx := float32(math.Cos(rad))
	vy := float32(math.Sin(rad))

	if vx == 0 && vy == 0 {
		return
	}

	vector := engo.Point{X: vx, Y: vy}
	speed := float32(b.Speed) / float32(math.Sqrt(float64(vx*vx+vy*vy)))
	vector.MultiplyScalar(speed)

	b.AddPosition(vector)
	b.Angle += b.AngleRate
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
	bullet.Angle = 0
	bullet.SpeedRate = 0
	bullet.AngleRate = 0
	return BulletBuilder{bullet}
}

func (bb *BulletBuilder) Build() Modeler {
	bullet := new(Bullet)
	copier.Copy(&bullet, bb.Bullet)
	bullet.basicEntity = ecs.NewBasic()

	return *bullet
}
