package system

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
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
}

// New is called when the system is added to the world
func (hud *HUDSystem) New(w *ecs.World) {
	c := color.RGBA{100, 100, 100, 255}
	r := common.Rectangle{BorderWidth: 4, BorderColor: c}

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
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&basic, render, space)
		}
	}

}

func (hud *HUDSystem) Update(dt float32) {

}

func (hud *HUDSystem) Remove(basic ecs.BasicEntity) {

}
