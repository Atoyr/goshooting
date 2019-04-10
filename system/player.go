package system

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent

	size     int
	lowSpeed int
	speed    int
}

type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
}

func (ps *PlayerSystem) New(w *ecs.World) {
	ps.world = w

	player := Player{BasicEntity: ecs.NewBasic()}

	player.lowSpeed = 4
	player.speed = 8

	// initialize position
	posX := int(engo.WindowWidth() / 2)
	posY := int(engo.WindowHeight() / 2)
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(posX), Y: float32(posY)},
		Width:    32,
		Height:   32,
	}

	// load image
	texture, err := common.LoadedSprite("textures/player.png")
	if err != nil {
		fmt.Println("Unable to load texture: " + err.Error())
	}
	player.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 0.1, Y: 0.1},
	}
	player.RenderComponent.SetZIndex(1)
	ps.playerEntity = &player
	ps.texture = texture

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}
}

func (ps *PlayerSystem) Update(dt float32) {
	// get inputs
	isleft := engo.Input.Button("MoveLeft").Down()
	isright := engo.Input.Button("MoveRight").Down()
	isup := engo.Input.Button("MoveUp").Down()
	isdown := engo.Input.Button("MoveDown").Down()
	islowSpeed := engo.Input.Button("LowSpeed").Down()

	var speed float32
	if islowSpeed {
		speed = float32(ps.playerEntity.lowSpeed)
	} else {
		speed = float32(ps.playerEntity.speed)
	}
	diagspeed := float32(speed) / 1.414

	if isleft {
		if !isup && isdown {
			ps.playerEntity.SpaceComponent.Position.X -= diagspeed
			ps.playerEntity.SpaceComponent.Position.Y += diagspeed
		} else if isup && !isdown {
			ps.playerEntity.SpaceComponent.Position.X -= diagspeed
			ps.playerEntity.SpaceComponent.Position.Y -= diagspeed
		} else {
			ps.playerEntity.SpaceComponent.Position.X -= speed
		}
	}
	if isright {
		if !isup && isdown {
			ps.playerEntity.SpaceComponent.Position.X += diagspeed
			ps.playerEntity.SpaceComponent.Position.Y += diagspeed
		} else if isup && !isdown {
			ps.playerEntity.SpaceComponent.Position.X += diagspeed
			ps.playerEntity.SpaceComponent.Position.Y -= diagspeed
		} else {
			ps.playerEntity.SpaceComponent.Position.X += speed
		}
	}
	if isup {
		if !isleft && isright {
			ps.playerEntity.SpaceComponent.Position.X += diagspeed
			ps.playerEntity.SpaceComponent.Position.Y -= diagspeed
		} else if isleft && !isright {
			ps.playerEntity.SpaceComponent.Position.X -= diagspeed
			ps.playerEntity.SpaceComponent.Position.Y -= diagspeed
		} else {
			ps.playerEntity.SpaceComponent.Position.Y -= speed
		}
	}
	if isdown {
		if !isleft && isright {
			ps.playerEntity.SpaceComponent.Position.X += diagspeed
			ps.playerEntity.SpaceComponent.Position.Y += diagspeed
		} else if isleft && !isright {
			ps.playerEntity.SpaceComponent.Position.X -= diagspeed
			ps.playerEntity.SpaceComponent.Position.Y += diagspeed
		} else {
			ps.playerEntity.SpaceComponent.Position.Y += speed
		}
	}
}

func (*PlayerSystem) Remove(ecs.BasicEntity) {
}
