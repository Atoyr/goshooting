package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type BulletBuilder struct {
	*EntityModel

	speed float32
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
	}
	em.MoveFunc = em.EntityMove
	return &BulletBuilder{
		EntityModel: &em,
		speed:       0,
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
func (b *BulletBuilder) Build() *Bullet {
	return &Bullet{
		BulletBuilder: b,
	}
}

func (b *Bullet) GetSpeed() float32 {
	return b.speed
}
