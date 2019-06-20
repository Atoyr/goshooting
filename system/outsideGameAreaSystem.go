package system

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	acommon "github.com/atoyr/goshooting/common"
)

type OutsideGameAreaSystem struct {
}

func (oga *OutsideGameAreaSystem) New(w *ecs.World) {
	color := color.RGBA{100, 100, 100, 255}
	rect := common.Rectangle{BorderWidth: 0}
	l := &common.RenderComponent{
		Drawable: rect,
		Color:    color,
	}

	r := &common.RenderComponent{
		Drawable: rect,
		Color:    color,
	}
	t := &common.RenderComponent{
		Drawable: rect,
		Color:    color,
	}
	b := &common.RenderComponent{
		Drawable: rect,
		Color:    color,
	}
	l.SetZIndex(500)
	r.SetZIndex(500)
	t.SetZIndex(500)
	b.SetZIndex(500)

	setting := acommon.NewSetting()
	leftBasic := ecs.NewBasic()
	leftSpace := &common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Width:    setting.AABB().Min.X,
		Height:   setting.RenderCanvas().Y,
	}
	rightBasic := ecs.NewBasic()
	rightSpace := &common.SpaceComponent{
		Position: engo.Point{X: setting.AABB().Max.X, Y: 0},
		Width:    setting.RenderCanvas().X - setting.AABB().Max.X,
		Height:   setting.RenderCanvas().Y,
	}
	topBasic := ecs.NewBasic()
	topSpace := &common.SpaceComponent{
		Position: engo.Point{X: setting.AABB().Min.X, Y: 0},
		Width:    setting.RenderGameAreaSize().X,
		Height:   setting.AABB().Min.Y,
	}
	bottomBasic := ecs.NewBasic()
	bottomSpace := &common.SpaceComponent{
		Position: engo.Point{X: setting.AABB().Min.X, Y: setting.AABB().Max.Y},
		Width:    setting.RenderGameAreaSize().X,
		Height:   setting.RenderCanvas().Y - setting.AABB().Max.Y,
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&leftBasic, l, leftSpace)
			sys.Add(&rightBasic, r, rightSpace)
			sys.Add(&topBasic, t, topSpace)
			sys.Add(&bottomBasic, b, bottomSpace)
		}
	}
}
func (oga *OutsideGameAreaSystem) Update(dt float32) {
}

func (oga *OutsideGameAreaSystem) Remove(basic ecs.BasicEntity) {
}
