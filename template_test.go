package goweb

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type Body struct {
	Heading   string
	Paragraph string
}

type Address struct {
	Street string
	City   string
}

type NavItem struct {
	Label string
	Link  string
}

type Content struct {
	Title     string
	Paragraph string
}

type Page struct {
	Name       string
	Title      string
	Header     string
	Navigation []NavItem
	Content    Content
}

//go:embed templates/*.gohtml
var templatesFiles embed.FS

// template.New("") will be used as base template for all template files
var parsedTemplates = template.Must(template.New("").Funcs(map[string]any{
	"upper": func(s string) string {
		return strings.ToUpper(s)
	},
}).ParseFS(templatesFiles, "templates/*.gohtml"))

// NOTICE: the page Page type cannot be pointer otherwise it wont work in the template. Bloody hell
// I spent too much time looking to this bug. The error message is not helpful at all.
func (page Page) SayWelcome() string {
	return "Welcome, " + page.Name + "!"
}

func TemplateHandler(writer http.ResponseWriter, request *http.Request) {
	// data could be passed as nested struct or nested maps. I choose to combine both
	parsedTemplates.ExecuteTemplate(writer, "name.gohtml", map[string]any{
		"Title": "Example Page",
		"Body": Body{
			Heading:   "Example Heading",
			Paragraph: "Example paragraph not too long tho",
		},
	})
}

func TemplateActionHandler(writer http.ResponseWriter, request *http.Request) {
	var values url.Values = request.URL.Query()
	var address Address            // nil struct, can be checks in template
	name := values.Get("name")     // single value
	street := values.Get("street") //single value
	city := values.Get("city")     // single value
	friends := values["friend"]    // array of friends

	if street != "" && city != "" {
		address = Address{
			Street: street,
			City:   city,
		}
	}
	parsedTemplates.ExecuteTemplate(writer, "template-action.gohtml", map[string]any{
		"Title":   "Page Title",
		"Name":    name,
		"Address": address,
		"Friends": friends,
	})
}

func TemplateLayoutHandler(writer http.ResponseWriter, request *http.Request) {

	parsedTemplates.ExecuteTemplate(writer, "layouts", Page{
		Name:   "Nabil",
		Title:  "Page Title",
		Header: "Page Header",
		Navigation: []NavItem{
			{
				Label: "Home",
				Link:  "/",
			},
			{
				Label: "About",
				Link:  "/about",
			},
			{
				Label: "Projects",
				Link:  "/projects",
			},
		},
		Content: Content{
			Title:     "Content Title",
			Paragraph: "Example Paragraph. Not too long tho.",
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

func TestActionTemplate(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080?name=Nabil&city=Bogor&street=Jl.%20Pengangsaan%20Timur&friend=Tina&friend=Joni", nil)
	recorder := httptest.NewRecorder()

	TemplateActionHandler(recorder, request)

	body, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(body))
}

func TestLayoutTemplate(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	recorder := httptest.NewRecorder()

	TemplateLayoutHandler(recorder, request)

	body, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(body))
}
