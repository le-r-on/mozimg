package main

import(
	"os"
	_ "io"
	"image"
	_ "io/ioutil"
    "image/color"
	_ "image/jpeg"
)

// generate mosaic given a target image, and an array of tile images
func generateMosaic(target image.Image, tiles []image.Image) color.YCbCr {
	_ := getTileIndex(tiles)
	return target
}


// get image objects (color.YCbCr) given file paths
func getImageObject(images []string) []image.Image {
	objects := make([]image.Image, len(images))
	for i, obj := range(images) {
		file, ferr := os.Open(obj)
		if ferr == nil {
			img, _, ierr := image.Decode(file)
			if ierr == nil {
				objects[i] = img
			}
		}
	}
	return objects
}


// create index mapping image object to its average YCbCr value
func getTileIndex(tiles []image.Image) map[image.Image]color.YCbCr {
	index := make(map[image.Image]color.YCbCr)
	for _, tile := range(tiles) {
		bounds := tile.Bounds()
        ar, ag, ab := 0.0, 0.0, 0.0
		numPix := float64((bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y))
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y:= bounds.Min.Y; y < bounds.Max.Y; y++ {
                pixel := tile.At(x, y)
                r, g, b, _ := pixel.RGBA()
                ar += float64(r)
                ag += float64(g)
                ab += float64(b)
			}
		}
        y, cb, cr := color.RGBToYCbCr(uint8(ar/numPix), uint8(ag/numPix), uint8(ab/numPix))
		average := color.YCbCr{Y: y, Cb: cb, Cr: cr}
        index[tile] = average
	}
	return index
}


// main
func main() {
    // read target image and tile images paths from the
    // command line
	images := os.Args[1:]
	objects := getImageObject(images)
	generateMosaic(objects[0], objects[0:])
}
