package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Img    *ebiten.Image
	Health int16
	X      float64
	Y      float64
}
