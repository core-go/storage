package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const contentTypeHeader = "Content-Type"

type FileHandler struct {
	Service StorageService
	KeyFile string
	Error   func(context.Context, string)
}

func NewFileHandler(service StorageService, keyFile string, options...func(context.Context, string)) *FileHandler {
	var logError func(context.Context, string)
	if len(options) > 0 && options[0] != nil {
		logError = options[0]
	}
	return &FileHandler{Service: service, KeyFile: keyFile, Error: logError}
}

func (s FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	i := strings.LastIndex(r.RequestURI, "/")
	filename := ""
	if i <= 0 {
		http.Error(w, "require id", http.StatusBadRequest)
		return
	}
	filename = r.RequestURI[i+1:]
	rs, err := s.Service.Delete(r.Context(), filename)
	if err != nil {
		log(s.Error, r, err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, r, http.StatusOK, rs)
}

func (s FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// 32 << 20 is the maxMemory param for ParseMultipartForm, equals to 32MB (1MB = 1024 * 1024 bytes = 2^20 bytes)
	// After you call ParseMultipartForm, the file will be saved in the server memory with maxMemory size.
	// If the file size is larger than maxMemory, the rest of the data will be saved in a system temporary file.
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log(s.Error, r, err.Error())
		http.Error(w, "not available", http.StatusInternalServerError)
		return
	}

	file, handler, err0 := r.FormFile(s.KeyFile)
	if err0 != nil {
		log(s.Error, r, "Cannot get key "+err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	bufferFile := bytes.NewBuffer(nil)
	if _, err1 := io.Copy(bufferFile, file); err != nil {
		log(s.Error, r, "error one post request "+err1.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	content := File{Name: handler.Filename, ContentType: handler.Header.Get(contentTypeHeader), Bytes: bufferFile.Bytes()}

	rs, err2 := s.Service.Upload(r.Context(), content)
	if err2 != nil {
		s.Error(r.Context(), err2.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, r, http.StatusOK, rs)
}

func respond(w http.ResponseWriter, r *http.Request, code int, result interface{}) {
	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func log(logError func(context.Context, string), r *http.Request, err string) {
	if logError != nil {
		logError(r.Context(), err)
	}
}
