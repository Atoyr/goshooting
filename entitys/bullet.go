package entitys

import (
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	acommon "github.com/atoyr/goshooting/common"
)

type BulletBuilder struct {
	*EntityModel

	speed     float32
	angle     float32
	speedRate float32
	angleRate float32
}

// Bullet is player shot
type Bullet struct {
	*BulletBuilder
}

func NewBulletBuilder(rc *common.RenderComponent, sc *common.SpaceComponent) *BulletBuilder {
	em := EntityModel{
		BasicEntity:     ecs.NewBasic(),
		RenderComponent: *rc,
		SpaceComponent:  *sc,
		VirtualPosition: engo.Point{X: 0, Y: 0},
		Size:            0,
		Mergin:          engo.Point{X: 0, Y: 0},
		IsRemoveTraget:  false,
	}
	em.MoveFunc = em.EntityMove
	return &BulletBuilder{
		EntityModel: &em,
		speed:       0,
		angle:       0,
		speedRate:   0,
		angleRate:   0,
	}
}

func (b *BulletBuilder) VirtualPosition(xy engo.Point) *BulletBuilder {
	b.EntityModel.VirtualPosition = xy
	return b
}
func (b *BulletBuilder) Size(s float32) *BulletBuilder {
	b.EntityModel.Size = s
	return b
}
func (b *BulletBuilder) Mergin(m engo.Point) *BulletBuilder {
	b.EntityModel.Mergin = m
	return b
}
func (b *BulletBuilder) Speed(s float32) *BulletBuilder {
	b.speed = s
	return b
}

func (b *BulletBuilder) Angle(a float32) *BulletBuilder {
	b.angle = a
	return b
}

func (b *BulletBuilder) SpeedRate(sr float32) *BulletBuilder {
	b.speedRate = sr
	return b
}

func (b *BulletBuilder) AngleRate(ar float32) *BulletBuilder {
	b.angleRate = ar
	return b
}

func (b *BulletBuilder) Build() *Bullet {
	s := acommon.NewSetting()
	p := s.ConvertRenderPosition(b.EntityModel.VirtualPosition)
	drawableSize := engo.Point{X: b.RenderComponent.Drawable.Width() / 2, Y: b.RenderComponent.Drawable.Height() / 2}
	drawableSize.Multiply(b.RenderComponent.Scale)
	drawableSize.MultiplyScalar(-1)
	p.Add(drawableSize)
	b.SpaceComponent.Position = p
	return &Bullet{
		BulletBuilder: b,
	}
}

func (b *Bullet) GetSpeed() float32 {
	return b.speed
}

func (b *Bullet) GetMoveInfo() (vx, vy, speed float32) {
	rad := float64((b.angle - 90) / float32(180) * math.Pi)
	vx = float32(math.Cos(rad))
	vy = float32(math.Sin(rad))
	speed = b.speed

	//fmt.Printf("%f %f %f %f  X :  %f  Y:  %f  \n", vx*speed, vy*speed, speed, rad, b.SpaceComponent.Position.X, b.SpaceComponent.Position.Y)

	b.angle += b.angleRate
	b.speed += b.speedRate

	return
}
