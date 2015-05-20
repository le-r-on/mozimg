package main

import (
    _ "github.com/polds/imgbase64"
    "github.com/franela/goreq"
    "fmt"
    "regexp"
    "bytes"
    "math"
    "image"
)

func fetchImage(url string) ([]byte, error) {
    res, err := goreq.Request{
        Uri: url,
    }.Do()

    if err != nil {
        return nil, err
    }

    result, _ := res.Body.ToString()

    return []byte(result), nil
}

func getUrls(size int) []string {
    query := fmt.Sprintf("http://natural.elastic.tubularlabs.net:9200/natural/vine_videos/_search?size=%d", size)
    res, _ := goreq.Request{
        Method: "POST",
        Uri: query,
        Body: 
        `{"query" : {"range": {"statistics.views_count": {"from": 10000000}}},
          "sort" : {
            "_script" : { 
                "script" : "Math.random()",
                "type" : "number",
                "params" : {},
                "order" : "asc"
            }
          }
        }`,
    }.Do()

    body, _ := res.Body.ToString()

    r, _ := regexp.Compile("\"0x0\":\"([^\"]*)\"")
    ress := r.FindAllStringSubmatch(body, -1)

    urls := make([]string, size, size);

    for i, value := range ress {
        urls[i] = value[1]
    }

    return urls
}

func randomThumbnails(size int) []*image.Image {
    // Take twice as much since some thumbnails will be incessible
    urls := getUrls(int(math.Max(float64(size * 2), float64(5))))

    images := make([]*image.Image, size, size)
    count := 0
    for _, url := range urls {
        buffer, err := fetchImage(url)
        if err == nil {
            tmp := imageFromReader(bytes.NewReader(buffer))
            if tmp != nil {
                images[count] = &tmp
                count++
            }
        } else {
            fmt.Println("Failed to fetch image", err)
        }
        if count == size {
            break
        }
    }

    return images
}