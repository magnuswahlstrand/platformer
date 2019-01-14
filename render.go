package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"text/tabwriter"
	"time"

	"github.com/kyeett/ebitenconsole"
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

var currentTime time.Time
var diffTime time.Duration

func (g *Game) drawScoreboard(screen *ebiten.Image) {
	fullHeart := image.Rect(0, 8, 8, 16)
	emptyHeart := image.Rect(8, 8, 16, 16)
	screen.DrawImage(scoreboardImg, &ebiten.DrawImageOptions{})

	counter := g.entities.GetUnsafe(playerID, components.CounterType).(*components.Counter)

	for l := 0.0; l < 3; l++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(8*l+4, 4)
		if int(l) >= (*counter)["lives"] {
			screen.DrawImage(miscImage.SubImage(emptyHeart).(*ebiten.Image), op)
		} else {
			screen.DrawImage(miscImage.SubImage(fullHeart).(*ebiten.Image), op)
		}
	}
}

func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	v := g.entities.GetUnsafe(playerID, components.VelocityType).(*components.Velocity)
	buf := bytes.NewBufferString("")
	wr := tabwriter.NewWriter(buf, 0, 1, 3, ' ', 0)
	fmt.Fprintf(wr, "x, y:\t(%4.0f,%4.0f)\t\n\t(%4.2f,%4.2f)\t", pos.X, pos.Y, v.X, v.Y)
	wr.Flush()

	// ebitenutil.DebugPrintAt(screen, buf.String(), 50, 60)

	ebitenutil.DebugPrintAt(screen, ebitenconsole.String(), 0, 40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f", ebiten.CurrentTPS()), 190, 0)
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
		if g.entities.HasComponents(e, components.RotatedType) {
			rot := g.entities.GetUnsafe(e, components.RotatedType).(*components.Rotated)
			w, h := s.Size()
			op.GeoM.Translate(-float64(w/2), float64(-h/2))
			op.GeoM.Rotate(rot.Angle)
			op.GeoM.Translate(float64(w/2), float64(h/2))
		}

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
