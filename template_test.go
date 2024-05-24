package goweb

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Body struct {
	Heading   string
	Paragraph string
}

func TemplateHandler(writer http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("./templates/name.gohtml"))

	// data could be passed as nested struct or nested maps. I choose to combine both
	template.ExecuteTemplate(writer, "name.gohtml", map[string]any{
		"Title": "Example Page",
		"Body": Body{
			Heading:   "Example Heading",
			Paragraph: "Example paragraph not too long tho",
		},
	})
}

func TestTemplate(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	recorder := httptest.NewRecorder()

	TemplateHandler(recorder, request)

	body, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(body))
}
