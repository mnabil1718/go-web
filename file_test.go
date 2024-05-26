package goweb

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func FileDownload(writer http.ResponseWriter, request *http.Request) {

	filename := request.URL.Query().Get("filename")
	if filename == "" {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, "Bad Request. Filename is required")
	}

	writer.Header().Add(`Content-Disposition`, `attachment; filename="`+filename+`"`)
	http.ServeFile(writer, request, "./resources/upload/"+filename)
}

func DisplayUploadForm(writer http.ResponseWriter, request *http.Request) {

	// parsedTemplates embed FS already declared in this block
	err := parsedTemplates.ExecuteTemplate(writer, "upload-form.gohtml", nil)
	if err != nil {
		panic(err)
	}

}

func FileUpload(writer http.ResponseWriter, request *http.Request) {

	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		panic(err)
	}

	var uploadPath string = "./resources/upload/"
	// automatically create necessary parent directories.
	// if already did, then it does nothing
	err = os.MkdirAll(uploadPath, 0755)
	if err != nil {
		panic(err)
	}
	// last dir path has to exists first. Make sure ./resources/upload exists
	// MKdirAll above ensures its parent dirs exists
	fileDestination, err := os.Create(uploadPath + fileHeader.Filename)

	if err != nil {
		panic(err)
	}

	_, err = io.Copy(fileDestination, file)

	if err != nil {
		panic(err)
	}

	name := request.PostFormValue("name")
	parsedTemplates.ExecuteTemplate(writer, "upload-success.gohtml", map[string]any{
		"Name": name,
		"File": map[string]string{
			"Filename": fileHeader.Filename,
			"FileURL":  "/static/" + fileHeader.Filename,
		},
	})
}

func TestFileRequestServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/form", DisplayUploadForm)
	mux.HandleFunc("/upload", FileUpload)
	mux.HandleFunc("/download", FileDownload)
	// without StripPrefix, the url would be: /resources/upload/static/filename
	// which doesnt exists. Thats why we do StripPrefix
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./resources/upload"))))

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}

//go:embed resources/img/img.jpg
var sampleFileContentAsSliceOfBytes []byte

func TestUpload(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("name", "Nabil")
	file, _ := writer.CreateFormFile("file", "test.jpg") // file name here act as file placeholder
	file.Write(sampleFileContentAsSliceOfBytes)          // writing content from embed to file placeholder as []byte
	writer.Close()                                       //  cannot use defer dont know why

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/form", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()

	FileUpload(recorder, request)

	bodyResponse, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(bodyResponse))
}
