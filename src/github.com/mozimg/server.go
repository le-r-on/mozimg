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
)

const MAX_FILE_SIZE = 1024 * 1024
const MAX_NUMBER_OF_FILES = 100

type Context struct {
    Image template.URL
    Message template.URL
    AvgColor template.URL
}

func (c *Context) ShowPicture(rw web.ResponseWriter, req *web.Request) {
    file, _ := os.Open("test.png")
    defer file.Close()

    imgObj, _, err := image.Decode(file)
    if err != nil {
        error_tmpl.Execute(rw, &Context{Message: template.URL("Failed to read file:" + err.Error())})
        return
    }
    displayImgObjAndAvg(rw, imgObj)
}

func (c *Context) RefreshPicture(rw web.ResponseWriter, req *web.Request) {
    imgs := randomThumbnails(100)
    fmt.Println(imgs)
    resImage := generateMosaic(imgs[0], imgs[1:], 30, 30)
    displayImgObjAndOrig(rw, resImage, imgs[0])
}



func (c *Context) UploadPicture(rw web.ResponseWriter, req *web.Request) {
    // the FormFile function takes in the POST input id file
    file, _, err := req.FormFile("file")
    defer file.Close()

    if err != nil {
        error_tmpl.Execute(rw, &Context{Message: template.URL("Failed to read from a user")})
        return
    }

    img := imageFromReader(file)
    tiles := randomThumbnails(50)
    resImage := generateMosaic(img, tiles, 30, 30)
    displayImgObjAndOrig(rw, resImage, img)

}

func (c *Context) UploadPictures(rw web.ResponseWriter, req *web.Request) {
    // the FormFile function takes in the POST input id file
    err := req.ParseMultipartForm(int64(MAX_FILE_SIZE * MAX_NUMBER_OF_FILES))
    if err != nil {
        error_tmpl.Execute(rw, &Context{Message: template.URL("Failed to parse multipart")})
        return
    }

    files := req.MultipartForm.File["files"]
    for _, file := range files {
        file, err := file.Open()
        defer file.Close()

        if err != nil {
            error_tmpl.Execute(rw, &Context{Message: template.URL("Failed to parse one of files")})
            return
        }

        img := imageFromReader(file)
        displayImgObjAndAvg(rw, img)
    }
}

func main() {
    router := web.New(Context{}).
        Get("/", (*Context).ShowPicture).
        Post("/", (*Context).RefreshPicture).
        Post("/upload", (*Context).UploadPicture).
        Post("/upload_dir", (*Context).UploadPictures)

    http.ListenAndServe("localhost:3000", router)
}
