package main

import (
	"fmt"

	"github.com/peterhellberg/gfx"

	"github.com/hajimehoshi/ebiten"

	"golang.org/x/image/colornames"

	"github.com/kyeett/gomponents/components"
)

func direction(v gfx.Vec) byte {
	var dir byte
	if v.X > 0 {
		dir |= components.DirRight
	}

	if v.X < 0 {
		dir |= components.DirLeft
	}

	if v.Y > 0 {
		dir |= components.DirUp
	}

	if v.Y < 0 {
		dir |= components.DirDown
	}
	return dir
}

func (g *Game) checkAndDrawTriggers(screen *ebiten.Image) {
	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	v := g.entities.GetUnsafe(playerID, components.VelocityType).(*components.Velocity)
	hb := g.entities.GetUnsafe(playerID, components.HitboxType).(*components.Hitbox)

	playerDirection := direction(v.Vec)
	playerShape := rectToShape(hb.Moved(pos.Vec))

	for _, e := range g.filteredEntities(components.TriggerType) {
		t := g.entities.GetUnsafe(e, components.TriggerType).(*components.Trigger)

		tRect := rectToShape(t.Rect)
		var collided bool
		if playerShape.WouldBeColliding(tRect, 0, 0) && t.Direction&playerDirection > 0 {
			fmt.Println(components.DirString(playerDirection), components.DirString(t.Direction), "res", components.DirString(t.Direction&playerDirection), t.Direction&playerDirection)
			fmt.Println("triggered!")
			collided = true
		}

		// Draw triggers
		if collided {
			drawPixelRect(screen, t.Rect, colornames.Red)
		} else {
			drawPixelRect(screen, t.Rect, colornames.Thistle)
		}

	}
}
