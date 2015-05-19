package main

import(
	"os"
	"fmt"
	_ "io"
	"sort"
	"image"
	_ "io/ioutil"
    "image/color"
	_ "image/jpeg"
)


// datastructure for sorting image.Image objects
type sortedMap struct {
	rgbm map[image.Image]color.RGBA
	ycbcrm map[image.Image]color.YCbCr
	i []image.Image
}


func (sm *sortedMap) Len() int {
	return len(sm.i)
}


func (sm *sortedMap) Less(j, k int) bool {
	return sm.ycbcrm[sm.i[j]].Y > sm.ycbcrm[sm.i[k]].Y
}


func (sm *sortedMap) Swap(j, k int) {
	sm.i[j], sm.i[k] = sm.i[k], sm.i[j]
}


// generate mosaic given a target image, and an array of tile images
func generateMosaic(target image.Image, tiles []image.Image) image.Image {
	index := getTileIndex(tiles)
	sort.Sort(index)
	return target
}


// get the closest tile for a given value
func getSimilarTile(value color.YCbCr, index *sortedMap) image.Image {
	images := index.i
	for len(images) > 1 {
		mid := uint32(len(images)/2)
		if value.Y > index.ycbcrm[images[mid]].Y {
			images = images[mid:]
		} else {
			images = images[:mid]
		}
	}
	return images[0]
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


// create index mapping image object to its average RGBA value
func getTileIndex(tiles []image.Image) *sortedMap {
	index := new(sortedMap)
	index.rgbm = make(map[image.Image]color.RGBA)
	index.ycbcrm = make(map[image.Image]color.YCbCr)
	index.i = make([]image.Image, len(tiles))
	for ind, tile := range(tiles) {
		bounds := tile.Bounds()
        ar, ag, ab, aa := 0.0, 0.0, 0.0, 0.0
		numPix := float64((bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y))
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y:= bounds.Min.Y; y < bounds.Max.Y; y++ {
                pixel := tile.At(x, y)
                r, g, b, a := pixel.RGBA()
                ar += float64(r)
                ag += float64(g)
                ab += float64(b)
				aa += float64(a)
			}
		}
		averageRGBA := color.RGBA{R: uint8(ar/numPix), G: uint8(ag/numPix),
		                          B: uint8(ab/numPix), A: uint8(aa/numPix)}
		y, cb, cr := color.RGBToYCbCr(uint8(ar/numPix), uint8(ag/numPix), uint8(ab/numPix))
		averageYCbCr := color.YCbCr{Y: y, Cb: cb, Cr: cr}
		index.rgbm[tile] = averageRGBA
		index.ycbcrm[tile] = averageYCbCr
		index.i[ind] = tile
	}
	return index
}


// main
func main() {
    // read target image and tile images paths from the
    // command line
	images := os.Args[1:]
	objects := getImageObject(images)
	generateMosaic(objects[0], objects[1:])
}
