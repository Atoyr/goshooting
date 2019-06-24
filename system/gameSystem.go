package system

import (
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

	entityList []entitys.Modeler

	framecount             uint64
	playerEntityID         uint64
	playerBulletIDs        []uint64
	enemyIDs               []uint64
	enemyBulletIDs         []uint64
	playerBulletStartCount uint64
	enemyBulletCount       uint64

	builderCollection map[string]entitys.Builder
}

// New is Startup Entity
func (gs *GameSystem) New(w *ecs.World) {
	gs.world = w
	gs.framecount = 0
	gs.playerBulletStartCount = 0
	gs.enemyBulletCount = 0
	gs.entityList = make([]entitys.Modeler, 131072)

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
	playerBuilder.Speed = 8
	playerBuilder.SetZIndex(10)
	playerBuilder.SetHitPoint(1)

	// enemy
	enemyBuilder := entitys.NewEnemyBuilder()
	enemyBuilder.SetDrawable(enemyTexture)
	enemyBuilder.SetPosition(engo.Point{X: 0, Y: -100})
	enemyBuilder.SetHitPoint(100)

	// Regist Entity

	player := playerBuilder.Build()
	gs.addModeler(player)

	enemy := enemyBuilder.Build()
	gs.addModeler(enemy)

	player.AddedRenderSystem(gs.renderSystem)
	enemy.AddedRenderSystem(gs.renderSystem)
}

// Update is Frame Update
func (gs *GameSystem) Update(dt float32) {
	gs.framecount += 1
	if gs.framecount == math.MaxInt64 {
		gs.framecount = 0
	}
	// get inputs
	// Get Input
	//	isleft := engo.Input.Button("MoveLeft").Down()
	//	isright := engo.Input.Button("MoveRight").Down()
	//	isup := engo.Input.Button("MoveUp").Down()
	//	isdown := engo.Input.Button("MoveDown").Down()
	//	islowspeed := engo.Input.Button("LowSpeed").Down()
	// isshot := engo.Input.Button("Shot").Down()

	// setSpeed
	// speed := float32(0)
	// if islowspeed {
	// speed = 4
	// } else {
	// speed = 8
	// }

	// Get vx vy
	//	vx := float32(0)
	//	vy := float32(0)
	//	if isleft && !isright {
	//		vx = -1
	//	} else if !isleft && isright {
	//		vx = 1
	//	}
	//	if isup && !isdown {
	//		vy = -1
	//	} else if !isup && isdown {
	//		vy = 1
	//	}
	//
	// Move Player
	// gs.playerEntity.Move(vx, vy, speed)

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
	//	if isshot {
	//		gs.playerEntity.Attack(gs.playerEntity, dt)
	//	} else if !isshot {
	//		gs.playerEntity.AttackStartFrame = -1
	//	}
	//
	// Collision
	//	if gs.framecount%60 == 0 {
	//		for _, e := range gs.enemyEntitys {
	//			if gs.playerEntity.IsCollision(e.EntityModel) {
	//				fmt.Println("Collision!!")
	//			}
	//		}
	//	}
	//	for _, e := range gs.enemyEntitys {
	//		e.Update(dt)
	//		for _, pb := range gs.playerBulletEntitys {
	//			if !pb.Hidden() {
	//				if e.IsCollision(pb) {
	//					pbHP := pb.HitPoint()
	//					eHp := e.HitPoint()
	//					e.AddHitPoint(-1 * pbHP)
	//					pb.AddHitPoint(-1 * eHp)
	//					if e.HitPoint() < 0 {
	//						e.SetHidden(true)
	//					}
	//					if pb.HitPoint() < 0 {
	//						pb.SetHidden(true)
	//					}
	//				}
	//			}
	//		}
	//	}

	// EnemyBullet Upate
	//	for _, eb := range gs.enemyBulletEntitys {
	//		//		if !eb.IsOverGameArea() {
	//		//			eb.Update(dt)
	//		//		}
	//		//		if eb.IsOverGameArea() && !eb.Hidden() {
	//		//			eb.SetHidden(true)
	//		//		}
	//	}
}

func (gs *GameSystem) Remove(b ecs.BasicEntity) {
	for _, system := range gs.world.Systems() {
		switch sys := system.(type) {
		case *engoCommon.RenderSystem:
			sys.Remove(b)
		}
	}
}

func (gs *GameSystem) addModeler(m entitys.Modeler) {
	id := int(m.ID())
	if c := cap(gs.entityList); c < id {
		l := make([]entitys.Modeler, id-c-1)
		gs.entityList = append(gs.entityList, l...)
	}
	gs.entityList[id] = m
}
