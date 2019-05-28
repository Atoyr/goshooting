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
	playerBuilder := entitys.NewPlayerBuilder()
	playerBuilder.SetDrawable(playerTexture)
	playerBuilder.SetVirtualPosition(engo.Point{X: 0, Y: 0})
	playerBuilder.SetSpeed(8)
	playerBuilder.SetCollisionDetectionSize(8)
	playerBuilder.SetZIndex(10)
	player := playerBuilder.Build()
	player.RenderCollisionDetection(true)

	playerBullets := map[uint64]*entitys.Entity{}
	enemys := map[uint64]*entitys.Entity{}
	enemyBullets := map[uint64]*entitys.Entity{}

	// enemy
	enemyBuilder := entitys.NewEnemyBuilder()
	enemyBuilder.SetDrawable(enemyTexture)
	enemyBuilder.SetZIndex(20)
	enemyBuilder.SetVirtualPosition(engo.Point{X: 300, Y: 200})
	enemyBuilder.SetAngle(70)
	enemyBuilder.SetCollisionDetectionSize(25)
	e := enemyBuilder.Build()

	// Regist Entity
	gs.playerEntity = &player
	gs.playerBulletEntitys = playerBullets
	gs.enemyEntitys = enemys
	gs.enemyBulletEntitys = enemyBullets
	gs.enemyBulletCount = 0

	enemys[e.ID()] = &e

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
	vx, vy := getPlayerMoveInfo(isleft, isright, isup, isdown)
	if islowspeed {
		gs.playerEntity.Speed = 4
	} else {
		gs.playerEntity.Speed = 8
	}
	gs.playerEntity.Move(gs.playerEntity, vx, vy)

	// Collision
	// if gs.framecount%60 == 0 {
	// 	for _, e := range gs.enemyEntitys {
	// 		if gs.playerEntity.IsCollision(e) {
	// 			fmt.Println("Collision!!")
	// 		}
	// 	}
	// }

	// PlayerBullet Update
	for _, pb := range gs.playerBulletEntitys {
		vx, vy := getPlayerMoveInfo(isleft, isright, isup, isdown)
		pb.Move(pb, vx, vy)
	}

	if isshot && gs.playerBulletStartCount%5 == 0 {
		gs.playerBulletStartCount += 1
		bulletTexture := common.GetTexture("textures/bullet3.png")
		bb := entitys.NewBulletBuilder()
		bb.SetDrawable(bulletTexture)
		bb.SetVirtualPosition(gs.playerEntity.VirtualPosition())
		bb.SetSpeed(16)
		bb.SetZIndex(10)

		b := bb.Build()
		bulletEntity = append(bulletEntity, &b)
		gs.playerBulletEntitys[b.ID()] = &b
	} else if !isshot {
		gs.playerBulletStartCount = 0
	} else {
		gs.playerBulletStartCount += 1

	}

	// EnemyBullet Upate
	for _, eb := range gs.enemyBulletEntitys {
		s := common.NewSetting()
		virtualPosition := eb.VirtualPosition()
		mergin := eb.Mergin()
		if (virtualPosition.X < -1*mergin.X || s.GetGameAreaSize().X+mergin.X < virtualPosition.X) || (virtualPosition.Y < -1*mergin.Y || s.GetGameAreaSize().Y+mergin.Y < virtualPosition.Y) {
			gs.Remove(eb.BasicEntity())
		}
	}

	if false {
		for i := range gs.enemyEntitys {
			e := gs.enemyEntitys[i]
			gs.enemyBulletCount += 1
			for i := 0; i < 4; i++ {
				bulletTexture := common.GetTexture("textures/bullet2.png")
				bb := entitys.NewEnemyBuilder()
				bb.SetDrawable(bulletTexture)
				bb.SetVirtualPosition(e.VirtualPosition())
				bb.SetSpeed(4)
				bb.SetAngle(float32(float32(gs.framecount%72*5) + float32(90*i) + float32(gs.enemyBulletCount)))
				bb.SetZIndex(10)
				// Angle(
				b := bb.Build()
				bulletEntity = append(bulletEntity, &b)
				gs.enemyBulletEntitys[b.ID()] = &b
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

func getPlayerMoveInfo(isleft, isright, isup, isdown bool) (vx, vy float32) {
	vx = 0
	vy = 0
	if isleft && !isright {
		vx = -1
	} else if !isleft && isright {
		vx = 1
	}
	if isup && !isdown {
		vy = -1
	} else if !isup && isdown {
		vy = 1
	}
	return vx, vy
}
