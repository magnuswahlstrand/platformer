package main

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/SolarLune/resolv/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
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
			v.Y = -5

		case ebiten.IsKeyPressed(ebiten.KeyRight):
			fmt.Print("Move right")
			v.X += 0.3
			d.D = 1

		case ebiten.IsKeyPressed(ebiten.KeyLeft):
			fmt.Print("-Move Left")
			v.X -= 0.3
			d.D = -1
		}
	}
}

func hitboxToRect(hb gfx.Rect) *resolv.Rectangle {
	return resolv.NewRectangle(int32(hb.Min.X), int32(hb.Min.Y), int32(hb.W()), int32(hb.H()))
}

const factor = 100000

func (g *Game) updateMovement(screen *ebiten.Image) {

	var space resolv.Space
	for _, e := range g.filteredEntities(components.HitboxType) {
		if e == playerID {
			continue
		}
		pos := g.entities.GetUnsafe(e, components.PosType).(*components.Pos)
		hb := g.entities.GetUnsafe(e, components.HitboxType).(*components.Hitbox)
		hbMoved := hb.Moved(pos.Vec)

		// Debug things
		// posBox := gfx.R(pos.X, pos.Y, pos.X+30, pos.Y+30)

		// fmt.Println("Before", hb, hb.Size(), gfx.V(0, 0))
		scaler := hbMoved.Size().Scaled(factor)
		// scaledBox := gfx.R(0, 0, scaler.X, scaler.Y)
		resizedBox := hbMoved.Resized(gfx.V(0, 0), scaler)
		drawPixelFilledRect(screen, resizedBox, colornames.Yellowgreen)

		// fmt.Println("After", resizedBox, scaler)

		s := hitboxToRect(resizedBox)
		s.SetTags(e)
		s.SetData(hb)
		space.AddShape(s)

		// drawPixelFilledRect(screen, posBox, colornames.Pink)
		// drawPixelFilledRect(screen, hbBox, colornames.Yellow)
		// drawPixelFilledRect(screen, scaledBox, colornames.Blue)
		// drawPixelFilledRect(screen, resizedBox, colornames.Wheat)
		// drawPixelRect(screen, resizedBox, colornames.Wheat)
	}
	// log.Fatal("")

	for i := 0; i < len(space); i++ {
		x, y := space.Get(i).GetXY()
		drawRect(screen, float64(x), float64(y), 15, 15, colornames.Turquoise)
	}

	for _, e := range []string{playerID} {
		pos := g.entities.GetUnsafe(e, components.PosType).(*components.Pos)
		v := g.entities.GetUnsafe(e, components.VelocityType).(*components.Velocity)
		hb := g.entities.GetUnsafe(e, components.HitboxType).(*components.Hitbox)
		hbMoved := hb.Moved(pos.Vec)
		scaler := hb.Size().Scaled(factor)
		r := hitboxToRect(hbMoved.Resized(gfx.V(0, 0), scaler))

		// Check collision vertically
		if res := space.Resolve(r, 0, int32(factor*v.Y)); res.Colliding() && !res.Teleporting {

			// Calculate distance to object
			// Todo, fix
			_, bY := res.ShapeB.GetXY()
			if v.Y > 0 { // Above
				// fmt.Println("Above")
				fac := hb.Max.Y
				// fmt.Println(hb.H(), fac)
				pos.Y = float64(bY/factor) - fac
			} else if v.Y < 0 { // Underneath
				// fmt.Println("Underneath")
			}

			v.Y = 0
			// log.Fatal("yay")
		} else {
			pos.Y += v.Y
		}

		if res := space.Resolve(r, int32(factor*v.X), 0); res.Colliding() && !res.Teleporting {
			// Bound if not jumping or falling
			if v.Y == 0 {

				v.X = -0.5 * v.X
			}

		} else {
			pos.X += v.X
		}
	}
}

func (g *Game) updatePostMovement() {
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
}
