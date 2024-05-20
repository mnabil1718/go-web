package goweb

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExampleHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello")
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
