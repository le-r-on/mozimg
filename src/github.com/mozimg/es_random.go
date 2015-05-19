package main

import (
    "github.com/polds/imgbase64"
    "github.com/franela/goreq"
    "fmt"
    "regexp"
    "bytes"
    "math"
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

func randomThumbnails(size int) []string {
    // Take twice as much since some thumbnails will be incessible
    urls := getUrls(int(math.Max(float64(size * 2), float64(5))))

    images := make([]string, size)
    count := 0
    for _, url := range urls {
        buffer, err := fetchImage(url)
        if err == nil {
            images[count] = imgbase64.FromBuffer(*bytes.NewBuffer(buffer))
        } else {
            fmt.Println("Failed to fetch image", err)
        }
    }

    return images
}