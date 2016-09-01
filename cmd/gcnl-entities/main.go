package main

import (
	"fmt"
	"os"

	"github.com/jlubawy/go-gcnl/entities"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "must provide URL")
		os.Exit(1)
	}

	key := os.Getenv("GOOGLE_API_KEY")
	if len(key) == 0 {
		fmt.Fprintln(os.Stderr, "must set GOOGLE_API_KEY environment variable")
		os.Exit(1)
	}

	es, err := entities.NewRequest(key).FromURL(os.Args[1])
	if err != nil {
		panic(err)
	}

	for _, e := range es {
		fmt.Println(e)
	}
}
