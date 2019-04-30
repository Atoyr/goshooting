package entitys

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	acommon "github.com/atoyr/goshooting/common"
)

type NumberBuilder struct {
	*EntityModel
	numbers []*common.RenderComponent
}

type Number struct {
	*NumberBuilder
	value int
}

func NewNumberBuilder(size acommon.NumberSize, scale engo.Point, sc *common.SpaceComponent) (*NumberBuilder, error) {
	n := make([]*common.RenderComponent, 10, 10)
	t, err := acommon.GetNumberTextures(size)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for i := range t {
		rc := common.RenderComponent{
			Drawable: t[i],
			Scale:    scale,
			Hidden:   true,
		}
		n[i] = &rc
	}
	em := &EntityModel{
		BasicEntity:     ecs.NewBasic(),
		RenderComponent: *n[0],
		SpaceComponent:  *sc,
		VirtualPosition: engo.Point{X: 0, Y: 0},
		Size:            0,
		Mergin:          engo.Point{X: 0, Y: 0},
	}
	return &NumberBuilder{
		EntityModel: em,
		numbers:     n,
	}, nil
}

// AddedRenderSystem is added render system
func (ne *Number) AddedRenderSystem(rs *common.RenderSystem) {
	for _, rc := range ne.numbers {
		b := ecs.NewBasic()
		ne.BasicEntity.AppendChild(&b)
		rs.Add(&b, rc, &ne.SpaceComponent)
	}
}

func (e *NumberBuilder) VirtualPosition(xy engo.Point) *NumberBuilder {
	e.EntityModel.VirtualPosition = xy
	return e
}
func (e *NumberBuilder) Size(s float32) *NumberBuilder {
	e.EntityModel.Size = s
	return e
}
func (e *NumberBuilder) Mergin(m engo.Point) *NumberBuilder {
	e.EntityModel.Mergin = m
	return e
}

func (nb *NumberBuilder) Build() *Number {
	return &Number{
		NumberBuilder: nb,
		value:         -1,
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
