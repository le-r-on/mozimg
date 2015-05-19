package main

import (
    "github.com/gocraft/web"
    "github.com/polds/imgbase64"
    _ "github.com/franela/goreq"
    "io"
    "html/template"
    _ "fmt"
    _ "log"
    "bytes"
    "image"
    _ "image/jpeg"
    "image/png"
    _ "image/color"
    "bufio"
)

func getAvgColorFromImg(imageObj image.Image) image.Image {
    avgColor, _ := getAverageColor(imageObj)
    avgColorImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{32, 32}})
    for y := 0; y < 32; y++ {
        for x := 0; x < 32; x++ {
            avgColorImg.SetRGBA(x, y, avgColor)       
        }
    }
    return avgColorImg
}

func imageToB64(imageObj image.Image) string {
    tmp_buf := new(bytes.Buffer)
    png.Encode(tmp_buf, imageObj)
    avgColorB64 := imgbase64.FromBuffer(*tmp_buf)
    return avgColorB64
}

func displayImgObjAndAvg(rw web.ResponseWriter, imageObj image.Image) {
    origImgB64 := imageToB64(imageObj)
    avgColorB64 := imageToB64(getAvgColorFromImg(imageObj))
    base_tmpl.Execute(
        rw,
        &Context{Image: template.URL(origImgB64), AvgColor: template.URL(avgColorB64)})
}

func displayImgObjAndOrig(rw web.ResponseWriter, imageObj image.Image, origImageObj image.Image) {
    origImgB64 := imageToB64(imageObj)
    origImageObjB64 := imageToB64(origImageObj)
    base_tmpl.Execute(
        rw,
        &Context{Image: template.URL(origImgB64), AvgColor: template.URL(origImageObjB64)})
}

func imageFromReader(reader io.Reader) image.Image {
    img, _, _ := image.Decode(reader)
    return img
}

func imageFromBuffer(buffer bytes.Buffer) image.Image {
    return imageFromReader(bufio.NewReader(&buffer))
}
