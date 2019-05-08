package system

import (
	"fmt"
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/common"
	"github.com/atoyr/goshooting/entitys"
)

// GameSystem is goshooting Game base system
type GameSystem struct {
	world        *ecs.World
	renderSystem *engoCommon.RenderSystem

	framecount             uint64
	playerEntity           *entitys.Entity
	playerBulletEntitys    map[uint64]*entitys.Entity
	enemyEntitys           map[uint64]*entitys.Entity
	enemyBulletEntitys     map[uint64]*entitys.Entity
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
	playerTexture := common.GetTexture("textures/player.png")
	enemyTexture := common.GetTexture("textures/enemy.png")
	common.GetTexture("textures/bullet3.png")

	// Create Entity
	// Player
	playerBuilder := entitys.NewEntityBuilder(
		&engoCommon.RenderComponent{
			Drawable: playerTexture,
			Scale:    engo.Point{X: 0.5, Y: 0.5},
		},
	)

	player := playerBuilder.BuildVirtualPosition(engo.Point{X: 0, Y: 0}).BuildSpeed(8).Build()
	player.Move = player.EntityMoveForPlayer

	playerBullets := map[uint64]*entitys.Entity{}
	enemys := map[uint64]*entitys.Entity{}
	enemyBullets := map[uint64]*entitys.Entity{}

	// enemy
	enemyBuilder := entitys.NewEntityBuilder(
		&engoCommon.RenderComponent{
			Drawable: enemyTexture,
			Scale:    engo.Point{X: 0.2, Y: 0.2},
		},
	)
	enemyBuilder.BuildZIndex(20)
	e := enemyBuilder.BuildVirtualPosition(engo.Point{X: 300, Y: 200}).Build()
	e.SetRotation(70)

	// Regist Entity
	gs.playerEntity = &player
	gs.playerBulletEntitys = playerBullets
	gs.enemyEntitys = enemys
	gs.enemyBulletEntitys = enemyBullets
	gs.enemyBulletCount = 0

	enemys[e.GetID()] = &e

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engoCommon.RenderSystem:
			gs.renderSystem = sys
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
	bulletEntity := make([]*entitys.Entity, 0, 10)

	// get inputs
	isleft := engo.Input.Button("MoveLeft").Down()
	isright := engo.Input.Button("MoveRight").Down()
	isup := engo.Input.Button("MoveUp").Down()
	isdown := engo.Input.Button("MoveDown").Down()
	islowspeed := engo.Input.Button("LowSpeed").Down()
	isshot := engo.Input.Button("Shot").Down()

	// Player Update
	gs.playerEntity.Move(gs.playerEntity.GetPlayerMoveInfo(isleft, isright, isup, isdown, islowspeed))

	// PlayerBullet Update
	for _, pb := range gs.playerBulletEntitys {
		pb.Move(0, -1, pb.GetSpeed())
	}

	if isshot && gs.playerBulletStartCount%5 == 0 {
		gs.playerBulletStartCount += 1
		bulletTexture := common.GetTexture("textures/bullet3.png")
		bb := entitys.NewEntityBuilder(
			&engoCommon.RenderComponent{
				Drawable: bulletTexture,
				Scale:    engo.Point{X: 0.5, Y: 0.5},
			})

		bb.BuildVirtualPosition(gs.playerEntity.GetVirtualPosition()).BuildSpeed(16)
		b := bb.Build()
		b.SetZIndex(10)
		bulletEntity = append(bulletEntity, &b)
		gs.playerBulletEntitys[b.GetID()] = &b
	} else if !isshot {
		gs.playerBulletStartCount = 0
	} else {
		gs.playerBulletStartCount += 1

	}

	// EnemyBullet Upate
	for _, eb := range gs.enemyBulletEntitys {
		s := common.NewSetting()
		eb.Move(eb.GetMoveInfo())
		virtualPosition := eb.GetVirtualPosition()
		mergin := eb.GetMergin()
		if (virtualPosition.X < -1*mergin.X || s.GetGameAreaSize().X+mergin.X < virtualPosition.X) || (virtualPosition.Y < -1*mergin.Y || s.GetGameAreaSize().Y+mergin.Y < virtualPosition.Y) {
			gs.Remove(eb.GetBasicEntity())
		}
	}

	if true {
		for _, e := range gs.enemyEntitys {
			gs.enemyBulletCount += 1
			for i := 0; i < 4; i++ {
				bulletTexture := common.GetTexture("textures/bullet2.png")
				bb := entitys.NewEntityBuilder(
					&engoCommon.RenderComponent{
						Drawable: bulletTexture,
						Scale:    engo.Point{X: 1.0, Y: 1.0},
					})
				bb.BuildVirtualPosition(e.GetVirtualPosition()).BuildSpeed(4).BuildAngle(float32(gs.framecount%72*5) + float32(90*i) + float32(gs.enemyBulletCount))
				// Angle(
				bb.SetZIndex(10)
				b := bb.Build()
				bulletEntity = append(bulletEntity, &b)
				gs.enemyBulletEntitys[b.GetID()] = &b
			}
		}
	}
	for _, system := range gs.world.Systems() {
		switch sys := system.(type) {
		case *engoCommon.RenderSystem:
			for _, b := range bulletEntity {
				b.AddedRenderSystem(sys)
			}
		}
	}

	if gs.framecount%60 == 0 {
		fmt.Printf("Entity count %d \n", len(gs.enemyBulletEntitys))
	}
}

func (gs *GameSystem) Remove(b ecs.BasicEntity) {
	for _, system := range gs.world.Systems() {
		switch sys := system.(type) {
		case *engoCommon.RenderSystem:
			sys.Remove(b)
		}
	}
	if _, ok := gs.playerBulletEntitys[b.ID()]; ok {

		delete(gs.playerBulletEntitys, b.ID())
	}
	if _, ok := gs.enemyEntitys[b.ID()]; ok {
		delete(gs.enemyEntitys, b.ID())
	}
	if _, ok := gs.enemyBulletEntitys[b.ID()]; ok {
		delete(gs.enemyBulletEntitys, b.ID())
	}
}
