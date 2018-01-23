package main

import (
	"bytes"
	"flag"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/muesli/smartcrop"
	"github.com/nfnt/resize"
)

func main() {
	flag.Parse()
	fileName := flag.Arg(0)
	width, _ := strconv.Atoi(flag.Arg(1))
	height, _ := strconv.Atoi(flag.Arg(2))

	f, _ := os.Open(fileName)

	img, _, _ := image.Decode(f)

	//cropSettings := smartcrop.CropSettings{DebugMode: false}
	analyzer := smartcrop.NewAnalyzer()
	topCrop, _ := analyzer.FindBestCrop(img, width, height)

	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	croppedimg := img.(SubImager).SubImage(topCrop)
	scaledImg := resize.Thumbnail(uint(width), uint(height), croppedimg, resize.Bicubic)

	buf := new(bytes.Buffer)
	options := jpeg.Options{Quality: 70}

	_ = jpeg.Encode(buf, scaledImg, &options)

	_ = ioutil.WriteFile("out.jpg", buf.Bytes(), 0644)
}
