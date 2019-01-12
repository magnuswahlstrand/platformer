package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"text/tabwriter"
	"time"

	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ebitenconsole"
)

var tmpImg *ebiten.Image
var traceImg *ebiten.Image

func init() {
	traceImg, _ = ebiten.NewImage(240, 240, ebiten.FilterDefault)
	tmpImg, _ = ebiten.NewImage(240, 240, ebiten.FilterDefault)
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
	// UpdatePostMovement

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// drawRect(sceen, h.X, y int32, w int32, h int32, c invalid type)
	// func drawRect(screen *ebiten.Image, x, y, w, h int32, c color.Color) {
	// ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), c)

	drawTrail(screen)

	diffTime = time.Since(currentTime)
	currentTime = time.Now()

	g.updatePostMovement()

	// Draw entities
	g.drawEntities(screen)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Current animation:   %s\nFrame (index/total): %d/%d",
	// 	g.player.Ase.CurrentAnimation.Name,
	// 	g.player.Ase.CurrentFrame-g.player.Ase.CurrentAnimation.Start,
	// 	g.player.Ase.CurrentAnimation.End-g.player.Ase.CurrentAnimation.Start+1))

	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	v := g.entities.GetUnsafe(playerID, components.VelocityType).(*components.Velocity)
	buf := bytes.NewBufferString("")
	wr := tabwriter.NewWriter(buf, 0, 1, 3, ' ', 0)
	fmt.Fprintf(wr, "x, y:\t(%4.0f,%4.0f)\t\n\t(%4.2f,%4.2f)\t", pos.X, pos.Y, v.X, v.Y)
	wr.Flush()

	ebitenutil.DebugPrintAt(screen, buf.String(), 50, 60)

	ebitenutil.DebugPrintAt(screen, ebitenconsole.String(), 0, 220)
	ebitenconsole.String()
	return nil
}

var currentTime time.Time
var diffTime time.Duration

type Game struct {
	Gravity    float64
	player     *Player
	entityList []string
	entities   *components.Map
}

func (g *Game) filteredEntities(types ...components.Type) []string {
	var IDs []string
	for _, ID := range g.entityList {
		if g.entities.HasComponents(ID, types...) {
			IDs = append(IDs, ID)
		}
	}
	return IDs
}

var playerID = "abc123"

func NewGame() Game {
	g := Game{
		Gravity:  0.18,
		entities: components.NewMap(),
	}

	hitbox := gfx.R(6, 10, 26, 26)
	g.entities.Add(playerID, components.Hitbox{hitbox})
	g.entities.Add(playerID, components.Pos{gfx.V(70, 170)})
	g.entities.Add(playerID, components.Velocity{gfx.V(0, 0)})
	g.entities.Add(playerID, components.Drawable{pImage})
	g.entities.Add(playerID, components.Direction{1.0})
	playerFile.Play("stand right")
	g.entities.Add(playerID, components.Animated{playerFile})

	hitbox = gfx.R(0, 0, 32, 32)
	x, y := 0, 1
	boxID := "cdf321"
	g.entities.Add(boxID, components.Hitbox{hitbox})
	g.entities.Add(boxID, components.Pos{gfx.V(70, 220)})
	g.entities.Add(boxID, components.Drawable{tileImage.SubImage((image.Rect(32*x, 32*y, 32*(x+1), 32*(y+1)))).(*ebiten.Image)})

	boxID2 := "cdf322"
	g.entities.Add(boxID2, components.Hitbox{hitbox})
	g.entities.Add(boxID2, components.Pos{gfx.V(70+32, 220)})
	g.entities.Add(boxID2, components.Drawable{tileImage.SubImage((image.Rect(32*x, 32*y, 32*(x+1), 32*(y+1)))).(*ebiten.Image)})

	boxID3 := "cdf323"
	g.entities.Add(boxID3, components.Hitbox{hitbox})
	g.entities.Add(boxID3, components.Pos{gfx.V(70+2*32, 220)})
	g.entities.Add(boxID3, components.Drawable{tileImage.SubImage((image.Rect(32*x, 32*y, 32*(x+1), 32*(y+1)))).(*ebiten.Image)})

	x, y = 1, 1
	boxID4 := "cdf324"
	g.entities.Add(boxID4, components.Hitbox{hitbox})
	g.entities.Add(boxID4, components.Pos{gfx.V(70+2*32, 220-32)})
	g.entities.Add(boxID4, components.Drawable{tileImage.SubImage((image.Rect(32*x, 32*y, 32*(x+1), 32*(y+1)))).(*ebiten.Image)})

	x, y = 1, 0
	boxID5 := "cdf325"
	g.entities.Add(boxID5, components.Hitbox{hitbox})
	g.entities.Add(boxID5, components.Pos{gfx.V(70, 220-3*32)})
	g.entities.Add(boxID5, components.Drawable{tileImage.SubImage((image.Rect(32*x, 32*y, 32*(x+1), 32*(y+1)))).(*ebiten.Image)})

	g.entityList = []string{playerID, boxID, boxID2, boxID3, boxID4, boxID5}
	return g
}

// Todo, handle Direction properly
