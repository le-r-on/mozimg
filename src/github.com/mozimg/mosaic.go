package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "io"
	"os"
	"sort"
)

// datastructure for sorting image.Image objects
type sortedMap struct {
	rgbm   map[image.Image]color.RGBA
	ycbcrm map[image.Image]color.YCbCr
	i      []image.Image
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
func generateMosaic(target image.Image, tiles []image.Image, rows int, columns int) image.Image {
	index := getTileIndex(tiles)
	sort.Sort(index)
	bounds := target.Bounds()
	x_length, y_length := int(bounds.Max.X/columns), int(bounds.Max.Y/rows)

	ycbcrImg := target.(*image.YCbCr)
	outImg := YCbCrToRGBA(ycbcrImg)

	// iterate through target image's cells and get tile
	for x := 0; x < bounds.Max.X-x_length; x += x_length {
		for y := 0; y < bounds.Max.Y-y_length; y += y_length {
			rect, dp := image.Rect(x, x+x_length, y, y+y_length), image.Point{X: x, Y: y}
			cell := ycbcrImg.SubImage(rect)
			_, averageYCbCr := getAverageColor(cell)
			tile := getSimilarTile(averageYCbCr, index)
			tile = resize.Resize(uint(x_length), uint(y_length), tile, resize.Lanczos3)
			r := image.Rectangle{dp, dp.Add(tile.Bounds().Size())}
			draw.Draw(outImg, r, tile, tile.Bounds().Min, draw.Src)
		}
	}
	out, _ := os.Create("./output.jpg")
	defer out.Close()
	var opt jpeg.Options
	opt.Quality = 80
	jpeg.Encode(out, outImg, &opt)

	return target
}

// get the closest tile for a given value
func getSimilarTile(value color.YCbCr, index *sortedMap) image.Image {
	images := index.i
	for len(images) > 1 {
		mid := uint32(len(images) / 2)
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
	for i, obj := range images {
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
	for ind, tile := range tiles {
		averageRGBA, averageYCbCr := getAverageColor(tile)
		index.rgbm[tile] = averageRGBA
		index.ycbcrm[tile] = averageYCbCr
		index.i[ind] = tile
	}
	return index
}

// get average RGBA, YCbCr color values given an image
func getAverageColor(img image.Image) (color.RGBA, color.YCbCr) {
	bounds := img.Bounds()
	ar, ag, ab, aa := 0.0, 0.0, 0.0, 0.0
	numPix := float64((bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y))
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixel := img.At(x, y)
			r, g, b, a := pixel.RGBA()
			ar += float64(r)
			ag += float64(g)
			ab += float64(b)
			aa += float64(a)
		}
	}
	r, g, b, a := uint8(ar/numPix), uint8(ag/numPix), uint8(ab/numPix), uint8(aa/numPix)
	averageRGBA := color.RGBA{R: r, G: g, B: b, A: a}
	y, cb, cr := color.RGBToYCbCr(r, g, b)
	averageYCbCr := color.YCbCr{Y: y, Cb: cb, Cr: cr}

	return averageRGBA, averageYCbCr
}

// Convert YCbCr image to RGBA
func YCbCrToRGBA(src *image.YCbCr) *image.RGBA {
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	rgb := image.NewRGBA(bounds)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)
			rgbColor := rgb.ColorModel().Convert(oldColor)
			rgb.Set(x, y, rgbColor)
		}
	}
	out, _ := os.Create("./output.jpg")
	defer out.Close()
	var opt jpeg.Options
	opt.Quality = 80
	jpeg.Encode(out, rgb, &opt)
	return rgb
}

// main
func main() {
	// read target image and tile images paths from the
	// command line
	images := os.Args[1:]
	objects := getImageObject(images)
	generateMosaic(objects[0], objects[1:], 10, 10)
}
