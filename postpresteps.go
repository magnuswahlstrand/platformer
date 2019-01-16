package main

import (
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/gomponents/components"
)

func (g *Game) updatePreMovement() {
	for _, e := range []string{playerID} {
		v := g.entities.GetUnsafe(e, components.VelocityType).(*components.Velocity)

		// Gravity
		v.Y += g.Gravity

		// Frictio
		v.X = 0.90 * v.X
		d := g.entities.GetUnsafe(e, components.DirectionType).(*components.Direction)
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyUp):
			v.Y = -jumpSpeed
			musicPlayer.PlayAudio(jumpSound)

		case ebiten.IsKeyPressed(ebiten.KeyRight):
			v.X += horizontalAcceleration
			d.D = 1

		case ebiten.IsKeyPressed(ebiten.KeyLeft):
			v.X -= horizontalAcceleration
			d.D = -1
		}
	}
}

func (g *Game) updatePostMovement() {
	diffTime = time.Since(currentTime)
	currentTime = time.Now()

	for _, e := range []string{playerID} {
		a := g.entities.GetUnsafe(e, components.AnimatedType).(*components.Animated)

		// Update animation time
		a.Ase.Update(float32(diffTime.Nanoseconds()) / 1000000000)

		// Update animation based on velocity
		v := g.entities.GetUnsafe(e, components.VelocityType).(*components.Velocity)
		d := g.entities.GetUnsafe(e, components.DirectionType).(*components.Direction)

		var direction string
		switch float64(d.D) {
		case -1.0:
			direction = "left"
		case 1.0:
			direction = "right"
		}

		switch {
		case v.Y > 0.03:
			a.Ase.Play("fall " + direction)
		case v.Y < -0.03:
			a.Ase.Play("jump " + direction)
		case v.X > 0.03:
			a.Ase.Play("walk " + direction)
		case v.X < -0.03:
			a.Ase.Play("walk " + direction)
		default:
			a.Ase.Play("stand " + direction)

		}
	}

	// Perform entity scenarios
	for _, e := range g.filteredEntities(components.ScenarioType) {
		scenario := g.entities.GetUnsafe(e, components.ScenarioType).(*components.Scenario)
		finished := scenario.F()

		if finished == true {
			g.entities.RemoveAll(e)
			var entities []string

			// Remove entity from list
			for _, s := range g.entityList {
				if s == e {
					continue
				}
				entities = append(entities, s)
			}
			g.entityList = entities
		}
	}
}
