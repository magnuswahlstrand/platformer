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
var miscImage *ebiten.Image

func init() {

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
	img, _, err = image.Decode(bytes.NewReader(resources.Misc_png))
	if err != nil {
		log.Fatal(err)
	}

	miscImage, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

}
