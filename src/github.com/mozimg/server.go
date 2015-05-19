package main

import (
    "github.com/gocraft/web"
    "github.com/polds/imgbase64"
    _ "github.com/franela/goreq"
    "html/template"
    "fmt"
    "net/http"
    _ "log"
    "bytes"
)

type Context struct {
    Image template.URL
}

func (c *Context) ShowPicture(rw web.ResponseWriter, req *web.Request) {
    img, _ := imgbase64.FromLocal("test.png")
    html.Execute(rw, &Context{Image: template.URL(img)})
}

func (c *Context) RefreshPicture(rw web.ResponseWriter, req *web.Request) {
    img := randomThumbnails(1)[0]
    html.Execute(rw, &Context{Image: template.URL(img)})
}

func (c *Context) UploadPicture(rw web.ResponseWriter, req *web.Request) {
    // the FormFile function takes in the POST input id file
    file, _, err := req.FormFile("file")

    if err != nil {
        //log.Warning("Failed to read from a user", err)
        fmt.Println("Failed to read from a user", err)
        return
    }
    defer file.Close()

    buffer := make([]byte, 1024*1024)
    _, err = file.Read(buffer)

    if err != nil {
        // log.Warning("File is too big.")
        fmt.Println("Failed to read file", err)
    }

    //if n == 1024*1024 { log.Warning("File is too big.")}
    img := imgbase64.FromBuffer(*bytes.NewBuffer(buffer))
    html.Execute(rw, &Context{Image: template.URL(img)})
}


func main() {
    router := web.New(Context{}).
        Get("/", (*Context).ShowPicture).
        Post("/", (*Context).RefreshPicture).
        Post("/upload", (*Context).UploadPicture)
    http.ListenAndServe("localhost:3000", router)
}
