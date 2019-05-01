package common

import (
	"sync"

	"github.com/EngoEngine/engo"
)

type Setting struct {
	canvas             engo.Point
	renderPositionRate engo.Point
	gameAreaSize       engo.Point
	renderScale        float32
	renderCanvas       engo.Point
}

var (
	setting     *Setting
	settingOnce sync.Once
)

func NewSetting() *Setting {
	settingOnce.Do(func() {
		canvas := engo.Point{X: 1280, Y: 720}
		renderPositionRate := engo.Point{X: 0.5, Y: 0.5}
		gameAreaSize := engo.Point{X: 516, Y: 688}
		renderScale := float32(1)
		setting = &Setting{
			canvas:             canvas,
			renderPositionRate: renderPositionRate,
			gameAreaSize:       gameAreaSize,
			renderScale:        renderScale,
		}

		setting.UpdateCanvas()
	})
	return setting
}

func (s *Setting) UpdateCanvas() engo.Point {
	xy := engo.Point{X: engo.CanvasWidth(), Y: engo.CanvasHeight()}
	s.renderCanvas = xy
	s.renderScale = xy.Y / s.canvas.Y
	return xy
}

func (s *Setting) AABB() engo.AABB {
	center := engo.Point{X: s.canvas.X * s.renderPositionRate.X, Y: s.canvas.Y * s.renderPositionRate.Y}
	half := engo.Point{X: s.gameAreaSize.X * 0.5, Y: s.gameAreaSize.Y * 0.5}
	return engo.AABB{Min: engo.Point{X: center.X - half.X, Y: center.Y - half.Y}, Max: engo.Point{X: center.X + half.X, Y: center.Y + half.Y}}
}

func (s *Setting) GetGameAreaSize() engo.Point {
	return s.gameAreaSize
}

func (s *Setting) GetCanvas() engo.Point {
	return s.canvas
}

func (s *Setting) ConvertRenderPosition(xy engo.Point) engo.Point {
	ret := engo.Point{X: 0, Y: 0}
	aabb := s.AABB()
	ret = *ret.Add(*aabb.Min.MultiplyScalar(s.renderScale))
	ret = *ret.Add(*xy.MultiplyScalar(s.renderScale))

	return ret
}
