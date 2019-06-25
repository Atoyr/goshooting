package common

import (
	"fmt"
	"sync"

	"github.com/EngoEngine/engo"
)

type Setting struct {
	// Base Setting
	renderPositionRate engo.Point
	// Convert Render
	renderScale  float32
	renderCanvas engo.Point
}

var (
	setting     *Setting
	settingOnce sync.Once
)

var baseCanvasHeight float32 = float32(846)
var baseGameAreaSize engo.Point = engo.Point{X: 600, Y: 800}

// NewSetting is create setting at once and return setting
func NewSetting() *Setting {
	settingOnce.Do(func() {
		renderPositionRate := engo.Point{X: 0.5, Y: 0.5}
		renderScale := float32(1)
		setting = &Setting{
			renderPositionRate: renderPositionRate,
			renderScale:        renderScale,
		}

		setting.UpdateRenderParams()
	})
	return setting
}

// UpdateCanvas is Update render Scale for CanvasWidth
func (s *Setting) UpdateRenderParams() engo.Point {
	xy := engo.Point{X: engo.CanvasWidth(), Y: engo.CanvasHeight()}
	s.renderCanvas = xy
	fmt.Println(s.renderScale)
	s.renderScale = xy.Y / baseCanvasHeight
	fmt.Println(s.renderScale)
	return xy
}

// AABB is return gameArea min and max position
func (s *Setting) AABB() engo.AABB {
	min := engo.Point{X: baseGameAreaSize.X, Y: baseGameAreaSize.Y}
	max := engo.Point{X: baseGameAreaSize.X, Y: baseGameAreaSize.Y}
	min.MultiplyScalar(-0.5)
	max.MultiplyScalar(0.5)
	return engo.AABB{Min: s.ConvertVirtualPositionToRenderPosition(min), Max: s.ConvertVirtualPositionToRenderPosition(max)}
}

func (s *Setting) RenderCanvas() engo.Point {
	return s.renderCanvas
}

func (s *Setting) GameAreaSize() engo.Point {
	return baseGameAreaSize
}

func (s *Setting) RenderGameAreaSize() engo.Point {
	gameArea := engo.Point{X: baseGameAreaSize.X, Y: baseGameAreaSize.Y}
	gameArea.MultiplyScalar(s.renderScale)
	return gameArea
}

func (s *Setting) Scale() engo.Point {
	ret := engo.Point{X: 0, Y: 0}
	ret.AddScalar(s.renderScale)
	return ret
}

func (s *Setting) ConvertVirtualPositionToRenderPosition(xy engo.Point) engo.Point {
	ret := engo.Point{X: 0, Y: 0}
	pxy := engo.Point{X: xy.X, Y: xy.Y}
	pxy.MultiplyScalar(s.renderScale)
	ret.Add(s.renderCanvas)
	ret.Multiply(s.renderPositionRate)
	ret.Add(pxy)

	return ret
}
