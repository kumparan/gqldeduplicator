# Graphql Deduplicator
[![godoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat)](https://pkg.go.dev/github.com/kumparan/gqldeduplicator) 
[![go report](https://goreportcard.com/badge/github.com/kumparan/gqldeduplicator)](https://goreportcard.com/report/github.com/kumparan/gqldeduplicator)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/kumparan/gqldeduplicator/blob/master/LICENSE)

GraphQL response deduplicator.

Javascript version: https://github.com/gajus/graphql-deduplicator 

### Usage

```
package main

import (
	"log"

	"github.com/kumparan/gqldeduplicator"
)

func main() {
    data := []byte(`
    {
        "root": [
            {
                "__typename": "foo",
                "id": 1,
                "name": "foo"
            },
            {
                "__typename": "foo",
                "id": 1,
                "name": "foo"
            }
        ]
    }`)

    deflate, err := gqldeduplicator.Deflate(data)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("deflate:", string(deflate.Data))

    inflate, err := gqldeduplicator.Inflate(deflate.Data)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("inflate:", string(inflate.Data))
}
```
