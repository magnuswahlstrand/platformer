package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type resource struct {
	path, name string
}

func main() {
	resources := []resource{
		resource{
			path: "assets/graphics/character.json",
			name: "Character_json",
		},
		resource{
			path: "assets/graphics/tiles.png",
			name: "Tiles_png",
		},
		resource{
			path: "assets/graphics/misc.png",
			name: "Misc_png",
		},
		resource{
			path: "assets/graphics/character.png",
			name: "Character_png",
		},
	}

	f, err := os.Create("assets/graphics.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "package resources")
	fmt.Fprintln(f, "")
	for _, r := range resources {
		content, err := ioutil.ReadFile(r.path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(f, "var %s    = []byte(%q)\n", r.name, string(content))
	}
}
