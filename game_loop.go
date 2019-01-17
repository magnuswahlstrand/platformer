package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

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

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		musicPlayer.audioPlayer.SetVolume(1 - musicPlayer.audioPlayer.Volume())
	}

	if musicPlayer.audioPlayer != nil && musicPlayer.audioPlayer.IsPlaying() {
		musicPlayer.current = musicPlayer.audioPlayer.Current()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {

		// Load initial size from first world map

		g.currentScene = "game"
		worldFile := "world6.tmx"
		worldMap := g.loadWorldMap(worldFile)
		g.initializeWorld(worldMap)
	}
	return g.scenes[g.currentScene](g, screen)
}

func (g *Game) loadWorldMap(filename string) *tiled.Map {
	// Load initial size from first world map

	// readFromMap := func(filename string) ([]byte, error) {
	// 	fmt.Println("load tileset", filename)
	// 	return resources.LookupFatal(g.baseDir + "/" + filename), nil
	// }
	// worldMap, err := tiled.MapFromBytes(resources.LookupFatal(g.baseDir+"/"+filepath.Base(filename)), readFromMap)

	fmt.Println("load tilemap", filename)
	worldMap, err := tiled.MapFromFile(g.baseDir + "/" + filepath.Base(filename))
	if err != nil {
		log.Fatal(err)
	}
	return worldMap
}

func victoryScreen(g *Game, screen *ebiten.Image) error {
	screen.Fill(colornames.Black)
	drawCenterText(screen, "Victory!", fontFace11, colornames.White)
	drawCenterText(screen, "Press R to restart game", fontFace5, colornames.White, 40)
	return nil
}

func lostScreen(g *Game, screen *ebiten.Image) error {
	screen.Fill(colornames.Black)
	drawCenterText(screen, "You lost!", fontFace11, colornames.White)
	drawCenterText(screen, "Press R to restart game", fontFace5, colornames.White, 40)
	return nil
}

var camera *ebiten.Image

var i1, i2 int
var tot1, tot2, tot3 time.Duration

func gameLoop(g *Game, screen *ebiten.Image) error {

	t1 := time.Now()
	camera.Clear()
	ebitenconsole.CheckInput()

	// Inspired by https://forums.tigsource.com/index.php?topic=46289.msg1386874#msg1386874
	// UpdatePreMovement
	// Apply friction, gravity and keypresses
	g.updatePreMovement()

	t2 := time.Now()
	g.updateMovement(camera)

	// For each Entity
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	g.updatePostMovement()
	t3 := time.Now()

	// Draw background

	// camera.DrawImage(foregroundImg, &ebiten.DrawImageOptions{})
	g.drawBackground(camera)
	t4 := time.Now()

	// Draw entities
	g.drawEntities(camera)
	t5 := time.Now()

	g.drawPlayerVision(camera)

	if hitbox {
		drawTrail(camera)
	}
	// Draw foreground

	// Check for collision with triggers
	g.checkAndDrawTriggers(camera)

	t6 := time.Now()
	g.drawHitboxes(camera)

	cr := g.getCameraPosition()

	t7 := time.Now()
	screen.DrawImage(camera.SubImage(cr).(*ebiten.Image), &ebiten.DrawImageOptions{})
	g.drawScoreboard(screen)
	g.drawDebugInfo(screen)
	t8 := time.Now()

	_ = fmt.Sprintf("%10s %10s %10s %10s %10s %10s %10s, tot: %10s\n", t2.Sub(t1), t3.Sub(t2), t4.Sub(t3), t5.Sub(t4), t6.Sub(t5), t7.Sub(t6), t8.Sub(t7), t8.Sub(t1))
	// fmt.Printf("%10s %10s %10s %10s %10s %10s %10s, tot: %10s\n", t2.Sub(t1), t3.Sub(t2), t4.Sub(t3), t5.Sub(t4), t6.Sub(t5), t7.Sub(t6), t8.Sub(t7), t8.Sub(t1))
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
