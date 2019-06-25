package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
	"github.com/jinzhu/copier"
	"math"
)

type Player struct {
	*EntityModel
	LowSpeed float32
	Speed    float32
	Attack   func(modeler Modeler, frame uint64) []Modeler
}

func (p *Player) Vector(isleft, isright, isup, isdown, islowspeed bool) engo.Point {
	vx := float32(0)
	vy := float32(0)
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
	vector := engo.Point{X: vx, Y: vy}
	if vx == 0 && vy == 0 {
		return vector
	}
	speed := p.Speed
	if islowspeed {
		speed = p.LowSpeed
	}
	speed = speed / float32(math.Sqrt(float64(vx*vx+vy*vy)))
	vector.MultiplyScalar(speed)

	return vector
}

type PlayerBuilder struct {
	*Player
}

func NewPlayerBuilder() PlayerBuilder {
	s := common.NewSetting()
	sc := engoCommon.SpaceComponent{Position: engo.Point{X: 0, Y: 0}}
	rc := engoCommon.RenderComponent{Scale: s.Scale()}
	model := EntityModel{
		spaceComponent:  sc,
		renderComponent: rc,
		virtualPosition: engo.Point{X: 0, Y: 0},
		scale:           1,
	}
	model.renderComponent.Scale.MultiplyScalar(model.scale)
	model.SetPosition(engo.Point{X: 0, Y: 0})

	player := new(Player)
	player.EntityModel = &model

	return PlayerBuilder{player}
}

func (pb *PlayerBuilder) Build() Modeler {
	player := new(Player)
	copier.Copy(&player, pb.Player)
	player.basicEntity = ecs.NewBasic()

	return *player
}
