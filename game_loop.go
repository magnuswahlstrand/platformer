package main

import (
	"errors"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ebitenconsole"
)

func (g *Game) update(screen *ebiten.Image) error {
	ebitenconsole.CheckInput()

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return errors.New("Player exited game")
	}

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
