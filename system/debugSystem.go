package system

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
	"github.com/atoyr/goshooting/entitys"
)

type DebugSystem struct {
	world        *ecs.World
	renderSystem *engoCommon.RenderSystem
}

func (ds *DebugSystem) New(w *ecs.World) {
	// debug message
	fmt.Printf("Canvas Width:%f Height:%f Scale:%f \n", engo.CanvasWidth(), engo.CanvasHeight(), engo.CanvasScale())
	fmt.Printf("Window Width:%f Height:%f  \n", engo.WindowWidth(), engo.WindowHeight())

	ds.world = w

	c := entitys.NewCharacterBuilder(common

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *engoCommon.RenderSystem:
			ds.renderSystem = sys
		}
	}
}

func (ds *DebugSystem) Update(dt float32) {

}

func (ds *DebugSystem) Remove(b ecs.BasicEntity) {

}
