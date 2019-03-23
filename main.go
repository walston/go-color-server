package main

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func getColor(s string, re *regexp.Regexp) *color.Color {
	m := re.FindStringSubmatch(s)

	r, _ := strconv.ParseUint(m[1][:2], 16, 8)
	g, _ := strconv.ParseUint(m[1][2:4], 16, 8)
	b, _ := strconv.ParseUint(m[1][4:], 16, 8)

	var c color.Color = color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}

	return &c
}

func respondSolid(w http.ResponseWriter, r *http.Request, re *regexp.Regexp) {
	img := image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{*getColor(r.URL.Path, re)})
	o := &gif.Options{NumColors: 1}

	img.SetColorIndex(0, 0, 0)
	w.Header().Set("Content-Type", "image/gif")
	gif.Encode(w, img, o)
}

func respondGrid(w http.ResponseWriter, r *http.Request, re *regexp.Regexp) {
	p := color.Palette{
		*getColor(r.URL.Path, re),
		color.RGBA{0x22, 0x22, 0x22, 0xff},
	}

	img := image.NewPaletted(image.Rect(0, 0, 72, 72), p)
	o := &gif.Options{NumColors: 2}

	m := img.Bounds().Max
	for i := 0; i < m.Y; i += 3 {
		img.SetColorIndex(0, i, 1)
	}

	for i := 0; i < m.X; i += 3 {
		img.SetColorIndex(i, 0, 1)
	}

	w.Header().Set("Content-Type", "image/gif")
	gif.Encode(w, img, o)
}

func respondRedirect(w http.ResponseWriter, r *http.Request) {
	prefix := ""
	if regexp.MustCompile("^/grid/").MatchString(r.URL.Path) {
		prefix = "/grid"
	}

	var c string
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		log.Panic(err)
	}
	c = hex.EncodeToString(bytes)

	fmt.Fprint(os.Stdout, c, "\n")
	http.Redirect(w, r, prefix+"/"+c[:6]+".gif", 302)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var solid = regexp.MustCompile("^/([0-9a-fA-F]{6})\\.gif$")
		var grid = regexp.MustCompile("^/grid/([0-9a-fA-F]{6})\\.gif$")

		if solid.MatchString(r.URL.Path) {
			respondSolid(w, r, solid)
		} else if grid.MatchString(r.URL.Path) {
			respondGrid(w, r, grid)
		} else {
			respondRedirect(w, r)
		}
	})
	http.ListenAndServe(":9000", nil)
}
