package system

import (
	"fmt"
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	acommon "github.com/atoyr/goshooting/common"
	"github.com/atoyr/goshooting/entitys"
)

// GameSystem is goshooting Game base system
type GameSystem struct {
	world *ecs.World

	framecount             uint64
	playerEntity           *entitys.Player
	playerBulletEntitys    map[uint64]*entitys.Bullet
	enemyEntitys           map[uint64]*entitys.Enemy
	enemyBulletEntitys     map[uint64]*entitys.Bullet
	playerBulletStartCount uint64
	enemyBulletCount       uint64
}

// New is Startup Entity
func (gs *GameSystem) New(w *ecs.World) {
	// debug message
	fmt.Printf("Canvas Width:%f Height:%f Scale:%f \n", engo.CanvasWidth(), engo.CanvasHeight(), engo.CanvasScale())
	fmt.Printf("Window Width:%f Height:%f  \n", engo.WindowWidth(), engo.WindowHeight())

	gs.world = w
	gs.framecount = 0
	gs.playerBulletStartCount = 0

	// load texture
	playerTexture := acommon.GetTexture("textures/player.png")
	enemyTexture := acommon.GetTexture("textures/enemy.png")
	acommon.GetTexture("textures/bullet3.png")

	// Create Entity
	// Player
	playerBuilder := entitys.NewPlayerBuilder(
		&common.RenderComponent{
			Drawable: playerTexture,
			Scale:    engo.Point{X: 0.5, Y: 0.5},
		},
		&common.SpaceComponent{
			Position: engo.Point{X: 0, Y: 0},
			Width:    8,
			Height:   8,
		},
	)

	player := playerBuilder.VirtualPosition(engo.Point{X: 0, Y: 0}).Size(32).Mergin(engo.Point{X: 32, Y: 32}).LowSpeed(4).Speed(8).Build()

	playerBullets := map[uint64]*entitys.Bullet{}
	enemys := map[uint64]*entitys.Enemy{}
	enemyBullets := map[uint64]*entitys.Bullet{}

	// enemy
	enemyBuilder := entitys.NewEnemyBuilder(
		&common.RenderComponent{
			Drawable: enemyTexture,
			Scale:    engo.Point{X: 0.2, Y: 0.2},
		},
		&common.SpaceComponent{
			Width:  8,
			Height: 8,
		},
	)
	enemyBuilder.RenderComponent.SetZIndex(20)
	e := enemyBuilder.VirtualPosition(engo.Point{X: 300, Y: 200}).Mergin(engo.Point{X: 32, Y: 32}).Build()

	// Regist Entity
	gs.playerEntity = player
	gs.playerBulletEntitys = playerBullets
	gs.enemyEntitys = enemys
	gs.enemyBulletEntitys = enemyBullets
	gs.enemyBulletCount = 0

	enemys[e.GetId()] = e

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			player.AddedRenderSystem(sys)
			e.AddedRenderSystem(sys)
		}
	}
}

// Update is Frame Update
func (gs *GameSystem) Update(dt float32) {
	gs.framecount += 1
	if gs.framecount == math.MaxInt64 {
		gs.framecount = 0
	}
	bulletEntity := make([]*entitys.Bullet, 0, 10)

	// get inputs
	isleft := engo.Input.Button("MoveLeft").Down()
	isright := engo.Input.Button("MoveRight").Down()
	isup := engo.Input.Button("MoveUp").Down()
	isdown := engo.Input.Button("MoveDown").Down()
	islowspeed := engo.Input.Button("LowSpeed").Down()
	isshot := engo.Input.Button("Shot").Down()

	// Player Update
	gs.playerEntity.Move(gs.playerEntity.GetMoveInfo(isleft, isright, isup, isdown, islowspeed))

	// PlayerBullet Update
	for _, pb := range gs.playerBulletEntitys {
		pb.Move(0, -1, pb.GetSpeed())
	}

	if isshot && gs.playerBulletStartCount%5 == 0 {
		gs.playerBulletStartCount += 1
		bulletTexture := acommon.GetTexture("textures/bullet3.png")
		bb := entitys.NewBulletBuilder(
			&common.RenderComponent{
				Drawable: bulletTexture,
				Scale:    engo.Point{X: 0.5, Y: 0.5},
			},
			&common.SpaceComponent{
				Width:  32,
				Height: 32,
			})

		bb.VirtualPosition(gs.playerEntity.GetVPoint()).Size(32).Speed(16).Mergin(engo.Point{X: 32, Y: 32})
		b := bb.Build()
		bulletEntity = append(bulletEntity, b)
		gs.playerBulletEntitys[b.GetId()] = b
	} else if !isshot {
		gs.playerBulletStartCount = 0
	} else {
		gs.playerBulletStartCount += 1

	}

	// EnemyBullet Upate
	for _, eb := range gs.enemyBulletEntitys {
		eb.Move(eb.GetMoveInfo())
	}

	if true {
		for _, e := range gs.enemyEntitys {
			gs.enemyBulletCount += 1
			for i := 0; i < 4; i++ {
				bulletTexture := acommon.GetTexture("textures/bullet2.png")
				bb := entitys.NewBulletBuilder(
					&common.RenderComponent{
						Drawable: bulletTexture,
						Scale:    engo.Point{X: 1.0, Y: 1.0},
					},
					&common.SpaceComponent{
						Width:  32,
						Height: 32,
					})
				bb.VirtualPosition(e.EntityModel.VirtualPosition).Size(32).Speed(4).Angle(float32(gs.framecount%72*5) + float32(90*i) + float32(gs.enemyBulletCount)).AngleRate(0)
				bb.RenderComponent.SetZIndex(10)
				b := bb.Build()
				bulletEntity = append(bulletEntity, b)
				gs.enemyBulletEntitys[b.GetId()] = b
			}
		}
	}
	for _, system := range gs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, b := range bulletEntity {
				b.AddedRenderSystem(sys)
			}
		}
	}

}

func (gs *GameSystem) Remove(ecs.BasicEntity) {}
