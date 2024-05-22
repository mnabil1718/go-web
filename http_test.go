package goweb

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SayHelloQueryParam(writer http.ResponseWriter, request *http.Request) {
	var values url.Values = request.URL.Query()
	names := values["name"]
	fmt.Fprintln(writer, strings.Join(names, " "))
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Hello")
}

func HeaderHandler(writer http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	fmt.Println(contentType)
	writer.Header().Add("X-Powered-By", "Tingkatin")
	writer.Header().Add("Content-Type-2", contentType)
	fmt.Fprintln(writer, "OK")
}

func FormPostHandler(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	// request.PostFormValue("first_name") # alternative way to get post form value without parsing
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(writer, "Hello, %s %s!\n", request.PostForm.Get("first_name"), request.PostForm.Get("last_name"))
}

func ResponseCodeHandler(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	if name == "" {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, "Bad request. Name is empty")
	} else {
		fmt.Fprintf(writer, "Hello %s", name)
	}
}

func TestHttp(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost:8080", nil)
	recorder := httptest.NewRecorder()

	Handler(recorder, request)

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

func TestHeader(t *testing.T) {
	request := httptest.NewRequest("GET", "http://localhost:8080?name=Muhammad&name=Eko", nil)
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	HeaderHandler(recorder, request)

	response := recorder.Result()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("header:", response.Header)
	fmt.Println("body:", string(body))
}

func TestFormPost(t *testing.T) {
	requestBody := strings.NewReader("first_name=Muhammad&last_name=Nabil")
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	record := httptest.NewRecorder()

	FormPostHandler(record, request)

	response := record.Result()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func TestResponseCode(t *testing.T) {
	request1 := httptest.NewRequest("GET", "http://localhost:8080", nil)
	request2 := httptest.NewRequest("GET", "http://localhost:8080?name=Nabil", nil)
	recorder1 := httptest.NewRecorder()
	recorder2 := httptest.NewRecorder()

	ResponseCodeHandler(recorder1, request1)
	ResponseCodeHandler(recorder2, request2)

	response1 := recorder1.Result()
	response2 := recorder2.Result()

	assert.Equal(t, http.StatusBadRequest, response1.StatusCode, "Response code should be 400")
	assert.Equal(t, http.StatusOK, response2.StatusCode, "Response code should be 200")

}
