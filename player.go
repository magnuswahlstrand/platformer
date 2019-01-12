package main

import (
	"bytes"
	"image"
	"log"

	resources "github.com/kyeett/platformer/assets"

	"github.com/hajimehoshi/ebiten"

	ase "github.com/kyeett/GoAseprite"
)

type Player struct {
	Ase       ase.File
	Sprite    *ebiten.Image
	direction string
}

var playerFile ase.File
var pImage *ebiten.Image
var tileImage *ebiten.Image

func init() {

	// goaseprite.Load() returns an AsepriteFile, assuming it finds the JSON file
	// playerFile = ase.Load("assets/graphics/character.json")

	// f, err := os.Open(playerFile.ImagePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// img, _, err := image.Decode(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pImage, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// playerFile = ase.Load("assets/graphics/character.json")
	playerFile = ase.LoadBytes(resources.Character_json)

	img, _, err := image.Decode(bytes.NewReader(resources.Character_png))
	if err != nil {
		log.Fatal(err)
	}

	pImage, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	//	img, _, err := image.Decode(bytes.NewReader(resources.Tiles_png))

	img, _, err = image.Decode(bytes.NewReader(resources.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}

	tileImage, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	// player.Sprite = pImage

	// padd := 2.0

}

// func NewPlayer() *Player {

// 	player := Player{}
// 	player.direction = "right"

// 	// goaseprite.Load() returns an AsepriteFile, assuming it finds the JSON file
// 	player.Ase = ase.Load("assets/graphics/character.json")

// 	// AsepriteFile.ImagePath will be the absolute path to the image file.
// 	// player.Texture = raylib.LoadTexture(player.Ase.ImagePath)

// 	f, err := os.Open(player.Ase.ImagePath)
// 	if err != nil {

// 	}

// 	img, _, err := image.Decode(f)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	pImage, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	player.Sprite = pImage

// 	padd := 2.0

// 	// // Set up the texture rectangle for drawing the sprite
// 	// player.TextureRect = raylib.Rectangle{0, 0, player.Ase.FrameWidth, player.Ase.FrameHeight}

// 	// Queues up the "Play" animation
// 	player.Ase.Play("stand right")

// 	return &player

// }

// func (p *Player) WhileJumping() {
// 	if p.V.Y > 0 {
// 		p.Ase.Play("fall " + p.direction)
// 	} else {
// 		p.Ase.Play("jump " + p.direction)
// 	}

// 	// Apply gravity

// 	p.Pos.Y += p.V.Y
// 	// Falling complete, go back to standing
// 	if p.Pos.X > 170 {
// 		p.Pos.Y = 170
// 		p.V.Y = 0
// 		p.Ase.Play("stand " + p.direction)
// 	}

// }

// func (p *Player) Draw(screen *ebiten.Image) {
// 	w, h := p.Ase.FrameWidth, p.Ase.FrameHeight
// 	x, y := p.Ase.GetFrameXY()
// 	op := &ebiten.DrawImageOptions{}
// 	op.GeoM.Translate(p.Pos.X, p.Pos.Y)
// 	screen.DrawImage(p.Sprite.SubImage(image.Rect(int(x), int(y), int(x+w), int(y+h))).(*ebiten.Image), op)
// }
// func (p *Player) Update() {
// 	p.Ase.Update(float32(diffTime.Nanoseconds()) / 1000000000)
// }

// func (p *Player) Move() {

// 	// Apply direction, wether jumping or not
// 	switch {
// 	case ebiten.IsKeyPressed(ebiten.KeyRight):
// 		p.Pos.X += 2
// 		p.direction = "right"
// 	case ebiten.IsKeyPressed(ebiten.KeyLeft):
// 		p.Pos.X -= 2
// 		p.direction = "left"
// 	}

// 	if p.Pos.X > screenWidth-float64(p.Ase.FrameWidth) {
// 		p.Pos.X = screenWidth - float64(p.Ase.FrameWidth)
// 	}
// 	if p.Pos.X < 0 {
// 		p.Pos.X = 0
// 	}

// 	if strings.Contains(p.Ase.CurrentAnimation.Name, "jump") || strings.Contains(p.Ase.CurrentAnimation.Name, "fall") {
// 		p.WhileJumping()
// 		return
// 	}

// 	switch {
// 	case ebiten.IsKeyPressed(ebiten.KeyUp):
// 		p.Ase.Play("jump " + p.direction)
// 		p.V.Y = -5
// 	case ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyLeft):
// 		p.Ase.Play("walk " + p.direction)
// 	default:
// 		p.Ase.Play("stand " + p.direction)
// 	}
// }
