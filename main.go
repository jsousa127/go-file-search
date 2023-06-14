package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./file-search <PATH> <SEARCH_STRING>")
		os.Exit(0)
	}

	path := os.Args[1]
	searchString := os.Args[2]

	ctx := context.Background()

	results, err := searchDir(ctx, path, searchString)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if len(results) == 0 {
		fmt.Println("No files found")
		os.Exit(0)
	}

	for _, path := range results {
		fmt.Println(path)
	}
}
