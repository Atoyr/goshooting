package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	acommon "github.com/atoyr/goshooting/common"
)

type PlayerBuilder struct {
	*EntityModel

	lowSpeed float32
	speed    float32
}

type Player struct {
	*PlayerBuilder
}

func NewPlayerBuilder(rc *common.RenderComponent, sc *common.SpaceComponent) *PlayerBuilder {
	em := &EntityModel{
		BasicEntity:     ecs.NewBasic(),
		RenderComponent: *rc,
		SpaceComponent:  *sc,
		VirtualPosition: engo.Point{X: 0, Y: 0},
		Size:            0,
		Mergin:          engo.Point{X: 0, Y: 0},
	}
	return &PlayerBuilder{
		EntityModel: em,
		lowSpeed:    0,
		speed:       0,
	}
}

func (b *PlayerBuilder) VirtualPosition(xy engo.Point) *PlayerBuilder {
	b.EntityModel.VirtualPosition = xy
	return b
}

func (b *PlayerBuilder) Size(s float32) *PlayerBuilder {
	b.EntityModel.Size = s
	return b
}

func (b *PlayerBuilder) Mergin(m engo.Point) *PlayerBuilder {
	b.EntityModel.Mergin = m
	return b
}

func (b *PlayerBuilder) LowSpeed(l float32) *PlayerBuilder {
	b.lowSpeed = l
	return b
}

func (b *PlayerBuilder) Speed(s float32) *PlayerBuilder {
	b.speed = s
	return b
}

func (b *PlayerBuilder) Build() *Player {
	return &Player{
		b,
	}
}

// Move is move
func (p *Player) Move(vx, vy, speed float32) {
	s := acommon.NewSetting()
	x := p.EntityModel.VirtualPosition.X
	y := p.EntityModel.VirtualPosition.Y

	if vx != 0 && vy != 0 {
		speed = float32(speed) / 1.414
	}
	x += speed * vx
	y += speed * vy
	max := s.GetGameAreaSize()
	if x < p.EntityModel.Mergin.X {
		x = p.EntityModel.Mergin.X
	} else if x > max.X-p.EntityModel.Mergin.X {
		x = max.X - p.EntityModel.Mergin.X
	}
	if y < p.EntityModel.Mergin.Y {
		y = p.EntityModel.Mergin.Y
	} else if y > max.Y-p.EntityModel.Mergin.Y {
		y = max.Y - p.EntityModel.Mergin.Y
	}
	p.EntityModel.VirtualPosition.X = x
	p.EntityModel.VirtualPosition.Y = y

	ret := s.ConvertRenderPosition(p.EntityModel.convertPosition())
	p.SpaceComponent.Position = ret

}

func (p *Player) GetMoveInfo(isleft, isright, isup, isdown, islowspeed bool) (vx, vy, speed float32) {
	vx = 0
	vy = 0
	if islowspeed {
		speed = float32(p.lowSpeed)
	} else {
		speed = float32(p.speed)
	}
	if isleft && !isright {
		vx = -1
	} else if !isleft && isright {
		vx = 1
	}
	if isup && !isdown {
		vy = -1
	} else if !isup && isdown {
		vy = 1
	}
	return vx, vy, speed
}
