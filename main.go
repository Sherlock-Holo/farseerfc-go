package main

import (
    "bytes"
    "fmt"
    "sync"
    "strings"
    "github.com/atotto/clipboard"
    "farseerfc/upload"
    "os"
    "time"
    "log"
    "bufio"
)

func main() {
    if len(os.Args[1:]) == 0 {
        reader := bufio.NewReader(os.Stdin)

        go func() {
            <-time.After(2 * time.Second)

            if reader.Buffered() == 0 {
                help()
                os.Exit(2)
            }
        }()

        result, err := upload.Pipeline(reader)

        if err != nil {
            log.Fatal(err)
        }

        fmt.Print(result)
        clipboard.WriteAll(strings.Replace(result, "\n", "", 1))

        return
    }

    if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "-help" {
        help()
        os.Exit(2)
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

    fmt.Println(resultString)
}

func help() {
    fmt.Fprintf(os.Stderr, "Usage: %s File1 File2 File3...\n", os.Args[0])
}
