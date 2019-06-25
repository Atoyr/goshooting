package system

import (
	"fmt"
	"math"
	"reflect"

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

	framecount     uint64
	playerEntityID uint64
	enemyIDs       []uint64

	builderCollection map[string]entitys.Builder
}

// New is Startup Entity
func (gs *GameSystem) New(w *ecs.World) {
	gs.world = w
	gs.framecount = 0
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
	playerBuilder.LowSpeed = 4
	playerBuilder.SetZIndex(10)
	playerBuilder.SetHitPoint(1)

	playerBuilder.Attack = func(modeler entitys.Modeler, frame uint64) []entitys.Modeler {
		modelers := make([]entitys.Modeler, 0)

		bulletTexture := common.GetTexture("textures/bullet3.png")
		bb := entitys.NewBulletBuilder()
		bb.SetDrawable(bulletTexture)
		bb.SetPosition(modeler.Position())
		bb.SetSpeed(16)
		bb.SetZIndex(10)
		bb.SetHitPoint(10)

		b := bb.Build()
		modelers = append(modelers, b)
		return modelers
	}

	// enemy
	enemyBuilder := entitys.NewEnemyBuilder()
	enemyBuilder.SetDrawable(enemyTexture)
	enemyBuilder.SetPosition(engo.Point{X: 0, Y: -100})
	enemyBuilder.SetHitPoint(100)

	// Regist Entity

	player := playerBuilder.Build()
	gs.playerEntityID = player.ID()
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

	// Get Input
	isleft := engo.Input.Button("MoveLeft").Down()
	isright := engo.Input.Button("MoveRight").Down()
	isup := engo.Input.Button("MoveUp").Down()
	isdown := engo.Input.Button("MoveDown").Down()
	islowspeed := engo.Input.Button("LowSpeed").Down()
	isshot := engo.Input.Button("Shot").Down()

	// Move Player
	if p, ok := gs.entityList[gs.playerEntityID].(entitys.Player); ok {
		v := p.Vector(isleft, isright, isup, isdown, islowspeed)
		p.AddPosition(v)

		// Attack for Player
		if isshot {
			modelers := p.Attack(p, gs.framecount)
			for _, m := range modelers {
				m.AddedRenderSystem(gs.renderSystem)
				p.AppendChild(m.BasicEntity())
				gs.addModeler(m)
			}
		}
	}

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

func (gs *GameSystem) addModeler(m entitys.Modeler) error {
	id := 0
	if tmpid := m.ID(); tmpid > math.MaxInt32 {
		return fmt.Errorf("ID is out of index : %d", m.ID())
	} else {
		id = int(tmpid)
	}

	if c := cap(gs.entityList); c < id {
		l := make([]entitys.Modeler, id-c-1)
		gs.entityList = append(gs.entityList, l...)
	}
	fmt.Printf("%d %s \n", id, reflect.TypeOf(m))
	gs.entityList[id] = m
	fmt.Printf("%d %s \n", id, reflect.TypeOf(gs.entityList[id]))
	return nil
}

func (gs *GameSystem) removeModeler(id uint64) (error, entitys.Modeler) {
	if id > math.MaxInt32 {
		return fmt.Errorf("ID is out of index : %d", id), nil
	}
	if i := int(id); cap(gs.entityList) < i {
		return fmt.Errorf("ID is out of index : %d", id), nil
	}
	e := gs.entityList[id]
	gs.entityList[id] = nil
	return nil, e
}
