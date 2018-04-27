package upload

import (
    "sync"
    "os"
    "log"
    "bytes"
    "mime/multipart"
    "io"
    "net/http"
    "mime"
    "path"
    "strings"
)

func Unity(fileName string, results chan string, n *sync.WaitGroup) {
    defer n.Done()

    file, err := os.Open(fileName)

    if err != nil {
        log.Println(err)
        return
    }
    defer file.Close()

    buf := new(bytes.Buffer)

    writer := multipart.NewWriter(buf)

    part, err := writer.CreateFormFile("c", "c")

    if err != nil {
        log.Println(err)
        return
    }
    _, err = io.Copy(part, file)

    if err != nil {
        log.Println(err)
        return
    }
    writer.Close()

    resp, err := http.Post(farseerfcUrl, writer.FormDataContentType(), buf)

    if err != nil {
        log.Println(err)
        return
    }

    stringBuffer := bytes.NewBufferString("")

    stringBuffer.ReadFrom(resp.Body)

    fileMime := mime.TypeByExtension(path.Ext(fileName))

    switch {
    case strings.HasPrefix(fileMime, "audio"), strings.HasPrefix(fileMime, "video"), strings.HasPrefix(fileMime, "image"):
        results <- fileName + ": " + stringBuffer.String()

    default:
        results <- fileName + ": " + strings.Replace(stringBuffer.String(), "\n", "/\n", 1)
    }
}
