package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

func main() {
	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			log.Fatal(err)
		}
		var filetype []string
		if filetype = strings.Split(file.Header["Content-Type"][0], "/"); filetype[0] != "image" {
			log.Fatal("error")
		}

		img, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer img.Close()

		var temp image.Image
		switch filetype[1] {
		case "jpeg":
			temp, err = jpeg.Decode(img)
			if err != nil {
				log.Fatal(err)
			}
			break
		case "png":
			temp, err = png.Decode(img)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
		log.Printf("weight,height : %T\n", temp.Bounds().Max.Y)

		_ = resizeImage(temp)
		_ = resizeImageTwo(temp)

		return c.Status(200).JSON(&fiber.Map{"success": true})
	})

	if err := app.Listen(":5000"); err != nil {
		log.Fatal(err)
	}
}

func resizeImage(file image.Image) image.Image {
	// // open file
	// file, err := os.Open("test0.jpg")
	// if err != nil {
	// 	log.Println("Error open file :", err)
	// }

	// // decode jpeg into image.Image
	// img, err := jpeg.Decode(file)
	// if err != nil {
	// 	log.Println("Error decode file :", err)
	// }
	// file.Close()

	// lib -> nfnt/resize
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	resizeImage := resize.Resize(1000, 0, file, resize.Lanczos3)
	thumbnail := resize.Thumbnail(1000, 1000, file, resize.Lanczos3)

	out, err := os.Create("test_resized_lanczos_0.jpg")
	if err != nil {
		log.Println("Error create resize file :", err)
	}
	defer out.Close()

	outThumbnail, err := os.Create("test_resized_thumbnail_lanczos_0.jpg")
	if err != nil {
		log.Println("Error create thumbmnail file :", err)
	}
	defer outThumbnail.Close()

	// write new image to file
	jpeg.Encode(out, resizeImage, nil)
	jpeg.Encode(outThumbnail, thumbnail, nil)

	return resizeImage
}

func resizeImageTwo(file image.Image) *image.NRGBA {
	// // lib -> disintegration/imaging
	// // Open a test image.
	// src, err := imaging.Open("test0.jpg")
	// if err != nil {
	// 	log.Println("Error open file :", err)
	// }

	// // Resize srcImage to size = 128x128px using the Lanczos filter.
	// dstImage1000 := imaging.Resize(src, 1000, 1000, imaging.Lanczos)

	// Resize srcImage to width = 800px preserving the aspect ratio.
	dstImage1000 := imaging.Resize(file, 1000, 0, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err := imaging.Save(dstImage1000, "out_example_lanczos.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	return dstImage1000
}
