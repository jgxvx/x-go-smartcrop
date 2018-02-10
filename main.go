package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/jgxvx/smartcrop"
	"gocv.io/x/gocv"
)

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()
	fileName := flag.Arg(0)
	width, _ := strconv.Atoi(flag.Arg(1))
	height, _ := strconv.Atoi(flag.Arg(2))
	var topCrop image.Rectangle

	f, err := os.Open(fileName)
	check(err)

	img, _, err := image.Decode(f)
	check(err)

	dc := gg.NewContextForImage(img)
	dc.SetLineWidth(3)

	faces := faceDetect(img, fileName)

	fmt.Printf("Found %d face(s)...\n", len(faces))

	for _, r := range faces {
		drawRectangleOnImage(img, r, color.RGBA{255, 0, 0, 255}, dc)
		fmt.Printf("%v\n", r)
	}

	if len(faces) > 0 {
		faceBounds := rectangleBounds(faces)
		fmt.Printf("Faces bounds: %v\n", faceBounds)
		drawRectangleOnImage(img, faceBounds, color.Black, dc)
		topCrop = topCropWithoutBoost(img, width, height)
	} else {
		topCrop = topCropWithoutBoost(img, width, height)
	}

	fmt.Printf("Crop: %v\n", topCrop)

	drawRectangleOnImage(img, topCrop, color.White, dc)

	dc.SavePNG("out.png")

	//croppedimg := img.(subImager).SubImage(topCrop)
	//scaledImg := resize.Thumbnail(uint(width), uint(height), croppedimg, resize.Bicubic)

	//buf := new(bytes.Buffer)
	//options := jpeg.Options{Quality: 70}

	//err = jpeg.Encode(buf, scaledImg, &options)
	//check(err)

	//err = ioutil.WriteFile("out.jpg", buf.Bytes(), 0644)
	//check(err)
}

func faceDetect(img image.Image, fileName string) []image.Rectangle {
	// Face Detection
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	classifier.Load("./data/haarcascade_frontalface_default.xml")

	mat := gocv.IMRead(fileName, gocv.IMReadColor)
	defer mat.Close()

	imgBounds := img.Bounds()
	minSize := image.Point{imgBounds.Max.X / 10, imgBounds.Max.Y / 10}
	maxSize := image.Point{imgBounds.Max.X - 10, imgBounds.Max.Y - 10}

	return classifier.DetectMultiScaleWithParams(mat, 1.1, 4, 0, minSize, maxSize)
}

/*func topCropWithBoost(img image.Image, width, height int, boost image.Rectangle) image.Rectangle {
	cropSettings := smartcrop.CropSettings{DebugMode: true}
	analyzer := smartcrop.NewAnalyzerWithCropSettings(cropSettings)
	topCrop, _ := analyzer.FindBestCrop(img, width, height)

	return topCrop
}*/

func topCropWithoutBoost(img image.Image, width, height int) image.Rectangle {
	cropSettings := smartcrop.CropSettings{DebugMode: true}
	analyzer := smartcrop.NewAnalyzerWithCropSettings(cropSettings)
	topCrop, _ := analyzer.FindBestCrop(img, width, height)

	return topCrop
}

func drawRectangleOnImage(img image.Image, rect image.Rectangle, color color.Color, dc *gg.Context) {
	dc.DrawRectangle(float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Dx()), float64(rect.Dy()))
	dc.SetColor(color)
	dc.Stroke()
}

func rectangleBounds(rects []image.Rectangle) image.Rectangle {
	first := rects[0]
	bounds := image.Rectangle{first.Min, first.Max}

	for _, rect := range rects[1:] {
		if bounds.Min.X > rect.Min.X {
			bounds.Min.X = rect.Min.X
		}

		if bounds.Min.Y > rect.Min.Y {
			bounds.Min.Y = rect.Min.Y
		}

		if bounds.Max.X < rect.Max.X {
			bounds.Max.X = rect.Max.X
		}

		if bounds.Max.Y < rect.Max.Y {
			bounds.Max.Y = rect.Max.Y
		}
	}

	return bounds
}
