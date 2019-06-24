package entitys

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type NumberBuilder struct {
	*Number
}

type Number struct {
	*EntityModel
	texture []engoCommon.Texture
	value   int
}

func NewNumberBuilder(size common.NumberSize) (*NumberBuilder, error) {
	t, err := common.GetNumberTextures(size)
	if err != nil {
		return nil, err
	}

	sc := engoCommon.SpaceComponent{Position: engo.Point{X: 0, Y: 0}}
	r := engoCommon.RenderComponent{Drawable: t[0], Hidden: true}
	em := EntityModel{spaceComponent: sc, renderComponent: r}
	emover := EntityMove{}
	attack := new(EntityAttack)
	collison := new(EntityCollision)
	e := Entity{EntityModel: &em, EntityMove: &emover, EntityAttack: attack, EntityCollision: collison}
	n := Number{
		Entity:  &e,
		value:   -1,
		texture: t,
	}
	return &NumberBuilder{
		&n,
	}, nil
}

func (nb *NumberBuilder) Build() Number {
	n := *nb.Number
	n.basicEntity = ecs.NewBasic()

	return n
}

func (n *Number) SetNumber(value int) {
	if value < 0 {
		n.renderComponent.Hidden = true
	} else {
		num := value % 10
		n.value = num
		n.SetDrawable(n.texture[num])
		n.renderComponent.Hidden = false
	}
}

func (n *Number) Add(value int) {
	if n.value == -1 {
		return
	}

	num := (n.value + value) % 10
	n.SetNumber(num)
}
