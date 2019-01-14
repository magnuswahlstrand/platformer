package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"

	"github.com/kyeett/gomponents/components"
	"github.com/kyeett/tiled"
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
	box1 := "cdf321"
	g.newBox(box1, gfx.V(70, 220), "green")

	box2 := "cdf322"
	g.newBox(box2, gfx.V(70+32, 220), "green")

	box3 := "cdf323"
	g.newBox(box3, gfx.V(70+2*32, 220), "green")

	box4 := "cdf324"
	g.newBox(box4, gfx.V(70+2*32, 220-32), "blue")

	box5 := "cdf325"
	g.newBox(box5, gfx.V(70, 220-3*32), "red")
	g.entityList = []string{playerID, box1, box2, box3, box4, box5}

	tmxPath := "../tiled/world6.tmx"
	worldMap, err := tiled.MapFromFile(tmxPath)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("assets/tilesheets/platformer2.png")
	if err != nil {
		log.Fatal("open image: %s")
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal("decode image: %s")
	}

	sImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(worldMap.FilteredLayers("yoyo"), sImg)

	for _, layer := range worldMap.FilteredLayers("") {

		for _, t := range worldMap.LayerTiles(layer) {

			sRect := image.Rect(t.SrcX, t.SrcY, t.SrcX+t.Width, t.SrcY+t.Height)

			// box := gfx.R(0, 0, 32, 32)
			id := fmt.Sprintf("%d", rand.Intn(10000))
			g.entities.Add(id, components.Pos{gfx.V(float64(t.X), float64(t.Y))})
			g.entities.Add(id, components.Drawable{sImg.SubImage(sRect).(*ebiten.Image)})
			// fmt.Println("Adding", t.X, t.Y)

			// fmt.Println("aaaa")
			// for _, o := range t.Objectgroup.Objects {

			// 	box := gfx.R(float64(o.X), float64(o.Y), float64(o.X+o.Width), float64(o.Y+o.Height))
			// 	g.entities.Add(id, components.Hitbox{box})
			// 	fmt.Println("Adding at", box)
			// }
			g.entityList = append(g.entityList, id)
		}
	}
	return g
}

func (g *Game) newBox(id string, v gfx.Vec, name string) {

	var x, y int
	switch name {
	case "red":
		x, y = 1, 0
	case "blue":
		x, y = 1, 1
	case "green":
		x, y = 0, 1
	default:
		log.Fatal("invalid name:", name)
	}

	box := gfx.R(0, 0, 32, 32)
	fmt.Println("Adding 2at", box)
	g.entities.Add(id, components.Hitbox{box})
	g.entities.Add(id, components.Pos{v})
	g.entities.Add(id, components.Drawable{tileImage.SubImage((image.Rect(32*x, 32*y, 32*(x+1), 32*(y+1)))).(*ebiten.Image)})
}

// Todo, handle Direction properly
