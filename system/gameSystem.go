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

	entityList []*Modeler

	framecount             uint64
	playerEntityID         uint64
	playerBulletIDs        uint64
	enemyIDs               uint64
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
	playerBuilder.SetPosition(engo.Point{X: 0, Y: 100})
	playerBuilder.SetSpeed(8)
	playerBuilder.SetCollisionDetectionSize(4)
	playerBuilder.SetZIndex(10)
	playerBuilder.RenderCollisionDetection(true)
	playerBuilder.SetCollisionBasicEntity(ecs.NewBasic())
	playerBuilder.SetHitPoint(1)

	player := playerBuilder.Build()

	player.Attack = func(e *entitys.Entity, frame float32) {
		if e.AttackStartFrame < 0 {
			e.AttackStartFrame = 0
		} else {
			e.AttackStartFrame++
		}
		if int(e.AttackStartFrame)%5 == 0 {
			bulletTexture := common.GetTexture("textures/bullet3.png")
			bb := entitys.NewBulletBuilder()
			bb.SetDrawable(bulletTexture)
			bb.SetPosition(gs.playerEntity.Position())
			bb.SetSpeed(16)
			bb.SetZIndex(10)
			bb.SetCollisionDetectionSize(16)
			bb.SetHitPoint(10)

			b := bb.Build()
			b.AddedRenderSystem(gs.renderSystem)
			gs.playerBulletEntitys[b.ID()] = &b
		} else {
			gs.playerBulletStartCount += 1
		}
	}

	playerBullets := map[uint64]*entitys.Entity{}
	enemys := map[uint64]*entitys.Entity{}
	enemyBullets := map[uint64]*entitys.Entity{}

	// enemy
	enemyBuilder := entitys.NewEnemyBuilder()
	enemyBuilder.SetDrawable(enemyTexture)
	enemyBuilder.SetPosition(engo.Point{X: 0, Y: -100})
	enemyBuilder.SetAngle(70)
	enemyBuilder.SetCollisionDetectionSize(16)
	enemyBuilder.SetHitPoint(100)
	enemyBuilder.RenderCollisionDetection(true)
	enemyBuilder.SetCollisionBasicEntity(ecs.NewBasic())
	enemy := enemyBuilder.Build()
	enemy.Attack = func(e *entitys.Entity, frame float32) {
		if int(e.AttackFrame)%5 == 0 {
			bulletTexture := common.GetTexture("textures/bullet2.png")

			bb := entitys.NewBulletBuilder()
			bb.SetPosition(e.Position())
			bb.SetSpeed(4)
			bb.SetZIndex(10)
			bb.SetHitPoint(1)
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
	// Get Input
	isleft := engo.Input.Button("MoveLeft").Down()
	isright := engo.Input.Button("MoveRight").Down()
	isup := engo.Input.Button("MoveUp").Down()
	isdown := engo.Input.Button("MoveDown").Down()
	islowspeed := engo.Input.Button("LowSpeed").Down()
	isshot := engo.Input.Button("Shot").Down()

	// setSpeed
	speed := float32(0)
	if islowspeed {
		speed = 4
	} else {
		speed = 8
	}

	// Get vx vy
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

	// Move Player
	gs.playerEntity.Move(vx, vy, speed)

	// Move PlayerBullet and hidden PlayerBullet
	//	for _, pb := range gs.playerBulletEntitys {
	//		//		if !pb.IsOverGameArea() {
	//		//			pb.Update(dt)
	//		//		}
	//		//		if pb.IsOverGameArea() && !pb.Hidden() {
	//		//			pb.SetHidden(true)
	//		//		}
	//	}

	// Attack for Player
	if isshot {
		gs.playerEntity.Attack(gs.playerEntity, dt)
	} else if !isshot {
		gs.playerEntity.AttackStartFrame = -1
	}

	// Collision
	if gs.framecount%60 == 0 {
		for _, e := range gs.enemyEntitys {
			if gs.playerEntity.IsCollision(e.EntityModel) {
				fmt.Println("Collision!!")
			}
		}
	}
	for _, e := range gs.enemyEntitys {
		e.Update(dt)
		for _, pb := range gs.playerBulletEntitys {
			if !pb.Hidden() {
				if e.IsCollision(pb) {
					pbHP := pb.HitPoint()
					eHp := e.HitPoint()
					e.AddHitPoint(-1 * pbHP)
					pb.AddHitPoint(-1 * eHp)
					if e.HitPoint() < 0 {
						e.SetHidden(true)
					}
					if pb.HitPoint() < 0 {
						pb.SetHidden(true)
					}
				}
			}
		}
	}

	// EnemyBullet Upate
	//	for _, eb := range gs.enemyBulletEntitys {
	//		//		if !eb.IsOverGameArea() {
	//		//			eb.Update(dt)
	//		//		}
	//		//		if eb.IsOverGameArea() && !eb.Hidden() {
	//		//			eb.SetHidden(true)
	//		//		}
	//	}

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
