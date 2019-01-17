package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"text/tabwriter"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"

	"github.com/kyeett/ebitenconsole"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
)

var fontFace5 font.Face
var fontFace7 font.Face
var fontFace9 font.Face
var fontFace11 font.Face

func init() {
	fnt, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal("loading font:", err)
	}
	const dpi = 144
	fontFace5 = truetype.NewFace(fnt, &truetype.Options{
		Size:    5,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	fontFace7 = truetype.NewFace(fnt, &truetype.Options{
		Size:    7,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	fontFace9 = truetype.NewFace(fnt, &truetype.Options{
		Size:    9,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	fontFace11 = truetype.NewFace(fnt, &truetype.Options{
		Size:    11,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
}

func drawCenterText(screen *ebiten.Image, txt string, face font.Face, c color.Color, offsetY ...int) {
	y := 0
	for _, o := range offsetY {
		y += o
	}
	size := face.Metrics().Height.Ceil() / 2
	width := int(1.135 * float64(len(txt)*size))
	w, h := screen.Size()
	text.Draw(screen, txt, face, (w-width)/2, (h+size)/2+y, c)
}

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

	ebitenutil.DebugPrintAt(screen, buf.String(), 50, 60)

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

		if g.entities.HasComponents(e, components.VelocityType) && e != playerID {
			v := g.entities.GetUnsafe(e, components.VelocityType).(*components.Velocity)
			w, _ := s.Size()
			if v.X > 0 {
				op.GeoM.Translate(float64(-w), 0)
				op.GeoM.Scale(-1, 1)
			}
		}

		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(img, op)

		if debug {
			ebitenutil.DebugPrintAt(screen, e, int(pos.X), int(pos.Y))
		}

		// For debug, add dot to mark position
		ebitenutil.DrawRect(traceImg, pos.X, pos.Y, 1, 1, colornames.Red)
	}
}

func (g *Game) drawBackground(camera *ebiten.Image) {
	cr := g.getCameraPosition()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cr.Min.X), 0)
	camera.DrawImage(backgroundImg.SubImage(cr).(*ebiten.Image), op)
}

func drawTrail(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.98)
	tmpImg.Clear()
	tmpImg.DrawImage(traceImg, op)
	op.ColorM.Scale(1, 1, 1, 1)
	traceImg.Clear()
	traceImg.DrawImage(tmpImg, op)

	// Draw trace
	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(traceImg, op)
}

func (g *Game) drawHitboxes(screen *ebiten.Image) {
	if hitbox {
		// Draw hitboxes
		for _, e := range g.filteredEntities(components.HitboxType, components.PosType) {
			pos := g.entities.GetUnsafe(e, components.PosType).(*components.Pos)
			hb := g.entities.GetUnsafe(e, components.HitboxType).(*components.Hitbox)

			switch {
			case hb.Properties["allow_from_down"]:
				drawPixelRect(screen, hb.Moved(pos.Vec), colornames.Turquoise)
			case hb.Properties["monsters_only"]:
				drawPixelRect(screen, hb.Moved(pos.Vec), colornames.Greenyellow)

			default:
				drawPixelRect(screen, hb.Moved(pos.Vec), colornames.Red)
			}
		}

	}
}

func (g *Game) getCameraPosition() image.Rectangle {
	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	cx := int(pos.X - cameraWidth/2)
	cx = min(g.Width-cameraWidth, cx)
	cx = max(0, cx)
	return image.Rect(cx, 0, cx+cameraWidth, cameraHeight)
}

func (g *Game) drawPlayerVision(screen *ebiten.Image) {
	// cr := g.getCameraPosition()

	// opt.Address = ebiten.AddressRepeat

	opt := &ebiten.DrawImageOptions{}
	tmpVisionImg.DrawImage(foregroundImg, opt) //.SubImage(cr).(*ebiten.Image)

	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	opt.GeoM.Translate(pos.X-float64(visionImg.Bounds().Dx())/2+15, pos.Y-float64(visionImg.Bounds().Dy())/2+20)
	// opt.GeoM.Translate(float64(cr.Min.X), 0)
	opt.CompositeMode = ebiten.CompositeModeDestinationOut

	// tmp, _ := ebiten.NewImage(200, 200, ebiten.FilterDefault)
	// tmp.Fill(colornames.Red)
	// backgroundImg.DrawImage(visionImg, &ebiten.DrawImageOptions{})

	tmpVisionImg.DrawImage(visionImg, opt)

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(tmpVisionImg, op)
}
