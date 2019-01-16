package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ebitenconsole"
	"github.com/sirupsen/logrus"
)

var (
	musicPlayer *Player
)

func main() {

	var world string
	flag.StringVar(&world, "world", "world6", "world to play in")
	flag.Parse()
	tmxPath := fmt.Sprintf("../tiled/%s.tmx", world)

	g := NewGame(tmxPath)
	// save = func() error {
	// 	g.entities
	// }
	// screenWidth, screenHeight := g.Width, g.Height
	ebitenconsole.FloatVar(&g.Gravity, "g", "world gravity")
	ebitenconsole.BoolVar(&hitbox, "h", "show hitboxes")
	ebitenconsole.BoolVar(&bnw, "bnw", "change game palette")
	// ebitenconsole.FuncVar(save, "save", "saves state")
	// ebitenconsole.FuncVar(load, "load", "load state")
	currentTime = time.Now()
	var err error
	musicPlayer, err = NewPlayer()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.Run(g.update, cameraWidth, cameraHeight, 3, "Aseprite demo"); err != nil {
		logrus.Error(err)
	}

}
