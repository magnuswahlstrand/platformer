package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/kyeett/gomponents/components"
	resources "github.com/kyeett/platformer/assets"
	"github.com/kyeett/tiled"
	"github.com/peterhellberg/gfx"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten"
)

var (
	sImg          *ebiten.Image
	tmpImg        *ebiten.Image
	backgroundImg *ebiten.Image
	foregroundImg *ebiten.Image
	visionImg     *ebiten.Image
	tmpVisionImg  *ebiten.Image
	traceImg      *ebiten.Image
	scoreboardImg *ebiten.Image
)

type Game struct {
	Gravity       float64
	currentScene  string
	scenes        map[string]func(*Game, *ebiten.Image) error
	entityList    []string
	entities      *components.Map
	Width, Height int
	baseDir       string
}

func NewGame(worldFile string) Game {

	g := Game{
		currentScene: "game",
		scenes: map[string]func(*Game, *ebiten.Image) error{
			"game":    gameLoop,
			"victory": victoryScreen,
			"lost":    lostScreen,
		},
		Gravity:    gravityConst,
		entities:   components.NewMap(),
		entityList: []string{},
		baseDir:    filepath.Dir(worldFile),
	}
	g.newPlayer()

	worldMap := g.loadWorldMap(worldFile)

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

	g.initializeWorld(worldMap)
	return g
}

func (g *Game) initializeWorld(worldMap *tiled.Map) {
	g.Width, g.Height = worldMap.Size()
	camera, _ = ebiten.NewImageFromImage(gfx.NewImage(g.Width, g.Height, colornames.Red), ebiten.FilterDefault)
	tmpVisionImg, _ = ebiten.NewImageFromImage(gfx.NewImage(g.Width, g.Height, color.Transparent), ebiten.FilterDefault)

	g.Width, g.Height = worldMap.Size()
	// Remove all existing entitites, except the player
	for _, e := range g.entityList {
		if e == playerID {
			continue
		}
		g.entities.RemoveAll(e)
	}
	g.entityList = []string{}

	path, err := worldMap.ImagePath(0)

	// Workaround
	path = strings.Replace(path, "../platformer/", "", -1)
	img, _, err := image.Decode(bytes.NewReader(resources.LookupFatal(path)))
	if err != nil {
		log.Fatal("decode image: %s")
	}

	sImg, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize background image, in case of no background layer
	img = gfx.NewImage(g.Width, g.Height, color.Black)
	backgroundImg, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	for _, layer := range worldMap.FilteredLayers() {
		if layer.Name != "background" {
			continue
		}
		img := gfx.NewImage(g.Width, g.Height, color.Black)
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
			g.parseTileObjects(id, t.Objectgroup.Objects)
			g.entities.Add(id, components.Pos{gfx.V(float64(t.X), float64(t.Y))})
			// for _, o := range t.Objectgroup.Objects {
			// 	box := gfx.R(float64(o.X), float64(o.Y), float64(o.X+o.Width), float64(o.Y+o.Height))
			// 	b := components.NewHitbox(box)

			// 	// Check for special hitboxes
			// 	for _, p := range o.Properties.Property {
			// 		if p.Name == "allow_from_down" && p.Value == "true" {
			// 			b.Properties[p.Name] = true
			// 		}
			// 		if p.Name == "target" && p.Value == "monster" { // Todo, make proper solution
			// 			b.Properties["monsters_only"] = true
			// 		}
			// 	}
			// 	g.entities.Add(id, b)
			// }
			g.parseTileProperty(id, t.Properties.Property)
			g.entityList = append(g.entityList, id)
		}
	}

	for _, layer := range worldMap.FilteredLayers() {
		if layer.Name == "background" || layer.Name == "hitboxes" || layer.Name == "foreground" {
			continue
		}

		fmt.Println(layer.Name)
		fmt.Println(layer.Width)
		fmt.Println(layer.Height)

		for _, t := range worldMap.LayerTiles(layer) {
			id := fmt.Sprintf("%d", rand.Intn(10000))
			g.parseTile(id, t)
		}
	}

	for _, o := range worldMap.FilteredObjectsType() {
		switch o.Type {
		case "player":
			g.setPlayerStartingPos(gfx.V(float64(o.X), float64(o.Y)))
		case "teleport":
			g.newTeleport(o)
		case "trigger":
			g.newTrigger(o)
		default:
			g.newEnemy(o, worldMap)
		}
	}

	g.entityList = append(g.entityList, playerID)
}

func (g *Game) parseTileObjects(id string, objects []tiled.Object) {
	for _, o := range objects {

		box := gfx.R(float64(o.X), float64(o.Y), float64(o.X+o.Width), float64(o.Y+o.Height))
		b := components.NewHitbox(box)

		// Check for special hitboxes
		for _, p := range o.Properties.Property {
			if p.Name == "allow_from_down" && p.Value == "true" {
				b.Properties[p.Name] = true
			}
			if p.Name == "target" && p.Value == "monster" { // Todo, make proper solution
				b.Properties["monsters_only"] = true
			}
		}
		g.entities.Add(id, b)
	}
}

func (g *Game) parseTile(id string, t tiled.TileData) {
	sRect := image.Rect(t.SrcX, t.SrcY, t.SrcX+t.Width, t.SrcY+t.Height)

	// box := gfx.R(0, 0, 32, 32)
	g.entities.Add(id, components.Pos{gfx.V(float64(t.X), float64(t.Y))})
	g.entities.Add(id, components.Drawable{sImg.SubImage(sRect).(*ebiten.Image)})
	g.parseTileObjects(id, t.Objectgroup.Objects)
	g.parseTileProperty(id, t.Properties.Property)
	g.entityList = append(g.entityList, id)
}

func (g *Game) newEnemy(o tiled.Object, worldMap *tiled.Map) {

	id := fmt.Sprintf("%d", rand.Intn(1000000))

	// https://doc.mapeditor.org/en/stable/reference/tmx-map-format/#tile-flipping
	const flippedHorizontallyFlag uint32 = 0x80000000
	const flippedVerticallyFlag uint32 = 0x40000000
	const flippedDiagonallyFlag uint32 = 0x20000000

	gid := ^(flippedHorizontallyFlag | flippedVerticallyFlag | flippedDiagonallyFlag) & o.Gid
	t := worldMap.TileData(int(gid))
	t.X = o.X
	t.Y = o.Y - 16 // Bug: 5: Offset - 1???

	sRect := image.Rect(t.SrcX, t.SrcY, t.SrcX+t.Width, t.SrcY+t.Height)
	g.entities.Add(id, components.Pos{Vec: gfx.V(float64(t.X), float64(t.Y))})
	if t.Type != "hitbox" {
		g.entities.Add(id, components.Drawable{sImg.SubImage(sRect).(*ebiten.Image)})
	}

	g.parseTileObjects(id, t.Objectgroup.Objects)
	g.parseTileProperty(id, t.TilesetTile.Properties.Property)
	g.parseTileProperty(id, o.Properties.Property)

	g.entityList = append(g.entityList, id)
	// fmt.Printf("adding teleport: %v\n", teleport)
}
