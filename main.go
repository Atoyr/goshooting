package main

import (
	"bytes"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/system"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

type myScene struct{}

func (*myScene) Type() string { return "goshooting" }

func (*myScene) Preload() {
	engo.Files.Load("textures/player.png")
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	common.SetBackground(color.White)
}

func (*myScene) Setup(u engo.Updater) {
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("MoveUp", engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton("MoveDown", engo.KeyS, engo.KeyArrowDown)
	engo.Input.RegisterButton("LowSpeed", engo.KeySpace)
	world, _ := u.(*ecs.World)

	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&system.PlayerSystem{})
}

func (*myScene) Exit() {
	engo.Exit()
}

func main() {
	opts := engo.RunOptions{
		Title:          "goshooting",
		Width:          400,
		Height:         300,
		StandardInputs: true,
		NotResizable:   true,
	}
	engo.Run(opts, &myScene{})
}
