package system

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	acommon "github.com/atoyr/goshooting/common"
	"github.com/atoyr/goshooting/entitys"
)

type HUDScoreMessage struct {
	Score int
}

const HUDScoreMessageType string = "HUDScoreMessage"

func (HUDScoreMessage) Type() string {
	return HUDScoreMessageType
}

type HUDSystem struct {
	score int
	frame uint64

	number *entitys.Number
}

// New is called when the system is added to the world
func (hud *HUDSystem) New(w *ecs.World) {
	c := color.RGBA{100, 100, 100, 255}
	r := common.Rectangle{BorderWidth: 4, BorderColor: c}
	hud.frame = 0

	basic := ecs.NewBasic()
	render := &common.RenderComponent{
		Drawable: r,
		Scale:    engo.Point{X: 1, Y: 1},
	}
	space := &common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Width:    100,
		Height:   8,
	}

	acommon.InitializeNumber_8_48("textures/number_8_48.png")
	space2 := &common.SpaceComponent{
		Position: engo.Point{X: 150, Y: 150},
		Width:    150,
		Height:   150,
	}
	numbuilder, err := entitys.NewNumberBuilder(acommon.Number_8_48, engo.Point{X: 2, Y: 2}, space2)
	if err != nil {
		return
	}
	hud.number = numbuilder.Build()
	hud.number.SetNumber(0)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&basic, render, space)
			hud.number.AddedRenderSystem(sys)
		}
	}

}

func (hud *HUDSystem) Update(dt float32) {
	hud.frame++
	if hud.frame%60 == 0 {
		hud.number.Add(1)
	}
}

func (hud *HUDSystem) Remove(basic ecs.BasicEntity) {

}
