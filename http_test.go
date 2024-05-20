package goweb

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func SayHelloQueryParam(writer http.ResponseWriter, request *http.Request) {
	var values url.Values = request.URL.Query()
	names := values["name"]
	fmt.Fprintln(writer, strings.Join(names, " "))
}

func ExampleHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Hello")
}

func TestHttp(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost:8080", nil)
	recorder := httptest.NewRecorder()

	ExampleHandler(recorder, request)

	response := recorder.Result()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.StatusCode)
	fmt.Println(string(body))
}

func TestQueryParam(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost:8080?name=Muhammad&name=Eko", nil)
	recorder := httptest.NewRecorder()

	SayHelloQueryParam(recorder, request)

	response := recorder.Result()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.StatusCode)
	fmt.Println(string(body))
}
