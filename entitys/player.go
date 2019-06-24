package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
	"github.com/jinzhu/copier"
)

type Player struct {
	*EntityModel
	LowSpeed float32
	Speed    float32
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
