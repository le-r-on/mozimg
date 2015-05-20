package main

import (
    "github.com/gocraft/web"
    _ "github.com/franela/goreq"
    "os"
    "html/template"
    "fmt"
    "net/http"
    _ "log"
    "image"
    _ "image/jpeg"
    _ "image/color"
    "strconv"
)

const MAX_FILE_SIZE = 1024 * 1024
const MAX_NUMBER_OF_FILES = 100

type Context struct {
    TiledImage template.URL
    OrigImage template.URL
    Message template.URL
}

func (c *Context) ShowPicture(rw web.ResponseWriter, req *web.Request) {
    file, _ := os.Open("test.png")
    defer file.Close()

    imgObj, _, err := image.Decode(file)
    if err != nil {
        error_tmpl.Execute(rw, &Context{Message: template.URL("Failed to read file:" + err.Error())})
        return
    }
    displayImgObj(rw, imgObj)
}

func (c *Context) RefreshPicture(rw web.ResponseWriter, req *web.Request) {
    req.ParseMultipartForm(1024)

    dimension, err := strconv.Atoi(req.Form["dimension"][0])
    if err != nil {
        dimension = 50
    }
    pic_num, err := strconv.Atoi(req.Form["pic_num"][0])
    if err != nil {
        pic_num = 100
    }

    fmt.Println("Getting images from elastic")
    imgs := randomThumbnails(pic_num)
    fmt.Println("Generating mosaic")
    resImage := generateMosaic(imgs[0], imgs[1:], dimension, dimension)
    fmt.Println("Displaying results")
    displayImgObjAndOrig(rw, resImage, imgs[0])
}

func (c *Context) TilePicture(rw web.ResponseWriter, req *web.Request) {
    fmt.Println(req.Form)
    err := req.ParseMultipartForm(int64(MAX_FILE_SIZE * MAX_NUMBER_OF_FILES))
    if err != nil {
        error_tmpl.Execute(rw, &Context{Message: template.URL("Failed to parse multipart")})
        return
    }

    fmt.Println("Fetching tiles from user")
    tiles := make([]image.Image, 0)
    files := req.MultipartForm.File["files"]
    for _, file := range files {
        file, err := file.Open()
        defer file.Close()

        if err != nil {
            error_tmpl.Execute(rw, &Context{Message: template.URL("Failed to parse one of files")})
            return
        }

        tiles = append(tiles, imageFromReader(file))
    }

    fmt.Println("Fetching base file from user")
    file, _, err := req.FormFile("file")
    defer file.Close()

    if err != nil {
        error_tmpl.Execute(rw, &Context{Message: template.URL("Failed to read from a user")})
        return
    }
    img := imageFromReader(file)

    dimension, _ := strconv.Atoi(req.Form["dimension"][0])

    fmt.Println("Generating mosaic")
    resImage := generateMosaic(img, tiles, dimension, dimension)
    fmt.Println("Displaying results")
    displayImgObjAndOrig(rw, resImage, img)
}

func main() {
    router := web.New(Context{}).
        Get("/", (*Context).ShowPicture).
        Post("/", (*Context).RefreshPicture).
        Post("/tile", (*Context).TilePicture)

    http.ListenAndServe("localhost:3000", router)
}
