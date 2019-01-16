package main

import (
	"fmt"
	"strings"

	"github.com/kyeett/gomponents/direction"

	"github.com/hajimehoshi/ebiten"

	"golang.org/x/image/colornames"

	"github.com/kyeett/gomponents/components"
)

func (g *Game) checkAndDrawTriggers(screen *ebiten.Image) {
	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	v := g.entities.GetUnsafe(playerID, components.VelocityType).(*components.Velocity)
	hb := g.entities.GetUnsafe(playerID, components.HitboxType).(*components.Hitbox)

	playerDirection := direction.FromVec(v.Vec)

	playerShape := rectToShape(hb.Moved(pos.Vec))

	for _, e := range g.filteredEntities(components.TriggerType) {
		t := g.entities.GetUnsafe(e, components.TriggerType).(*components.Trigger)

		tRect := rectToShape(t.Rect)
		var collided bool
		if playerShape.WouldBeColliding(tRect, 0, 0) && t.Direction&playerDirection > 0 {
			fmt.Println(playerDirection, t.Direction, "res", t.Direction&playerDirection, t.Direction&playerDirection)
			fmt.Println("triggered!")

			switch {
			case t.Scenario == "victory":
				g.currentScene = "victory"
			case strings.Contains(t.Scenario, "world:"):
				// Load initial size from first world map

				worldFile := strings.Replace(t.Scenario, "world:", "", -1) + ".tmx"
				worldMap := g.loadWorldMap(worldFile)
				g.initializeWorld(worldMap)
			}

			collided = true
		}

		// Draw triggers
		if hitbox {
			if collided {
				drawPixelRect(screen, t.Rect, colornames.Red)
			} else {
				drawPixelRect(screen, t.Rect, colornames.Thistle)
			}
		}
	}
}
