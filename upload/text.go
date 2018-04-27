package upload

import (
    "sync"
    "os"
    "log"
    "io/ioutil"
    "net/http"
    "net/url"
    "bytes"
)

func Text(fileName string, results chan string, n *sync.WaitGroup) {
    defer n.Done()

    file, err := os.Open(fileName)

    if err != nil {
        log.Println(err)
        return
    }

    defer file.Close()

    data, err := ioutil.ReadAll(file)

    if err != nil {
        log.Println(err)
        return
    }

    resp, err := http.PostForm(farseerfcUrl, url.Values{"c": []string{string(data)}})

    if err != nil {
        log.Println(err)
        return
    }

    defer resp.Body.Close()

    stringBuffer := bytes.NewBufferString("")

    stringBuffer.ReadFrom(resp.Body)

    result := []rune(stringBuffer.String())
    result = result[:len(result)-1]

    results <- fileName + ": " + string(result) + "/\n"
}
