package main

import (
	"errors"
	"fmt"
	"log"

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

func GameLoop(g *Game, screen *ebiten.Image) error {
	ebitenconsole.CheckInput()

	// Inspired by https://forums.tigsource.com/index.php?topic=46289.msg1386874#msg1386874
	// UpdatePreMovement
	// Apply friction, gravity and keypresses
	g.updatePreMovement()

	g.updateMovement(screen)

	// For each Entity
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	g.updatePostMovement()

	// Draw background
	screen.DrawImage(backgroundImg, &ebiten.DrawImageOptions{})

	// screen.DrawImage(foregroundImg, &ebiten.DrawImageOptions{})
	// Draw entities
	g.drawEntities(screen)

	g.drawPlayerVision(screen)

	// Draw foreground

	// Check for collision with triggers
	g.checkAndDrawTriggers(screen)

	g.drawHitboxes(screen)

	g.drawScoreboard(screen)
	g.drawDebugInfo(screen)
	return nil
}
