package main

import (
	"bytes"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/system"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

type myScene struct{}

func (*myScene) Type() string { return "goshooting" }

func (*myScene) Preload() {
	// File Load
	engo.Files.Load("textures/player.png")
	engo.Files.Load("textures/bullet1.png")
	engo.Files.Load("textures/bullet2.png")
	engo.Files.Load("textures/bullet3.png")
	engo.Files.Load("textures/enemy.png")
	engo.Files.Load("textures/number_8_48.png")
	engo.Files.Load("textures/number_16_16.png")
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	//common.SetBackground(color.Black)
}

func (*myScene) Setup(u engo.Updater) {
	// Button Register
	engo.Input.RegisterButton("MoveRight", engo.KeyL, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyJ, engo.KeyArrowLeft)
	engo.Input.RegisterButton("MoveUp", engo.KeyI, engo.KeyArrowUp)
	engo.Input.RegisterButton("MoveDown", engo.KeyK, engo.KeyArrowDown)
	engo.Input.RegisterButton("LowSpeed", engo.KeyLeftShift)
	engo.Input.RegisterButton("Shot", engo.KeyZ)
	world, _ := u.(*ecs.World)
	fps := common.FPSSystem{Display: false, Terminal: true}

	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.AnimationSystem{})
	world.AddSystem(&system.GameSystem{})
	world.AddSystem(&system.HUDSystem{})
	world.AddSystem(&system.OutsideGameAreaSystem{})
	world.AddSystem(&fps)
}

func (*myScene) Exit() {
	engo.Exit()
}

func main() {
	opts := engo.RunOptions{
		Title:          "goshooting",
		Width:          1280,
		Height:         720,
		StandardInputs: true,
		NotResizable:   true,
	}
	engo.Run(opts, &myScene{})
}
