package main

import (
	"game_in_go/entities"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 240
)

type Game struct {
	Player      *entities.Sprite
	Enemies     []*entities.Enemy
	Potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
	cam         *Camera
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

	for _, potion := range g.Potions {
		if (g.Player.X < potion.X+16) && (g.Player.X+16 > potion.X) &&
			(g.Player.Y < potion.Y+16) && (g.Player.Y+16 > potion.Y) {

			g.Player.Health = min(100, g.Player.Health+potion.AmtHeal)
			if g.Player.Health > 100 {
				g.Player.Health = 100
			}

			// Remove the potion from the game by moving it off-screen. Ideally should be removed from the slice.
			potion.X = -100
			potion.Y = -100
		}
	}

	// Center camera on player
	// fmt.Printf("Testing screenpos: %s, %s", float64(g.tilemapJSON.Layers[0].Width)*16.0, float64(g.tilemapJSON.Layers[0].Width)*16.0)
	g.cam.CenterOn(g.Player.X+8, g.Player.Y+8, ScreenWidth, ScreenHeight)
	g.cam.Constraint(
		float64(g.tilemapJSON.Layers[0].Width)*16.0,
		float64(g.tilemapJSON.Layers[0].Width)*16.0,
		ScreenWidth,
		ScreenHeight,
	)

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
			drawOptions.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				drawOptions,
			)

			// Reset the transformation matrix for the next draw, otherwise the next tranlation wil be wrong.
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
		drawOptions.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(enemyImg, drawOptions)

		drawOptions.GeoM.Reset()
	}

	for _, potion := range g.Potions {
		potionImge := potion.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image)

		drawOptions.GeoM.Translate(potion.X, potion.Y)
		drawOptions.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(potionImge, drawOptions)

		drawOptions.GeoM.Reset()
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImage, _, err := ebitenutil.NewImageFromFile("assets/images/ninja.png")
	if err != nil {
		log.Fatal("Failed to load player image:", err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal("Failed to load skeleton image:", err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/Heart.png")
	if err != nil {
		log.Fatal("Failed to load potion image:", err)
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
		Player: &entities.Sprite{
			Img:    playerImage,
			Health: 100,
			X:      50,
			Y:      50,
		},
		Enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   100.4,
					Y:   100.4,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   200.4,
					Y:   200.4,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   300.4,
					Y:   300.4,
				},
				FollowsPlayer: false,
			},
		},
		Potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   200,
					Y:   200,
				},
				AmtHeal: 20,
			},
		},
		tilemapJSON: tilemapJson,
		tilemapImg:  tilemapImg,
		cam:         NewCamera(0, 0),
	}

	// Start the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
