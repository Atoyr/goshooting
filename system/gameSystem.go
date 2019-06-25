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
	enemybulletTexture := common.GetTexture("textures/bullet2.png")
	bulletTexture := common.GetTexture("textures/bullet3.png")

	// Create Entity
	// Player
	playerBuilder := entitys.NewPlayerBuilder()
	playerBuilder.SetDrawable(playerTexture)
	playerBuilder.SetPosition(engo.Point{X: 0, Y: 100})
	playerBuilder.Speed = 8
	playerBuilder.LowSpeed = 4
	playerBuilder.SetZIndex(200)
	playerBuilder.SetHitPoint(1)

	bb := entitys.NewBulletBuilder()
	bb.SetDrawable(bulletTexture)
	bb.Speed = 16
	bb.SetZIndex(10)
	bb.SetHitPoint(10)
	playerBuilder.AttackBuilderList = append(playerBuilder.AttackBuilderList, &bb)
	if bb2, ok := bb.Clone().(*entitys.BulletBuilder); ok {
		bb2.SetRotation(30)
		playerBuilder.AttackBuilderList = append(playerBuilder.AttackBuilderList, bb2)
	}
	if bb3, ok := bb.Clone().(*entitys.BulletBuilder); ok {
		bb3.SetRotation(-30)
		playerBuilder.AttackBuilderList = append(playerBuilder.AttackBuilderList, bb3)
	}

	playerBuilder.Attack = func(modeler entitys.Modeler, frame uint64) []entitys.Modeler {
		modelers := make([]entitys.Modeler, 0)

		if p, ok := modeler.(entitys.Player); ok {

			for _, bb := range p.AttackBuilderList {
				b := bb.Build()
				b.SetPosition(p.Position())
				modelers = append(modelers, b)
			}
		}
		return modelers
	}

	// enemy
	enemyBuilder := entitys.NewEnemyBuilder()
	enemyBuilder.SetDrawable(enemyTexture)
	enemyBuilder.SetPosition(engo.Point{X: 0, Y: -100})
	enemyBuilder.SetHitPoint(100)
	eb := entitys.NewBulletBuilder()
	eb.SetDrawable(enemybulletTexture)
	eb.Speed = 16
	eb.SetZIndex(10)
	eb.SetHitPoint(10)
	enemyBuilder.AttackBuilderList = append(enemyBuilder.AttackBuilderList, &eb)
	enemyBuilder.Attack = func(modeler entitys.Modeler, frame uint64) []entitys.Modeler {
		modelers := make([]entitys.Modeler, 0)

		if e, ok := modeler.(entitys.Enemy); ok {

			for _, ab := range e.AttackBuilderList {
				if x, ok := ab.(*entitys.BulletBuilder); ok {
					x.AddRotation(10)
					b := x.Build()
					b.SetPosition(e.Position())
					modelers = append(modelers, b)
				}
			}
		}
		return modelers
	}

	// Regist Entity

	player := playerBuilder.Build()
	gs.playerEntityID = player.ID()
	gs.addModeler(player)

	enemy := enemyBuilder.Build()
	gs.enemyIDs = append(gs.enemyIDs, enemy.ID())
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

	// Player action
	if p, ok := gs.entityList[gs.playerEntityID].(entitys.Player); ok {
		// Move Player
		v := p.Vector(isleft, isright, isup, isdown, islowspeed)
		p.AddPosition(v)

		// Move Player Bullet
		for _, child := range p.BasicEntity().Children() {
			if b, ok := gs.entityList[child.ID()].(entitys.Bullet); ok {
				b.Move()
			}
		}

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

	// Enemy action
	for _, enemyid := range gs.enemyIDs {
		if e, ok := gs.entityList[enemyid].(entitys.Enemy); ok {
			// Move Enemy

			// Move Enemy Bullet
			for _, child := range e.BasicEntity().Children() {
				if b, ok := gs.entityList[child.ID()].(entitys.Bullet); ok {
					b.Move()
				}
			}

			// Attack for Enemy
			modelers := e.Attack(e, gs.framecount)
			for _, m := range modelers {
				m.AddedRenderSystem(gs.renderSystem)
				e.AppendChild(m.BasicEntity())
				gs.addModeler(m)
			}
		}
	}

	// Collision
	// TODO : go chan

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
	gs.entityList[id] = m
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
