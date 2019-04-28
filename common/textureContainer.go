package common

import (
	"fmt"
	"sync"

	"github.com/EngoEngine/engo/common"
	engoCommon "github.com/EngoEngine/engo/common"
)

var (
	textureContainer map[string]*engoCommon.Texture
	textureOnce      sync.Once
)

func NewTextureContainer() map[string]*engoCommon.Texture {
	textureOnce.Do(func() {
		textureContainer = map[string]*engoCommon.Texture{}
	})
	return textureContainer
}

func GetTexture(url string) *engoCommon.Texture {
	if t, ok := textureContainer[url]; ok {
		return t
	} else {
		t, err := common.LoadedSprite(url)
		if err != nil {
			fmt.Println("Unable to load texture: " + err.Error())
		} else {
			tc := NewTextureContainer()
			tc[url] = t
		}
		return t
	}
}
