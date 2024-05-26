package goweb

import (
	"fmt"
	"net/http"
	"testing"
)

type LogMiddleware struct {
	Handler http.Handler
}

func (middleware *LogMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Incoming Request:", request.Method, request.URL.Path)
	middleware.Handler.ServeHTTP(writer, request)
	fmt.Println("Outgoing Response:", writer.Header().Get("Content-Type"))
}

type ErrorHandlingMiddleware struct {
	Handler http.Handler
}

func (middleware *ErrorHandlingMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, "error: %s", err)
		}
	}()
	middleware.Handler.ServeHTTP(writer, request)
}

func TestMiddleware(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello root")
	})

	mux.HandleFunc("/panic", func(writer http.ResponseWriter, request *http.Request) {
		panic("Intentional Internal Server Error")
	})

	errorMiddleware := &ErrorHandlingMiddleware{Handler: &LogMiddleware{Handler: mux}}

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: errorMiddleware,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
