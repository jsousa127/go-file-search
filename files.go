package main

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sync/errgroup"
)

func searchFile(path, searchString string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if strings.Contains(line, searchString) {
			return true, nil
		}
	}
	return false, nil
}

func searchDir(ctx context.Context, root, searchString string) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	resultChan := make(chan string, 1)

	checkFile := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() || info.IsDir() {
			return nil
		}

		g.Go(func() error {
			containsSearch, err := searchFile(path, searchString)
			if err != nil {
				return err
			}

			if containsSearch {
				select {
				case resultChan <- path:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})

		return nil
	}

	g.Go(func() error {
		return filepath.Walk(root, checkFile)
	})

	go func() {
		g.Wait()
		close(resultChan)
	}()

	var results []string
	for path := range resultChan {
		results = append(results, path)
	}

	return results, g.Wait()
}