package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
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
	em.MoveFunc = em.EntityMove
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

func (e *Enemy) GetSpeed() float32 {
	return e.speed
}