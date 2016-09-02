# go-gcnl

This package can be used to easily access the [Google Cloud Natural Language API](https://cloud.google.com/natural-language/).

## Quick-Start

### Analyze Entities Method

    package main

    import (
        "fmt"
        "log"

        "github.com/jlubawy/go-gcnl/entities"
    )

    var apiKey = "<your API key goes here>"
    var content = "Plain text content to analyze"

    func main() {
        req := entities.NewRequest(apiKey)
        entityMap, err := req.FromPlainText(content)
        if err != nil {
            log.Fatalln(err)
        }

        for _, es := range entityMap {
            for _, e := range es {
                fmt.Println(e)
            }
        }
    }

For a more in-depth example see [https://github.com/jlubawy/go-gcnl/tree/master/cmd/gcnl-entities](https://github.com/jlubawy/go-gcnl/tree/master/cmd/gcnl-entities).

To analyze an HTML document given a URL see [entities.FromURL](https://github.com/jlubawy/go-gcnl/blob/master/entities/entities.go#L76).

## TODO

- [x] analyzeEntities
- [ ] analyzeSentiment
- [ ] annotateText
