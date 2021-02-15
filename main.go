package main

import (
	"bytes"
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
		resizeFile := resizeImage(temp)
		_ = resizeImageTwo(temp)

		_, err = encodeImage(resizeFile, filetype[1])
		if err != nil {
			log.Fatal(err)
		}

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
	var resizeImage, resizeImageTwo image.Image
	if file.Bounds().Dy() > file.Bounds().Dx() {
		resizeImage = resize.Resize(0, 1000, file, resize.Lanczos3)
		resizeImageTwo = resize.Resize(0, 2000, file, resize.Lanczos3)
		// thumbnail = resize.Thumbnail(1000, 1000, file, resize.Lanczos3)
	} else {
		resizeImage = resize.Resize(1000, 0, file, resize.Lanczos3)
		resizeImageTwo = resize.Resize(2000, 0, file, resize.Lanczos3)
		// thumbnail = resize.Thumbnail(1000, 1000, file, resize.Lanczos3)
	}

	out, err := os.Create("test_resized_lanczos_0.jpg")
	if err != nil {
		log.Println("Error create resize file :", err)
	}
	defer out.Close()

	outTwo, err := os.Create("test_resized_big_lanczos_0.jpg")
	if err != nil {
		log.Println("Error create thumbmnail file :", err)
	}
	defer outTwo.Close()

	// write new image to file
	jpeg.Encode(out, resizeImage, nil)
	jpeg.Encode(outTwo, resizeImageTwo, nil)

	return resizeImage
}

func resizeImageTwo(file image.Image) *image.NRGBA {
	// // lib -> disintegration/imaging
	// // Open a test image.
	// src, err := imaging.Open("test0.jpg")
	// if err != nil {
	// 	log.Println("Error open file :", err)
	// }

	// // Resize image to size = 128x128px using the Lanczos filter.
	// dstImage1000 := imaging.Resize(src, 1000, 1000, imaging.Lanczos)

	dstImage, dstImageTwo := new(image.NRGBA), new(image.NRGBA)
	if file.Bounds().Dy() > file.Bounds().Dx() {
		// Resize image to height = 1000,2000px preserving the aspect ratio.
		dstImage = imaging.Resize(file, 0, 1000, imaging.Lanczos)
		dstImageTwo = imaging.Resize(file, 0, 2000, imaging.Lanczos)
	} else {
		// Resize image to width = 2000px preserving the aspect ratio.
		dstImage = imaging.Resize(file, 1000, 0, imaging.Lanczos)
		dstImageTwo = imaging.Resize(file, 2000, 0, imaging.Lanczos)
	}

	// Save the resulting image as JPEG.
	err := imaging.Save(dstImage, "out_example_lanczos.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	err = imaging.Save(dstImageTwo, "out_example_lanczos_big.jpg")
	if err != nil {
		log.Fatalf("failed to save image Two: %v", err)
	}
	return dstImage
}

func encodeImage(file image.Image, typesFile string) (*bytes.Buffer, error) {
	encoded := &bytes.Buffer{}
	switch typesFile {
	case "jpeg":
		if err := jpeg.Encode(encoded, file, nil); err != nil {
			return nil, err
		}
		break
	case "png":
		if err := png.Encode(encoded, file); err != nil {
			return nil, err
		}
		break
	}
	return encoded, nil
}
