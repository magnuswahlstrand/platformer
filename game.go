package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/kyeett/gomponents/components"
	"github.com/kyeett/tiled"
	"github.com/peterhellberg/gfx"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten"
)

var (
	tmpImg        *ebiten.Image
	backgroundImg *ebiten.Image
	foregroundImg *ebiten.Image
	visionImg     *ebiten.Image
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

func (g *Game) drawHitboxes(screen *ebiten.Image) {
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

func (g *Game) drawPlayerVision(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	// opt.Address = ebiten.AddressRepeat

	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	opt.GeoM.Translate(pos.X-float64(visionImg.Bounds().Dx())/2+15, pos.Y-float64(visionImg.Bounds().Dy())/2+20)
	opt.CompositeMode = ebiten.CompositeModeDestinationOut
	tmp, _ := ebiten.NewImageFromImage(foregroundImg, ebiten.FilterDefault)
	// tmp, _ := ebiten.NewImage(200, 200, ebiten.FilterDefault)
	// tmp.Fill(colornames.Red)
	// backgroundImg.DrawImage(visionImg, &ebiten.DrawImageOptions{})

	tmp.DrawImage(visionImg, opt)
	screen.DrawImage(tmp, &ebiten.DrawImageOptions{})
}

type Game struct {
	Gravity       float64
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

func NewGame(worldFile string) Game {
	g := Game{
		Gravity:    gravityConst,
		entities:   components.NewMap(),
		entityList: []string{},
	}

	fmt.Println("nay")
	worldMap, err := tiled.MapFromFile(worldFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("yay")

	g.Width, g.Height = worldMap.Size()

	traceImg, _ = ebiten.NewImage(g.Width, g.Height, ebiten.FilterDefault)
	backgroundImg, _ = ebiten.NewImage(g.Width, g.Height, ebiten.FilterDefault)
	foregroundImg, _ = ebiten.NewImage(g.Width, g.Height, ebiten.FilterDefault)

	tmpImg, _ = ebiten.NewImage(g.Width, g.Height, ebiten.FilterDefault)
	scoreboardImg, _ = ebiten.NewImage(g.Width, 16, ebiten.FilterDefault)
	scoreboardImg.Fill(color.Black)

	tmpImg2 := gfx.NewImage(200, 200, color.Transparent)
	gfx.DrawCircle(tmpImg2, gfx.V(100, 100), 50, 0, color.White)
	visionImg, _ = ebiten.NewImageFromImage(tmpImg2, ebiten.FilterDefault)

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
		if layer.Name != "foreground" {
			continue
		}

		img := gfx.NewImage(g.Width, g.Height, color.Transparent)
		for _, t := range worldMap.LayerTiles(layer) {
			sRect := image.Rect(t.SrcX, t.SrcY, t.SrcX+t.Width, t.SrcY+t.Height)
			dstRect := image.Rect(t.X, t.Y, g.Width+100, g.Height)
			draw.Draw(img, dstRect, sImg.SubImage(sRect), image.Pt(t.SrcX, t.SrcY), draw.Over)
		}

		foregroundImg, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

		if err != nil {
			log.Fatal(err)
		}
	}

	for _, layer := range worldMap.FilteredLayers() {
		if layer.Name != "hitboxes" {
			continue
		}

		for _, t := range worldMap.LayerTiles(layer) {
			id := fmt.Sprintf("%d", rand.Intn(10000))
			g.entities.Add(id, components.Pos{gfx.V(float64(t.X), float64(t.Y))})
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
			}
			g.parseTileProperty(id, t.Properties.Property)
			g.entityList = append(g.entityList, id)
		}
	}

	for _, layer := range worldMap.FilteredLayers() {
		if layer.Name == "background" || layer.Name == "hitboxes" || layer.Name == "foreground" {
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

		switch o.Type {
		case "player":
			g.newPlayer(gfx.V(float64(o.X), float64(o.Y)))
		case "teleport":
			g.newTeleport(o)
		case "trigger":
			g.newTrigger(o)
		}
	}

	return g
}

func parseDirections(s string) byte {
	if s == "" {
		return components.DirUp | components.DirDown | components.DirLeft | components.DirRight
	}

	var dir byte
	if strings.Contains(s, "U") {
		dir |= components.DirUp
	}
	if strings.Contains(s, "D") {
		dir |= components.DirDown
	}
	if strings.Contains(s, "L") {
		dir |= components.DirLeft
	}
	if strings.Contains(s, "R") {
		dir |= components.DirRight
	}
	return dir
}

func (g *Game) newTrigger(o tiled.Object) {
	id := fmt.Sprintf("%d", rand.Intn(1000000))

	trigger := components.Trigger{
		Rect: gfx.R(float64(o.X), float64(o.Y), float64(o.X+o.Width), float64(o.Y+o.Height)),
	}

	for _, p := range o.Properties.Property {
		switch p.Name {
		case "scenario":
			trigger.Scenario = p.Value
		case "dir":
			trigger.Direction = parseDirections(p.Value)
		}
	}

	// g.entities.Add(id, components)

	g.entities.Add(id, trigger)
	g.entityList = append(g.entityList, id)
	fmt.Printf("adding trigger: %v\n", trigger)
}

func (g *Game) newTeleport(o tiled.Object) {

	id := fmt.Sprintf("%d", rand.Intn(1000000))

	teleport := components.Teleporting{
		Name: o.Name,
		Pos:  gfx.V(float64(o.X), float64(o.Y)),
	}

	for _, p := range o.Properties.Property {
		switch p.Name {
		case "target":
			teleport.Target = p.Value
		case "dx":
			dx, _ := strconv.Atoi(p.Value)
			teleport.Pos.X = float64(o.X + dx)
		case "dy":
			dy, _ := strconv.Atoi(p.Value)
			teleport.Pos.Y = float64(o.Y + dy)
		}
	}

	g.entities.Add(id, teleport)
	g.entities.Add(id, components.Pos{Vec: gfx.V(float64(o.X), float64(o.Y))})
	g.entities.Add(id, components.NewHitbox(gfx.R(0, 0, float64(o.Width), float64(o.Height))))
	g.entityList = append(g.entityList, id)
	fmt.Printf("adding teleport: %v\n", teleport)
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
