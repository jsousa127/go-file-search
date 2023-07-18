package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

type Response struct {
	Data   []string `json:"data"`
	Error  string   `json:"error"`
	Status int      `json:"status"`
}

func TestServer(t *testing.T) {
	err := os.Mkdir("testing", os.ModePerm)
	check(err, t)
	defer os.Remove("testing")

	err = os.WriteFile("testing/searchTestFile", []byte("testing\nsearch\n"), 0644)
	check(err, t)
	defer os.Remove("testing/searchTestFile")

	engine := Engine()

	t.Run("search for keyword present in file", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/search?path=testing&keyword=testing", nil)
		//check(err, t)

		engine.ServeHTTP(recorder, request)
		if recorder.Code != 200 {
			t.Fatalf("status code: got %d, want %d", recorder.Code, statusOK)
		}

		var response Response
		body := recorder.Body.String()

		err = json.Unmarshal([]byte(body), &response)
		check(err, t)

		got := response.Data
		want := []string{"testing/searchTestFile"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("search for keyword not present in file", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/search?path=testing&keyword=noresults", nil)
		check(err, t)

		engine.ServeHTTP(recorder, request)
		if recorder.Code != 200 {
			t.Fatalf("status code: got %d, want %d", recorder.Code, statusOK)
		}

		var response Response
		body := recorder.Body.String()

		err = json.Unmarshal([]byte(body), &response)
		check(err, t)

		if len(response.Data) > 0 {
			t.Errorf("got %v, want []", response.Data)
		}
	})

	t.Run("malformed request missing path", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/search?keyword=noresults", nil)
		check(err, t)

		engine.ServeHTTP(recorder, request)
		if recorder.Code != 400 {
			t.Fatalf("status code: got %d, want %d", recorder.Code, statusBadRequest)
		}

		var response Response
		body := recorder.Body.String()

		err = json.Unmarshal([]byte(body), &response)
		check(err, t)

		got := response.Error
		want := "Malformed Request"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/search?path=testing2&keyword=noresults", nil)
		check(err, t)

		engine.ServeHTTP(recorder, request)
		if recorder.Code != 500 {
			t.Fatalf("status code: got %d, want %d", recorder.Code, statusInternalServerError)
		}

		var response Response
		body := recorder.Body.String()

		err = json.Unmarshal([]byte(body), &response)
		check(err, t)

		got := response.Error
		want := "lstat testing2: no such file or directory"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
