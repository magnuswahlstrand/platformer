package main

import (
	"image"
	"image/color"

	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func drawPixelRect(screen *ebiten.Image, r gfx.Rect, c color.Color) {
	ebitenutil.DrawLine(screen, r.Min.X+1, r.Min.Y+1, r.Min.X+1, r.Max.Y-1, c)
	ebitenutil.DrawLine(screen, r.Min.X+1, r.Max.Y-1, r.Max.X-1, r.Max.Y-1, c)
	ebitenutil.DrawLine(screen, r.Max.X-1, r.Max.Y-1, r.Max.X-1, r.Min.Y+1, c)
	ebitenutil.DrawLine(screen, r.Max.X-1, r.Min.Y+1, r.Min.X+1, r.Min.Y+1, c)
}

func drawPixelFilledRect(screen *ebiten.Image, r gfx.Rect, c color.Color) {
	ebitenutil.DrawRect(screen, r.Min.X, r.Min.Y, r.W(), r.H(), c)
}

func drawRect(screen *ebiten.Image, x, y, w, h float64, c color.Color) {
	ebitenutil.DrawLine(screen, x, y, x, y+h, c)
	ebitenutil.DrawLine(screen, x, y+h, x+w, y+h, c)
	ebitenutil.DrawLine(screen, x+w, y+h, x+w, y, c)
	ebitenutil.DrawLine(screen, x+w, y, x, y, c)
}

func (g *Game) drawEntities(screen *ebiten.Image) {

	// Draw entitief
	for _, e := range g.filteredEntities(components.DrawableType, components.PosType) {
		pos := g.entities.GetUnsafe(e, components.PosType).(*components.Pos)
		s := g.entities.GetUnsafe(e, components.DrawableType).(*components.Drawable)
		img := s.Image
		// If animated
		if g.entities.HasComponents(e, components.AnimatedType) {
			a := g.entities.GetUnsafe(e, components.AnimatedType).(*components.Animated)
			w, h := a.Ase.FrameWidth, a.Ase.FrameHeight
			x, y := a.Ase.GetFrameXY()
			img = img.SubImage(image.Rect(int(x), int(y), int(x+w), int(y+h))).(*ebiten.Image)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(img, op)

		// For debug, add dot to mark position
		ebitenutil.DrawRect(traceImg, pos.X, pos.Y, 1, 1, colornames.Red)
	}

	if hitbox {
		// Draw hitboxes
		for _, e := range g.filteredEntities(components.HitboxType, components.PosType) {
			pos := g.entities.GetUnsafe(e, components.PosType).(*components.Pos)
			hb := g.entities.GetUnsafe(e, components.HitboxType).(*components.Hitbox)

			if hb.Properties["allow_from_down"] {
				drawPixelRect(screen, hb.Moved(pos.Vec), colornames.Turquoise)
			} else {
				drawPixelRect(screen, hb.Moved(pos.Vec), colornames.Red)
			}
		}
	}
}
