package main

import (
	"context"
	"os"
	"reflect"
	"sort"
	"testing"
)

func check(err error, t *testing.T) {
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestSearchFile(t *testing.T) {
	err := os.Mkdir("testing", os.ModePerm)
	check(err, t)
	defer os.Remove("testing")

	err = os.WriteFile("testing/searchTestFile", []byte("testing\nsearch\n"), 0644)
	check(err, t)
	defer os.Remove("testing/searchTestFile")

	t.Run("search file not containing the keyword", func(t *testing.T) {
		got, _ := searchFile("testing/searchTestFile", "not present")
		want := false

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("search file containing the keyword", func(t *testing.T) {
		got, _ := searchFile("testing/searchTestFile", "testing")
		want := true

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("search file containing the keyword not on 1st line", func(t *testing.T) {
		got, _ := searchFile("testing/searchTestFile", "search")
		want := true

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestSearchDir(t *testing.T) {
	err := os.Mkdir("testing", os.ModePerm)
	check(err, t)
	defer os.Remove("testing")

	err = os.Mkdir("testing/nested", os.ModePerm)
	check(err, t)
	defer os.Remove("testing/nested")

	err = os.WriteFile("testing/searchTestFile1", []byte("testing\nsearch\n"), 0644)
	check(err, t)
	err = os.WriteFile("testing/searchTestFile2", []byte("search\ndir\n"), 0644)
	check(err, t)
	err = os.WriteFile("testing/searchTestFile3", []byte("hello\nworld\n"), 0644)
	check(err, t)
	err = os.WriteFile("testing/nested/searchTestFile", []byte("world\n"), 0644)
	check(err, t)

	defer os.Remove("testing/searchTestFile1")
	defer os.Remove("testing/searchTestFile2")
	defer os.Remove("testing/searchTestFile3")
	defer os.Remove("testing/nested/searchTestFile")

	t.Run("search directory for non present keyword", func(t *testing.T) {
		ctx := context.Background()
		got, _ := searchDir(ctx, "testing", "not present")

		if len(got) > 0 {
			t.Errorf("got %v, want []", got)
		}
	})

	t.Run("search directory for present keyword", func(t *testing.T) {
		ctx := context.Background()
		got, _ := searchDir(ctx, "testing", "search")
		sort.Strings(got)
		want := []string{"testing/searchTestFile1", "testing/searchTestFile2"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("search for present keyword in nested directory", func(t *testing.T) {
		ctx := context.Background()
		got, _ := searchDir(ctx, "testing", "world")
		sort.Strings(got)
		want := []string{"testing/nested/searchTestFile", "testing/searchTestFile3"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
