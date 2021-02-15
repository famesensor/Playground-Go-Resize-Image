package main

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

func main() {
	// open file
	file, err := os.Open("test0.jpg")
	if err != nil {
		log.Println("Error open file :", err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Println("Error decode file :", err)
	}
	file.Close()

	// lib -> nfnt/resize
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	resizeImage := resize.Resize(1000, 0, img, resize.NearestNeighbor)
	thumbnail := resize.Thumbnail(1000, 1000, img, resize.NearestNeighbor)

	out, err := os.Create("test_resized_Nearest_0.jpg")
	if err != nil {
		log.Println("Error create resize file :", err)
	}
	defer out.Close()

	outThumbnail, err := os.Create("test_resized_thumbnail_Nearest_0.jpg")
	if err != nil {
		log.Println("Error create thumbmnail file :", err)
	}
	defer outThumbnail.Close()

	// write new image to file
	jpeg.Encode(out, resizeImage, nil)
	jpeg.Encode(outThumbnail, thumbnail, nil)

	// lib -> disintegration/imaging
	// Open a test image.
	src, err := imaging.Open("test0.jpg")
	if err != nil {
		log.Println("Error open file :", err)
	}

	// // Resize srcImage to size = 128x128px using the Lanczos filter.
	// dstImage1000 := imaging.Resize(src, 1000, 1000, imaging.Lanczos)

	// Resize srcImage to width = 800px preserving the aspect ratio.
	dstImage1000 := imaging.Resize(src, 1000, 0, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err = imaging.Save(dstImage1000, "out_example_preser.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}
