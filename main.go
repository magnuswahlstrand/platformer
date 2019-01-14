package main

import (
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ebitenconsole"
	"github.com/sirupsen/logrus"
)

var hitbox = true
var bnw = false

func main() {

	g := NewGame()
	screenWidth, screenHeight := g.Width, g.Height
	ebitenconsole.FloatVar(&g.Gravity, "g", "world gravity")
	ebitenconsole.BoolVar(&hitbox, "h", "show hitboxes")
	ebitenconsole.BoolVar(&bnw, "bnw", "change game palette")
	currentTime = time.Now()

	if err := ebiten.Run(g.update, screenWidth, screenHeight, 3, "Aseprite demo"); err != nil {
		logrus.Error(err)
	}
}
