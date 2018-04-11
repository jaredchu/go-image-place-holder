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
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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
		addLabel(img,width/2,height/2,strconv.Itoa(width))

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
	col := color.RGBA{134, 134, 134, 1}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	basicFont := basicfont.Face7x13

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicFont,
		Dot:  point,
	}
	d.DrawString(label)
}
