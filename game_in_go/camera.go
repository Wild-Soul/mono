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
