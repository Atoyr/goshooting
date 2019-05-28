package entitys

import (
	"fmt"
	"math"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	acommon "github.com/atoyr/goshooting/common"
)

type NumberGroupBuilder struct {
	digitComponent []*Number
	value          int

	spaceComponent common.SpaceComponent
}

type NumberGroup struct {
	*NumberGroupBuilder
}

func NewNumberGroupBuilder(digit int, size acommon.NumberSize, scale float32, sc *common.SpaceComponent, margin float32) (*NumberGroupBuilder, error) {
	dc := make([]*Number, digit, digit)
	p := sc.Position
	for i := 1; i <= digit; i++ {
		nb, err := NewNumberBuilder(size)
		nb.SetVirtualPosition(p)
		if err != nil {
			return nil, err
		}
		nb.SetScale(scale)
		n := nb.Build()
		dc[digit-i] = &n
		p.Add(engo.Point{X: (acommon.GetNumberSize(size).X + margin) * scale, Y: 0})
	}
	return &NumberGroupBuilder{digitComponent: dc, value: -1}, nil
}

func (ngb *NumberGroupBuilder) Value(value int) {
	valueStr := fmt.Sprint(value)
	digit := len(valueStr)
	if max := int(math.Pow(10, float64(len(ngb.digitComponent)))); max <= value {
		value = max - 1
	}
	ngb.value = value

	for i := 0; i < len(ngb.digitComponent); i++ {
		var n int
		if i < digit {
			n = value / int(math.Pow(10, float64(i))) % 10
		} else {
			n = -1
		}
		ngb.digitComponent[i].SetNumber(n)
	}
}

func (ngb *NumberGroupBuilder) Build() NumberGroup {
	return NumberGroup{
		NumberGroupBuilder: ngb,
	}
}

func (ng *NumberGroup) AddedRenderSystem(rs *common.RenderSystem) {
	for _, dc := range ng.digitComponent {
		dc.AddedRenderSystem(rs)
	}
}
func (ng *NumberGroup) Add(value int) {
	ng.Value(ng.value + value)
}
