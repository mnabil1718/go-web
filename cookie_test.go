package goweb

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetCookie(writer http.ResponseWriter, request *http.Request) {

	cookie := new(http.Cookie)
	cookie.Name = "X-Powered-By"
	cookie.Value = request.URL.Query().Get("name")
	cookie.Path = "/"

	http.SetCookie(writer, cookie)
	fmt.Fprintf(writer, "Cookie set successfully")

}

func GetCookies(writer http.ResponseWriter, request *http.Request) {

	cookies := request.Cookies()

	if len(cookies) == 0 {
		fmt.Fprintf(writer, "Cookie not found")
	} else {
		for _, cookie := range cookies {
			fmt.Fprintf(writer, "%s:%s\n", cookie.Name, cookie.Value)
		}
	}
}

func TestCookie(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080?name=Tingkatin", nil)
	recorder := httptest.NewRecorder()

	request2 := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	request2.AddCookie(&http.Cookie{Name: "X-Powered-By", Value: "Tingkatin"})
	request2.AddCookie(&http.Cookie{Name: "Email", Value: "tingkatin@gmail.com"})
	recorder2 := httptest.NewRecorder()

	SetCookie(recorder, request)
	GetCookies(recorder2, request2)

	// setcookie
	for _, cookie := range recorder.Result().Cookies() {
		fmt.Printf("%s:%s\n", cookie.Name, cookie.Value)
	}

	// getcookie
	response2 := recorder2.Result()
	body, _ := io.ReadAll(response2.Body)
	fmt.Println(string(body))

}
