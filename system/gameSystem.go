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
	gs.world = w
	gs.framecount = 0
	gs.playerBulletStartCount = 0

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engoCommon.RenderSystem:
			gs.renderSystem = sys
		}
	}
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
	playerBuilder.DenyOverArea = true
	playerBuilder.RenderCollisionDetection(true)
	playerBuilder.SetCollisionBasicEntity(ecs.NewBasic())

	player := playerBuilder.Build()

	player.Attack = func(e *entitys.Entity, frame float32) {
		isshot := engo.Input.Button("Shot").Down()
		if isshot {
			if e.AttackStartFrame < 0 {
				e.AttackStartFrame = 0
			} else {
				e.AttackStartFrame++
			}
			if int(e.AttackStartFrame)%5 == 0 {
				bulletTexture := common.GetTexture("textures/bullet3.png")
				bb := entitys.NewBulletBuilder()
				bb.SetDrawable(bulletTexture)
				bb.SetVirtualPosition(gs.playerEntity.VirtualPosition())
				bb.SetSpeed(16)
				bb.SetZIndex(10)

				b := bb.Build()
				b.AddedRenderSystem(gs.renderSystem)
				gs.playerBulletEntitys[b.ID()] = &b
			} else if !isshot {
				e.AttackStartFrame = -1
			} else {
				gs.playerBulletStartCount += 1
			}
		}
	}

	playerBullets := map[uint64]*entitys.Entity{}
	enemys := map[uint64]*entitys.Entity{}
	enemyBullets := map[uint64]*entitys.Entity{}

	// enemy
	enemyBuilder := entitys.NewEnemyBuilder()
	enemyBuilder.SetDrawable(enemyTexture)
	enemyBuilder.SetVirtualPosition(engo.Point{X: 0, Y: -100})
	enemyBuilder.SetAngle(70)
	enemyBuilder.SetCollisionDetectionSize(16)
	enemyBuilder.RenderCollisionDetection(true)
	enemyBuilder.SetCollisionBasicEntity(ecs.NewBasic())
	enemy := enemyBuilder.Build()
	enemy.Attack = func(e *entitys.Entity, frame float32) {
		if int(e.AttackFrame)%5 == 0 {
			bulletTexture := common.GetTexture("textures/bullet2.png")

			bb := entitys.NewBulletBuilder()
			bb.SetVirtualPosition(e.VirtualPosition())
			bb.SetSpeed(4)
			bb.SetZIndex(10)
			for i := 0; i < 4; i++ {
				b := bb.Build()
				b.SetDrawable(bulletTexture)
				b.Angle = float32(90*i) + e.AttackFrame
				b.AddedRenderSystem(gs.renderSystem)
				gs.enemyBulletEntitys[b.ID()] = &b
			}
		}
		e.AttackFrame++
	}
	enemys[enemy.ID()] = &enemy

	// Regist Entity
	gs.playerEntity = &player
	gs.playerBulletEntitys = playerBullets
	gs.enemyEntitys = enemys
	gs.enemyBulletEntitys = enemyBullets
	gs.enemyBulletCount = 0

	player.AddedRenderSystem(gs.renderSystem)
	player.AddedRenderSystemToCollisionComponent(gs.renderSystem)
	enemy.AddedRenderSystem(gs.renderSystem)
	enemy.AddedRenderSystemToCollisionComponent(gs.renderSystem)
}

// Update is Frame Update
func (gs *GameSystem) Update(dt float32) {
	gs.framecount += 1
	if gs.framecount == math.MaxInt64 {
		gs.framecount = 0
	}
	// get inputs
	isleft := engo.Input.Button("MoveLeft").Down()
	isright := engo.Input.Button("MoveRight").Down()
	isup := engo.Input.Button("MoveUp").Down()
	isdown := engo.Input.Button("MoveDown").Down()
	islowspeed := engo.Input.Button("LowSpeed").Down()

	speed := float32(0)

	// setSpeed
	if islowspeed {
		speed = 4
	} else {
		speed = 8
	}

	vx := float32(0)
	vy := float32(0)
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

	gs.playerEntity.Move(vx, vy, speed)

	// Collision
	if gs.framecount%60 == 0 {
		for _, e := range gs.enemyEntitys {
			if gs.playerEntity.IsCollision(e) {
				fmt.Println("Collision!!")
			}
		}
	}

	// PlayerBullet Update
	for _, pb := range gs.playerBulletEntitys {
		pb.Update(dt)
	}

	// EnemyBullet Upate
	for _, eb := range gs.enemyBulletEntitys {
		eb.Update(dt)
		if eb.IsOverGameArea() {
			gs.Remove(eb.BasicEntity())
		}
	}

	for _, e := range gs.enemyEntitys {
		e.Update(dt)
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
