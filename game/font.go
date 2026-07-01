package game

import (
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

var gameFont font.Face
var smallFont font.Face
var fontLoaded bool

func initFont() {
	if fontLoaded {
		return
	}
	fontLoaded = true
	fontPaths := []string{
		"C:\\Windows\\Fonts\\msyh.ttc",
		"C:\\Windows\\Fonts\\msyh.ttf",
		"C:\\Windows\\Fonts\\simhei.ttf",
		"C:\\Windows\\Fonts\\simsun.ttc",
		"resource/msyh.ttc",
		"resource/simhei.ttf",
	}
	for _, fp := range fontPaths {
		data, err := os.ReadFile(fp)
		if err != nil {
			continue
		}
		tt, err := opentype.Parse(data)
		if err != nil {
			continue
		}
		face, err := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    16,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			continue
		}
		gameFont = face
		small, err := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    12,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err == nil {
			smallFont = small
		} else {
			smallFont = face
		}
		return
	}
	gameFont = basicfont.Face7x13
	smallFont = basicfont.Face7x13
}

func fontNormal() font.Face {
	if gameFont == nil {
		initFont()
	}
	return gameFont
}

func fontSmall() font.Face {
	if smallFont == nil {
		initFont()
	}
	return smallFont
}
