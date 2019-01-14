package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"strconv"

	"github.com/kyeett/gomponents/components"
	"github.com/kyeett/tiled"
	"github.com/peterhellberg/gfx"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ebitenconsole"
)

var (
	tmpImg        *ebiten.Image
	backgroundImg *ebiten.Image
	traceImg      *ebiten.Image
	scoreboardImg *ebiten.Image
)

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
	screen.DrawImage(backgroundImg, &ebiten.DrawImageOptions{})
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.updatePostMovement()

	// Draw entities
	g.drawEntities(screen)
	g.drawScoreboard(screen)
	g.drawDebugInfo(screen)
	return nil
}

type Game struct {
	Gravity       float64
	player        *Player
	entityList    []string
	entities      *components.Map
	Width, Height int
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
		Gravity:    0.18,
		entities:   components.NewMap(),
		entityList: []string{},
	}

	tmxPath := "../tiled/world6.tmx"
	fmt.Println("nay")
	worldMap, err := tiled.MapFromFile(tmxPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("yay")

	g.Width, g.Height = worldMap.Size()
	fmt.Println(g.Width, g.Height)

	traceImg, _ = ebiten.NewImage(g.Width, g.Height, ebiten.FilterDefault)
	tmpImg, _ = ebiten.NewImage(g.Width, g.Height, ebiten.FilterDefault)
	scoreboardImg, _ = ebiten.NewImage(g.Width, 16, ebiten.FilterDefault)
	scoreboardImg.Fill(color.Black)

	img, err := worldMap.LoadImage(0)
	if err != nil {
		log.Fatal("decode image: %s")
	}

	sImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	for _, layer := range worldMap.FilteredLayers() {
		if layer.Name != "background" {
			continue
		}

		img := gfx.NewImage(g.Width, g.Height, color.Transparent)
		for _, t := range worldMap.LayerTiles(layer) {
			sRect := image.Rect(t.SrcX, t.SrcY, t.SrcX+t.Width, t.SrcY+t.Height)
			dstRect := image.Rect(t.X, t.Y, g.Width+100, g.Height)
			draw.Draw(img, dstRect, sImg.SubImage(sRect), image.Pt(t.SrcX, t.SrcY), draw.Over)
		}

		backgroundImg, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

		if err != nil {
			log.Fatal(err)
		}

	}

	for _, layer := range worldMap.FilteredLayers() {
		if layer.Name == "background" {
			continue
		}

		for _, t := range worldMap.LayerTiles(layer) {

			sRect := image.Rect(t.SrcX, t.SrcY, t.SrcX+t.Width, t.SrcY+t.Height)

			// box := gfx.R(0, 0, 32, 32)
			id := fmt.Sprintf("%d", rand.Intn(10000))
			g.entities.Add(id, components.Pos{gfx.V(float64(t.X), float64(t.Y))})
			g.entities.Add(id, components.Drawable{sImg.SubImage(sRect).(*ebiten.Image)})
			// fmt.Println("Adding", t.X, t.Y)

			// fmt.Println("aaaa")
			for _, o := range t.Objectgroup.Objects {

				box := gfx.R(float64(o.X), float64(o.Y), float64(o.X+o.Width), float64(o.Y+o.Height))
				b := components.NewHitbox(box)

				// Check for special hitboxes
				for _, p := range o.Properties.Property {
					if p.Name == "allow_from_down" && p.Value == "true" {
						b.Properties[p.Name] = true
					}
				}
				g.entities.Add(id, b)

				fmt.Println("Adding at", box)
			}

			g.parseTileProperty(id, t.Properties.Property)

			g.entityList = append(g.entityList, id)
		}
	}

	for _, o := range worldMap.FilteredObjectsType() {
		fmt.Println("filtered", o)
		g.newPlayer(gfx.V(float64(o.X), float64(o.Y)))
	}

	return g
}

func (g *Game) parseTileProperty(id string, props []tiled.Property) {
	for _, p := range props {
		switch p.Name {
		case "hazard":
			val, _ := strconv.ParseBool(p.Value)
			if val {
				g.entities.Add(id, components.Hazard{})
			}
		case "bouncy":
			val, _ := strconv.ParseBool(p.Value)
			if val {
				g.entities.Add(id, components.Bouncy{})
			}
		case "killable":
			val, _ := strconv.ParseBool(p.Value)
			if val {
				g.entities.Add(id, components.Killable{})
			}
		}
		fmt.Println(p.Name, p.Value)
	}
}

var initialPos gfx.Vec

func (g *Game) Reset() {
	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	v := g.entities.GetUnsafe(playerID, components.VelocityType).(*components.Velocity)
	counter := g.entities.GetUnsafe(playerID, components.CounterType).(*components.Counter)

	pos.Vec = initialPos
	v.Vec = gfx.V(0, 0)
	(*counter)["lives"]--

}

func (g *Game) newPlayer(pos gfx.Vec) {
	initialPos = pos
	hitbox := gfx.R(6, 10, 26, 26)
	g.entityList = append(g.entityList, playerID)
	g.entities.Add(playerID, components.NewHitbox(hitbox))
	g.entities.Add(playerID, components.Pos{pos})
	g.entities.Add(playerID, components.Velocity{gfx.V(0, 0)})
	g.entities.Add(playerID, components.Drawable{pImage})
	g.entities.Add(playerID, components.Direction{1.0})
	counters := components.Counter{}
	counters["lives"] = 3
	g.entities.Add(playerID, counters)
	playerFile.Play("stand right")
	g.entities.Add(playerID, components.Animated{playerFile})

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
	g.entities.Add(id, components.NewHitbox(box))
	g.entities.Add(id, components.Pos{v})
	g.entities.Add(id, components.Drawable{tileImage.SubImage((image.Rect(32*x, 32*y, 32*(x+1), 32*(y+1)))).(*ebiten.Image)})
}

// Todo, handle Direction properly
