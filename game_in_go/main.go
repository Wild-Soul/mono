package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img *ebiten.Image
	X   float64
	Y   float64
}

type Player struct {
	*Sprite
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Game struct {
	Player      *Sprite
	Enemies     []*Enemy
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Player.Y += 2
	}

	for _, enemy := range g.Enemies {
		if enemy.FollowsPlayer {
			if g.Player.X > enemy.X {
				enemy.X += 1
			} else if g.Player.X < enemy.X {
				enemy.X -= 1
			}

			if g.Player.Y > enemy.Y {
				enemy.Y += 1
			} else if g.Player.Y < enemy.Y {
				enemy.Y -= 1
			}

		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	drawOptions := &ebiten.DrawImageOptions{}

	// Draw background
	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {

			x := index % layer.Width
			y := index / layer.Width
			x *= 16
			y *= 16

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22
			srcX *= 16
			srcY *= 16

			drawOptions.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				drawOptions,
			)
			drawOptions.GeoM.Reset()
		}
	}

	playerImage := g.Player.Img.SubImage(image.Rect(0, 0, 16, 16))
	drawOptions.GeoM.Translate(g.Player.X, g.Player.Y)

	screen.DrawImage(playerImage.(*ebiten.Image), drawOptions)
	drawOptions.GeoM.Reset()

	for _, enemy := range g.Enemies {
		enemyImg := enemy.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image)

		drawOptions.GeoM.Translate(enemy.X, enemy.Y)
		screen.DrawImage(enemyImg, drawOptions)

		drawOptions.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImage, _, err := ebitenutil.NewImageFromFile("assets/images/ninja.png")
	if err != nil {
		log.Fatal("Failed to load player image:", err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal("Failed to load skeleton image:", err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal("Failed to load floor tileset image:", err)
	}

	tilemapJson, err := NewTilemapJSON("assets/maps/tilesets/spawn.json")
	if err != nil {
		log.Fatal("Failed to load tilemap json:", err)
	}

	game := &Game{
		Player: &Sprite{
			Img: playerImage,
			X:   100,
			Y:   300,
		},
		Enemies: []*Enemy{
			{
				Sprite: &Sprite{
					Img: skeletonImg,
					X:   100.4,
					Y:   100.4,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &Sprite{
					Img: skeletonImg,
					X:   200.4,
					Y:   200.4,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &Sprite{
					Img: skeletonImg,
					X:   300.4,
					Y:   300.4,
				},
				FollowsPlayer: false,
			},
		},
		tilemapJSON: tilemapJson,
		tilemapImg:  tilemapImg,
	}

	// Start the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
