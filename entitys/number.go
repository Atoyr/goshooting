package entitys

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
)

type NumberBuilder struct {
	*Entity
	numbers []*engoCommon.RenderComponent
}

type Number struct {
	*NumberBuilder
	value int
}

func NewNumberBuilder(size common.NumberSize, scale engo.Point, sc *engoCommon.SpaceComponent) (*NumberBuilder, error) {
	n := make([]*engoCommon.RenderComponent, 10, 10)
	t, err := common.GetNumberTextures(size)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for i := range t {
		rc := engoCommon.RenderComponent{
			Drawable: t[i],
			Scale:    scale,
			Hidden:   true,
		}
		n[i] = &rc
	}
	em := &Entity{
		basicEntity:            ecs.NewBasic(),
		renderComponent:        n[0],
		spaceComponent:         sc,
		virtualPosition:        engo.Point{X: 0, Y: 0},
		collisionDetectionSize: 0,
		mergin:                 engo.Point{X: 0, Y: 0},
	}
	em.MoveFunc = em.EntityMove
	return &NumberBuilder{
		Entity:  em,
		numbers: n,
	}, nil
}

// AddedRenderSystem is added render system
func (n *Number) AddedRenderSystem(rs *engoCommon.RenderSystem) {
	for _, rc := range n.numbers {
		b := ecs.NewBasic()
		n.basicEntity.AppendChild(&b)
		rs.Add(&b, rc, n.spaceComponent)
	}
}

func (e *NumberBuilder) SetVirtualPosition(xy engo.Point) *NumberBuilder {
	e.virtualPosition = xy
	return e
}

func (nb *NumberBuilder) Build() *Number {
	return &Number{
		NumberBuilder: nb,
		value:         -1,
	}
}

func (n *Number) SetZIndex(value float32) {
	for i := range n.numbers {
		n.numbers[i].SetZIndex(value)
	}
}

func (n *Number) SetNumber(value int) {
	num := value % 10
	n.value = num
	for i := range n.numbers {
		n.numbers[i].Hidden = value == -1 || num != i
	}
}

func (n *Number) Add(value int) {
	if n.value == -1 {
		return
	}

	num := (n.value + value) % 10
	n.SetNumber(num)
}
