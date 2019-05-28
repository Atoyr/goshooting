package common

import (
	"errors"

	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
)

type NumberSize string

const Number_8_48 NumberSize = "number_8_48"
const Number_16_16 NumberSize = "number_16_16"
const Number_16_32 NumberSize = "number_16_32"

var number_8_48_spritsheet *engoCommon.Spritesheet
var number_16_16_spritsheet *engoCommon.Spritesheet
var number_16_32_spritsheet *engoCommon.Spritesheet

func InitializeNumber(numberSize NumberSize, url string) {
	var s *engoCommon.Spritesheet
	switch numberSize {
	case Number_8_48:
		s = engoCommon.NewSpritesheetFromFile(url, 8, 48)
		number_8_48_spritsheet = s
	case Number_16_16:
		s = engoCommon.NewSpritesheetFromFile(url, 16, 16)
		number_16_16_spritsheet = s
	case Number_16_32:
		s = engoCommon.NewSpritesheetFromFile(url, 16, 32)
		number_16_32_spritsheet = s
	default:
		return
	}

	tc := NewTextureContainer()
	for i := 0; i < s.CellCount(); i++ {
		str := string(numberSize) + string(i)
		t := s.Cell(i)
		tc[str] = &t
	}
}

func isNumberInitialize(numSize NumberSize) error {
	switch numSize {
	case Number_8_48:
		if number_8_48_spritsheet == nil {
			return errors.New("Not Initialize ")
		}
	case Number_16_16:
		if number_16_16_spritsheet == nil {
			return errors.New("Not Initialize ")
		}
	case Number_16_32:
		if number_16_32_spritsheet == nil {
			return errors.New("Not Initialize ")
		}
	}
	return nil
}

func GetNumberSize(numSize NumberSize) engo.Point {
	p := engo.Point{X: 0, Y: 0}
	switch numSize {
	case Number_8_48:
		p.Add(engo.Point{X: 8, Y: 48})
	case Number_16_16:
		p.Add(engo.Point{X: 16, Y: 16})
	case Number_16_32:
		p.Add(engo.Point{X: 16, Y: 32})
	}
	return p
}

func GetNumberTexture(value int, numSize NumberSize) (*engoCommon.Texture, error) {
	if err := isNumberInitialize(numSize); err != nil {
		return nil, err
	}
	return NewTextureContainer()[string(numSize)+string(value)], nil
}

func GetNumberTextures(numSize NumberSize) ([]engoCommon.Texture, error) {
	textures := make([]engoCommon.Texture, 10, 10)
	for i := 0; i < 10; i++ {
		t, err := GetNumberTexture(i, numSize)
		if err != nil {
			return nil, err
		}
		textures[i] = *t
	}
	return textures, nil
}

func GetNumberAnimationComponent(numSize NumberSize, rate float32) (engoCommon.AnimationComponent, error) {
	if err := isNumberInitialize(numSize); err != nil {
		return engoCommon.AnimationComponent{}, err
	}
	drawables := make([]engoCommon.Drawable, 10, 10)
	for i := 0; i < 10; i++ {
		drawables[i] = NewTextureContainer()[string(numSize)+string(i)]
	}
	return engoCommon.NewAnimationComponent(drawables, rate), nil
}
