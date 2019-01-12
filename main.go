package main

import (
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ebitenconsole"
)

const (
	screenWidth  = 240
	screenHeight = 240
)

var hitbox = true

func main() {

	g := NewGame()

	ebitenconsole.FloatVar(&g.Gravity, "g", "world gravity")
	ebitenconsole.BoolVar(&hitbox, "h", "show hitboxes")
	currentTime = time.Now()

	if err := ebiten.Run(g.update, screenWidth, screenHeight, 2, "Aseprite demo"); err != nil {
		log.Fatal(err)
	}
}
