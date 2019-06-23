package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type Character struct {
	*EntityModel
	characterSize common.CharacterSize
	value         string
}

func (n *Character) SetCharacter(value string) error {
	n.value = value
	t, err := common.GetCharacterTexture(value, n.characterSize)
	if err != nil {
		return err
	}
	n.SetDrawable(t)
	n.renderComponent.Hidden = false
	return nil
}

type CharacterBuilder struct {
	*Character
}

func NewCharacterBuilder(size common.CharacterSize) (*CharacterBuilder, error) {
	sc := engoCommon.SpaceComponent{Position: engo.Point{X: 0, Y: 0}}
	r := engoCommon.RenderComponent{Hidden: true}
	em := EntityModel{spaceComponent: sc, renderComponent: r}
	c := Character{
		EntityModel:   &em,
		value:         "",
		characterSize: size,
	}
	return &CharacterBuilder{
		&c,
	}, nil
}

func (nb *CharacterBuilder) Build() Character {
	n := *nb.Character
	n.basicEntity = ecs.NewBasic()

	return n
}
