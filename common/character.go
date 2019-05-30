package common

import (
	"errors"

	"github.com/EngoEngine/engo"
	engoCommon "github.com/EngoEngine/engo/common"
)

type CharacterSize string

// Character_16_16 is character graphic 16"16
const Character_16_16 CharacterSize = "Character_16_16"

// character_16_16Spritsheet is character graphic 16"16 sprit sheet
var character_16_16_Spritsheet *engoCommon.Spritesheet

func InitializeCharacter(characterSize CharacterSize, url string) {
	var s *engoCommon.Spritesheet
	switch characterSize {
	case Character_16_16:
		s = engoCommon.NewSpritesheetFromFile(url, 16, 16)
		character_16_16_Spritsheet = s
	default:
		return
	}

	tc := NewTextureContainer()

	// TODO : i -> char
	for i := 0; i < s.CellCount(); i++ {
		str := string(characterSize) + string(i)
		t := s.Cell(i)
		tc[str] = &t
	}
}

func isCharacterInitialize(characterSize CharacterSize) error {
	switch characterSize {
	case Character_16_16:
		if character_16_16_Spritsheet == nil {
			return errors.New("Not Initialize ")
		}
	}
	return nil
}

func GetCharacterSize(characterSize CharacterSize) engo.Point {
	p := engo.Point{X: 0, Y: 0}
	switch characterSize {
	case Character_16_16:
		p.Add(engo.Point{X: 16, Y: 16})
	}
	return p
}

func GetCharacterTexture(char string, characterSize CharacterSize) (*engoCommon.Texture, error) {
	if err := isCharacterInitialize(characterSize); err != nil {
		return nil, err
	}
	return NewTextureContainer()[string(characterSize)+string(getCharacterIndex(char))], nil
}

func GetCharacterTextures(str string, characterSize CharacterSize) ([]engoCommon.Texture, error) {
	textures := make([]engoCommon.Texture, len(str))
	for i := 0; i < len(str); i++ {
		r := str[i]
		t, err := GetCharacterTexture(string(r), characterSize)
		if err != nil {
			return nil, err
		}
		textures[i] = *t
	}
	return textures, nil
}

func getCharacterIndex(char string) int {
	ret := 0
	switch char {
	case "0":
		ret = 0
	case "1":
		ret = 1
	case "2":
		ret = 2
	case "3":
		ret = 3
	case "4":
		ret = 4
	case "5":
		ret = 5
	case "6":
		ret = 6
	case "7":
		ret = 7
	case "8":
		ret = 8
	case "9":
		ret = 9
	case "s0":
		ret = 10
	case "s1":
		ret = 11
	case "s2":
		ret = 12
	case "s3":
		ret = 13
	case "s4":
		ret = 14
	case "s5":
		ret = 15
	case "s6":
		ret = 16
	case "s7":
		ret = 17
	case "s8":
		ret = 18
	case "s9":
		ret = 19
	case "A":
		ret = 26
	case "B":
		ret = 27
	case "C":
		ret = 28
	case "D":
		ret = 29
	case "E":
		ret = 30
	case "F":
		ret = 31
	case "G":
		ret = 32
	case "H":
		ret = 33
	case "I":
		ret = 34
	case "J":
		ret = 35
	case "K":
		ret = 36
	case "L":
		ret = 37
	case "M":
		ret = 38
	case "N":
		ret = 39
	case "O":
		ret = 40
	case "P":
		ret = 41
	case "Q":
		ret = 42
	case "R":
		ret = 43
	case "S":
		ret = 44
	case "T":
		ret = 45
	case "U":
		ret = 46
	case "V":
		ret = 47
	case "W":
		ret = 48
	case "X":
		ret = 49
	case "Y":
		ret = 50
	case "Z":
		ret = 51
	case "a":
		ret = 52
	case "b":
		ret = 53
	case "c":
		ret = 54
	case "d":
		ret = 55
	case "e":
		ret = 56
	case "f":
		ret = 57
	case "g":
		ret = 58
	case "h":
		ret = 59
	case "i":
		ret = 60
	case "j":
		ret = 61
	case "k":
		ret = 62
	case "l":
		ret = 63
	case "m":
		ret = 64
	case "n":
		ret = 65
	case "o":
		ret = 66
	case "p":
		ret = 67
	case "q":
		ret = 68
	case "r":
		ret = 69
	case "s":
		ret = 70
	case "t":
		ret = 71
	case "u":
		ret = 72
	case "v":
		ret = 73
	case "w":
		ret = 74
	case "x":
		ret = 75
	case "y":
		ret = 76
	case "z":
		ret = 77
	case ".":
		ret = 78
	case ",":
		ret = 79
	case "-":
		ret = 80
	case "+":
		ret = 81
	case "//":
		ret = 82
	case "**":
		ret = 83
	case "=":
		ret = 84
	case "/":
		ret = 85
	case "^2":
		ret = 86
	case "^3":
		ret = 87
	case "'":
		ret = 88
	case "r\"":
		ret = 89
	case "l\"":
		ret = 90
	case "!":
		ret = 91
	case "?":
		ret = 92
	case ":":
		ret = 93
	case "^":
		ret = 94
	case "&":
		ret = 95
	case "(":
		ret = 96
	case ")":
		ret = 97
	case "[":
		ret = 98
	case "]":
		ret = 99
	case "_":
		ret = 100
	case "#":
		ret = 101
	case "po":
		ret = 102
	case "bl":
		ret = 103
	case " ":
		ret = 104
	}
	return ret
}
