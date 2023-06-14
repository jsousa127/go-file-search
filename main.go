package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	path := os.Args[1]
	searchString := os.Args[2]

	ctx := context.Background()

	results, err := searchDir(ctx, path, searchString)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if len(results) == 0 {
		fmt.Println("No files found")
	}

	for _, path := range results {
		fmt.Println(path)
	}
}
