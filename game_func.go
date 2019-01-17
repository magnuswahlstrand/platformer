package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/gomponents/components"
	"github.com/kyeett/gomponents/direction"
	"github.com/kyeett/tiled"
	"github.com/peterhellberg/gfx"
)

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

func (g *Game) newTrigger(o tiled.Object) {
	id := fmt.Sprintf("%d", rand.Intn(1000000))

	trigger := components.Trigger{
		Rect:      gfx.R(float64(o.X), float64(o.Y), float64(o.X+o.Width), float64(o.Y+o.Height)),
		Direction: direction.All,
	}

	for _, p := range o.Properties.Property {
		switch p.Name {
		case "scenario":
			trigger.Scenario = p.Value
		case "dir":
			trigger.Direction = direction.FromString(p.Value)
		}
	}

	// g.entities.Add(id, components)

	g.entities.Add(id, trigger)
	g.entityList = append(g.entityList, id)
	// fmt.Printf("adding trigger: %v\n", trigger)
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
	// fmt.Printf("adding teleport: %v\n", teleport)
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
		case "velocity":
			s := strings.Split(p.Value, ",")
			if len(s) < 2 {
				logrus.Error("parsed non-float values:", p.Value)
				return
			}
			fx, err := strconv.ParseFloat(s[0], 64)
			if err != nil {
				logrus.Error("parsed non-float value:", s[0])
			}
			fy, err := strconv.ParseFloat(s[1], 64)
			if err != nil {
				logrus.Error("parsed non-float value:", s[1])
			}
			g.entities.Add(id, components.Velocity{Vec: gfx.V(fx, fy)})
		}
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

	if (*counter)["lives"] <= 0 {
		g.currentScene = "lost"
		(*counter)["lives"] = 3
	}
}

func (g *Game) newPlayer() {
	hitbox := gfx.R(10, 10, 22, 26)
	g.entityList = append(g.entityList, playerID)
	g.entities.Add(playerID, components.NewHitbox(hitbox))
	g.entities.Add(playerID, components.Pos{Vec: gfx.V(0, 0)})
	g.entities.Add(playerID, components.Velocity{Vec: gfx.V(0, 0)})
	g.entities.Add(playerID, components.Drawable{pImage})
	g.entities.Add(playerID, components.Direction{1.0})
	counters := components.Counter{}
	counters["lives"] = 3
	counters["jumps"] = 2
	g.entities.Add(playerID, counters)
	playerFile.Play("stand right")
	g.entities.Add(playerID, components.Animated{playerFile})

}

func (g *Game) setPlayerPos(p gfx.Vec) {
	pos := g.entities.GetUnsafe(playerID, components.PosType).(*components.Pos)
	pos.Vec = p
}

func (g *Game) setPlayerStartingPos(p gfx.Vec) {
	initialPos = p
	g.setPlayerPos(p)
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
	g.entities.Add(id, components.NewHitbox(box))
	g.entities.Add(id, components.Pos{v})
	g.entities.Add(id, components.Drawable{tileImage.SubImage((image.Rect(32*x, 32*y, 32*(x+1), 32*(y+1)))).(*ebiten.Image)})
}

// Todo, handle Direction properly
