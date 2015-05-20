package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "io"
	"math"
	"os"
	"sort"
)

// datastructure for sorting image.Image objects
type sortedMap struct {
	rgbm   []color.RGBA
	ycbcrm []color.YCbCr
	i      []*image.Image
}

func (sm *sortedMap) Len() int {
	return len(sm.i)
}

func (sm *sortedMap) Less(j, k int) bool {
	return sm.ycbcrm[j].Y > sm.ycbcrm[k].Y
}

func (sm *sortedMap) Swap(j, k int) {
	sm.i[j], sm.i[k] = sm.i[k], sm.i[j]
	sm.ycbcrm[j], sm.ycbcrm[k] = sm.ycbcrm[k], sm.ycbcrm[j]
	sm.rgbm[j], sm.rgbm[k] = sm.rgbm[k], sm.rgbm[j]
}

// generate mosaic given a target image, and an array of tile images
func generateMosaic(target *image.Image, tiles []*image.Image, rows int, columns int) image.Image {
	fmt.Println("Mosaic: generating and sorting index")

	bounds := (*target).Bounds()
	x_length, y_length := int(bounds.Max.X/columns), int(bounds.Max.Y/rows)

	index := getTileIndex(tiles, uint(x_length), uint(y_length))
	sort.Sort(index)

	ycbcrImg := (*target).(*image.YCbCr)
	outImg := YCbCrToRGBA(ycbcrImg)

	fmt.Println("Mosaic: iterating through tiles")
	// iterate through target image's cells and get tile
	for x := 0; x < bounds.Max.X; x += x_length {
		fmt.Printf("Mosaic: iterating through tiles, %d left\n", (bounds.Max.X-x)/x_length)
		for y := 0; y < bounds.Max.Y; y += y_length {
			rect, dp := image.Rect(x, y, x+x_length, y+y_length), image.Point{X: x, Y: y}
			cell := ycbcrImg.SubImage(rect)
			_, averageYCbCr := getAverageColor(&cell)
			tile := *getSimilarTile(averageYCbCr, index)
			r := image.Rectangle{dp, dp.Add(tile.Bounds().Size())}
			draw.Draw(outImg, r, tile, tile.Bounds().Min, draw.Src)
		}
	}
	return outImg
}

// get the closest tile for a given value
func getSimilarTile(value color.YCbCr, index *sortedMap) *image.Image {
	images := index.i
	left, right := 0, len(images)
	for right - left > 1 {
		mid := (right + left) / 2
		if value.Y > index.ycbcrm[mid].Y {
			right = mid
		} else {
			left = mid
		}
	}
	return images[left]
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
func getTileIndex(tiles []*image.Image, x_length uint, y_length uint) *sortedMap {
	index := new(sortedMap)
	index.rgbm = make([]color.RGBA, len(tiles) + 1)
	index.ycbcrm = make([]color.YCbCr, len(tiles) + 1)
	index.i = make([]*image.Image, len(tiles) + 1)
	for ind, tile := range tiles {
		norm_tile := resize.Resize(x_length, y_length, *tile, resize.Lanczos3)
		if norm_tile == nil {
			continue
		}
		averageRGBA, averageYCbCr := getAverageColor(&norm_tile)
		index.rgbm[ind] = averageRGBA
		index.ycbcrm[ind] = averageYCbCr
		index.i[ind] = &norm_tile
	}
	return index
}

// get average RGBA, YCbCr color values given an image
func getAverageColor(img *image.Image) (color.RGBA, color.YCbCr) {
	bounds := (*img).Bounds()
	ar, ag, ab, aa := 0.0, 0.0, 0.0, 0.0
	numPix := float64((bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y))
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pixel := (*img).At(x, y)
			r, g, b, a := pixel.RGBA()
			asqrt := math.Sqrt(float64(a))
			ar += math.Floor(float64(r) / asqrt)
			ag += math.Floor(float64(g) / asqrt)
			ab += math.Floor(float64(b) / asqrt)
			aa += math.Floor(asqrt)
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
