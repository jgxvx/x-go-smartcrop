package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/muesli/smartcrop"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	args := os.Args[1:]
	fileName := args[0]

	f, err := os.Open(fileName)
	check(err)

	img, _, err := image.Decode(f)
	check(err)

	analyzer := smartcrop.NewAnalyzer()
	topCrop, _ := analyzer.FindBestCrop(img, 250, 250)

	// The crop will have the requested aspect ratio, but you need to copy/scale it yourself
	fmt.Printf("Top crop: %+v\n", topCrop)

	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	croppedimg := img.(SubImager).SubImage(topCrop)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, croppedimg)
	check(err)

	err = ioutil.WriteFile("out.png", buf.Bytes(), 0644)
	check(err)
}
