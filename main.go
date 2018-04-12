package main

import (
	"net/http"
	"image"
	"image/color"
	"image/jpeg"
	"strings"
	"strconv"
	"fmt"
	"google.golang.org/appengine/memcache"
	"bytes"
	"google.golang.org/appengine"
	"github.com/golang/freetype"
	"image/draw"
	"log"
	"io/ioutil"
	"flag"
)

var (
	dpi      = flag.Float64("dpi", 100, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "./luxisr.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 24, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

func init() {
	http.HandleFunc("/", home)
}

func home(w http.ResponseWriter, r *http.Request) {
	width, height := getDimension(strings.Trim(r.URL.Path, "/"), 100, 100)
	if (width > 4000 || height > 4000) {
		fmt.Fprint(w, "Image size is too big")
		return
	}

	debug := r.FormValue("debug")
	key := strconv.Itoa(width + height)
	context := appengine.NewContext(r)
	if item, err := memcache.Get(context, key); len(debug) > 0 || err == memcache.ErrCacheMiss {
		// Not cached yet
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				img.Set(x, y, color.RGBA{204, 204, 204, 1})
			}
		}
		addLabel(img,10,10,strconv.Itoa(width))

		// cache the image
		item := &memcache.Item{
			Key:   key,
			Value: imgToByte(img),
		}
		if err := memcache.Set(context, item); err != nil {
			fmt.Fprint(w, err)
			return
		}

		w.Header().Set("Content-Type", "image/jpg")
		jpeg.Encode(w, img, nil)

	} else if err != nil {
		fmt.Fprint(w, err)
		return
	} else {
		// got item from memcache
		w.Header().Set("Content-Type", "image/jpg")
		w.Write(item.Value);
	}
}

func getDimension(path string, defaultWidth int, defaultHeight int) (int, int) {
	dArray := strings.Split(path, "x")
	if (len(dArray) == 2) {
		width, err := strconv.Atoi(dArray[0])
		if (err != nil) {
			return defaultWidth, defaultHeight
		}

		height, err := strconv.Atoi(dArray[1])
		if (err != nil) {
			return defaultWidth, defaultHeight
		}

		return width, height
	} else {
		width, err := strconv.Atoi(dArray[0])
		if (err != nil) {
			return defaultWidth, defaultHeight
		}

		return width, width
	}
}

func imgToByte(img image.Image) []byte {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return []byte("")
	}
	return buf.Bytes()
}

func addLabel(img *image.RGBA, x, y int, label string) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	fg, bg := image.NewUniform(color.RGBA{152,152,152,1}), img
	draw.Draw(img, img.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)

	size := 18.0 // font size in pixels
	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))

	if _, err := c.DrawString(label, pt); err != nil {
		// handle error
	}
}
