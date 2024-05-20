package goweb

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWebServer(t *testing.T) {

	server := http.Server{
		Addr: "localhost:8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func TestServerHandler(t *testing.T) {

	var handler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintf(writer, request.RequestURI)
		if err != nil {
			panic(err)
		}
	}

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func TestServeMux(t *testing.T) {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you are at %s", r.URL.Path)
	})

	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you are at projects page")
	})

	mux.HandleFunc("/about/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you are at about page. But this path is not prioritized.")
	})

	mux.HandleFunc("/about/us", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you are inside a nested about us path. This path is prioritized.")
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
