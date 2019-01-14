package main

import (
	"fmt"

	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

func hitboxToRect(hb gfx.Rect) *resolv.Rectangle {
	return resolv.NewRectangle(int32(hb.Min.X), int32(hb.Min.Y), int32(hb.W()), int32(hb.H()))
}

const factor = 100

func (g *Game) updateMovement(screen *ebiten.Image) {

	var space resolv.Space
	// Add possible collision entities
	for _, e := range g.filteredEntities(components.HitboxType) {
		if e == playerID {
			continue
		}
		pos := g.entities.GetUnsafe(e, components.PosType).(*components.Pos)
		hb := g.entities.GetUnsafe(e, components.HitboxType).(*components.Hitbox)
		hbMoved := hb.Moved(pos.Vec)

		// Debug things

		scaler := hbMoved.Size().Scaled(factor)
		resizedBox := hbMoved.Resized(gfx.V(0, 0), scaler)

		s := hitboxToRect(resizedBox)
		// s.SetTags(e)
		// s.SetData(hb)
		s.SetTags(e)
		if hb.Properties["allow_from_down"] {
			s.SetTags("allow_from_down")
		}

		space.AddShape(s)
	}

	for _, e := range []string{playerID} {
		pos := g.entities.GetUnsafe(e, components.PosType).(*components.Pos)
		v := g.entities.GetUnsafe(e, components.VelocityType).(*components.Velocity)
		hb := g.entities.GetUnsafe(e, components.HitboxType).(*components.Hitbox)
		hbMoved := hb.Moved(pos.Vec)
		scaler := hb.Size().Scaled(factor)
		r := hitboxToRect(hbMoved.Resized(gfx.V(0, 0), scaler))

		// Check collision vertically

		filterFunc := func(s resolv.Shape) bool { return true }
		if v.Y < 0 {
			filterFunc = func(s resolv.Shape) bool {
				return !s.HasTags("allow_from_down")
			}
		}
		verticalSpace := space.Filter(filterFunc)

		if res := verticalSpace.Resolve(r, 0, int32(factor*v.Y)); res.Colliding() && !res.Teleporting {
			t := res.ShapeB.GetTags()[0]
			// Calculate distance to object
			// Todo, fix
			_, bY := res.ShapeB.GetXY()

			entityUnderneath := v.Y > 0

			if entityUnderneath {
				fac := hb.Max.Y
				pos.Y = float64(bY/factor) - fac
			} else if v.Y < 0 { // Underneath
				fmt.Println("Underneath")
			}

			if g.entities.HasComponents(t, components.BouncyType) {
				v.Y = -4
			} else {
				v.Y = 0
			}

			// Killed!
			if g.entities.HasComponents(t, components.KillableType) {
				g.handleKilled(t)
			}

		} else {
			pos.Y += v.Y
		}

		if res := space.Resolve(r, int32(factor*v.X), 0); res.Colliding() && !res.Teleporting {
			t := res.ShapeB.GetTags()[0]

			if g.entities.HasComponents(t, components.HazardType) {
				g.Reset()
				return
			}

			// Bounce if not jumping or falling
			if v.Y == 0 {
				v.X = -0.5 * v.X
			}

		} else {
			pos.X += v.X
		}
	}
}

func (g *Game) handleKilled(t string) {
	pos := g.entities.GetUnsafe(t, components.PosType).(*components.Pos)
	pos.Y += 6
	g.entities.Remove(t, components.HitboxType)
	g.entities.Add(t, components.Rotated{0.0})
	g.entities.Add(t, components.Scenario{
		F: func() bool {
			pas := g.entities.GetUnsafe(t, components.PosType).(*components.Pos)
			pas.Y++

			rot := g.entities.GetUnsafe(t, components.RotatedType).(*components.Rotated)
			rot.Rotate(0.1)

			return pas.Y > float64(g.Height)
		},
	})
}
