package main

import (
	"net/http"
	"image"
	"image/color"
	"image/jpeg"
	"strings"
	"strconv"
	"fmt"
)

func init() {
	http.HandleFunc("/", home)
}

func home(w http.ResponseWriter, r *http.Request) {
	width, height := getDimension(r, 100, 100)
	if (width > 4000 || height > 4000){
		fmt.Fprint(w,"Dismension is too big")
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{204, 204, 204, 1})
		}
	}

	w.Header().Set("Content-Type", "image/jpg")
	jpeg.Encode(w, img, nil)
}

func getDimension(r *http.Request, defaultWidth int, defaultHeight int) (int, int) {
	path := strings.Trim(r.URL.Path,"/");
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
