package upload

import (
    "bytes"
    "io"
    "mime/multipart"
    "log"
    "net/http"
    "bufio"
)

func Pipeline(reader *bufio.Reader) (string, error) {
    mimeBuf := new(bytes.Buffer)

    writer := multipart.NewWriter(mimeBuf)

    part, err := writer.CreateFormFile("c", "c")

    if err != nil {
        log.Println(err)
        return "", err
    }
    _, err = io.Copy(part, reader)

    if err != nil {
        log.Println(err)
        return "", err
    }
    writer.Close()

    resp, err := http.Post(farseerfcUrl, writer.FormDataContentType(), mimeBuf)

    if err != nil {
        log.Println(err)
        return "", err
    }

    stringBuf := bytes.NewBufferString("")

    _, err = stringBuf.ReadFrom(resp.Body)

    if err != nil {
        return "", err
    }

    return stringBuf.String(), nil
}
