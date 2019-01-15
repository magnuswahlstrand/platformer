package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/peterhellberg/gfx"

	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ebitenconsole"
	"github.com/kyeett/tiled"
)

func (g *Game) update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return errors.New("Player exited game")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {

		// Load initial size from first world map

		g.currentScene = "game"
		worldFile := "world6"
		worldMap, err := tiled.MapFromFile(fmt.Sprintf("%s/%s.tmx", g.baseDir, worldFile))
		if err != nil {
			log.Fatal(err)
		}
		g.initializeWorld(worldMap)
	}
	return g.scenes[g.currentScene](g, screen)
}

func VictoryScreen(g *Game, screen *ebiten.Image) error {
	screen.Fill(colornames.Black)
	drawCenterText(screen, "Victory!", fontFace11, colornames.White)
	drawCenterText(screen, "Press R to restart game", fontFace5, colornames.White, 40)
	return nil
}

func LostScreen(g *Game, screen *ebiten.Image) error {
	screen.Fill(colornames.Black)
	drawCenterText(screen, "You lost!", fontFace11, colornames.White)
	drawCenterText(screen, "Press R to restart game", fontFace5, colornames.White, 40)
	return nil
}

var camera *ebiten.Image

var i1, i2 int
var tot1, tot2, tot3 time.Duration

func GameLoop(g *Game, screen *ebiten.Image) error {

	camera, _ = ebiten.NewImageFromImage(gfx.NewImage(g.Width, g.Height, colornames.Red), ebiten.FilterDefault)
	ebitenconsole.CheckInput()

	// Inspired by https://forums.tigsource.com/index.php?topic=46289.msg1386874#msg1386874
	// UpdatePreMovement
	// Apply friction, gravity and keypresses
	g.updatePreMovement()

	g.updateMovement(camera)

	// For each Entity
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	g.updatePostMovement()

	// Draw background

	// camera.DrawImage(foregroundImg, &ebiten.DrawImageOptions{})
	g.drawBackground(camera)

	// Draw entities
	g.drawEntities(camera)

	g.drawPlayerVision(camera)

	// Draw foreground

	// Check for collision with triggers
	g.checkAndDrawTriggers(camera)

	g.drawHitboxes(camera)

	cr := g.getCameraPosition()
	screen.DrawImage(camera.SubImage(cr).(*ebiten.Image), &ebiten.DrawImageOptions{})
	g.drawScoreboard(screen)
	g.drawDebugInfo(screen)

	return nil
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}
func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}
