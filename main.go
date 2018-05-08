package main

import (
    "flag"
    "github.com/Sherlock-Holo/farseerfc-go/upload"
    "log"
    "fmt"
    "github.com/atotto/clipboard"
    "strings"
    "sync"
    "os"
    "bytes"
)

var pipeline = flag.Bool("p", false, "pipeline mode")
var normal = flag.Bool("n", false, "normal mode")

func main() {
    flag.Parse()

    switch {
    case *pipeline:
        result, err := upload.Pipeline()

        if err != nil {
            log.Fatal(err)
        }

        fmt.Print(result)
        clipboard.WriteAll(strings.Replace(result, "\n", "", 1))

        return

    case *normal:
        var n sync.WaitGroup

        files := flag.Args()

        results := make(chan string, len(files))

        for _, file := range files {
            n.Add(1)
            go upload.Unity(file, results, &n)
        }

        go func() {
            n.Wait()
            close(results)
        }()

        var buffer = bytes.NewBufferString("")

        var resultn int

        for result := range results {
            resultn++
            buffer.WriteString(result)
        }

        resultString := buffer.String()
        resultString = string([]rune(resultString)[:len(resultString)-1])

        if resultn == 1 {
            list := strings.Split(resultString, ": ")
            clipboard.WriteAll(list[1])
        } else {
            clipboard.WriteAll(resultString)
        }

        fmt.Println(resultString)

    default:
        flag.Usage()
        os.Exit(2)
    }

    /*if *pipeline {
        result, err := upload.Pipeline()

        if err != nil {
            log.Fatal(err)
        }

        fmt.Print(result)
        clipboard.WriteAll(strings.Replace(result, "\n", "", 1))

        return
    }

    var n sync.WaitGroup

    files := os.Args[1:]

    results := make(chan string, len(files))

    for _, file := range files {
        n.Add(1)
        go upload.Unity(file, results, &n)
    }

    go func() {
        n.Wait()
        close(results)
    }()

    var buffer = bytes.NewBufferString("")

    var resultn int

    for result := range results {
        resultn++
        buffer.WriteString(result)
    }

    resultString := buffer.String()
    resultString = string([]rune(resultString)[:len(resultString)-1])

    if resultn == 1 {
        list := strings.Split(resultString, ": ")
        clipboard.WriteAll(list[1])
    } else {
        clipboard.WriteAll(resultString)
    }

    fmt.Println(resultString)*/
}
