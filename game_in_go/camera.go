package main

type Camera struct {
	X float64
	Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) CenterOn(x, y, screendWidth, screenHeight float64) {
	c.X = -x + screenHeight/2.0
	c.Y = -y + screenHeight/2.0
}

func (c *Camera) Constraint(tilemapWidthPixels, tilemapHeightPixels, screenWidth, screenHeight float64) {
	c.X = min(c.X, 0.0)
	c.Y = min(c.Y, 0.0)

	c.X = max(c.X, screenWidth-tilemapWidthPixels)
	c.Y = max(c.Y, screenHeight-tilemapHeightPixels)
}
