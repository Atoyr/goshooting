package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	acommon "github.com/atoyr/goshooting/common"
)

type EnemyBuilder struct {
	*EntityModel

	speed float32
}

// Enemy is player shot
type Enemy struct {
	*EnemyBuilder
}

func NewEnemyBuilder(rc *common.RenderComponent, sc *common.SpaceComponent) *EnemyBuilder {
	em := EntityModel{
		BasicEntity:     ecs.NewBasic(),
		RenderComponent: *rc,
		SpaceComponent:  *sc,
		VirtualPosition: engo.Point{X: 0, Y: 0},
		Size:            0,
		Mergin:          engo.Point{X: 0, Y: 0},
	}
	return &EnemyBuilder{
		EntityModel: &em,
		speed:       0,
	}
}

func (e *EnemyBuilder) VirtualPosition(xy engo.Point) *EnemyBuilder {
	e.EntityModel.VirtualPosition = xy
	return e
}
func (e *EnemyBuilder) Size(s float32) *EnemyBuilder {
	e.EntityModel.Size = s
	return e
}
func (e *EnemyBuilder) Mergin(m engo.Point) *EnemyBuilder {
	e.EntityModel.Mergin = m
	return e
}
func (e *EnemyBuilder) Speed(s float32) *EnemyBuilder {
	e.speed = s
	return e
}
func (e *EnemyBuilder) Build() *Enemy {
	return &Enemy{
		EnemyBuilder: e,
	}
}

func (e *Enemy) Move(vx, vy, speed, angle float32) engo.Point {
	s := acommon.NewSetting()
	e.EntityModel.VirtualPosition.X += vx * speed
	e.EntityModel.VirtualPosition.Y += vy * speed
	ret := s.ConvertRenderPosition(e.EntityModel.convertPosition())
	e.EntityModel.SpaceComponent.Position = ret

	return ret
}
func (e *Enemy) GetSpeed() float32 {
	return e.speed
}

func (e *Enemy) Atack(frame uint64, playerPoint engo.Point) []*Bullet {
	bullets := make([]*Bullet, 1)

	return bullets
}
