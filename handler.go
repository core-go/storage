package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const ContentTypeHeader = "Content-Type"

type UploadFileHandler struct {
	CloudService StorageService
	KeyFile      string
	LogError     func(context.Context, string)
}

func (s UploadFileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	i := strings.LastIndex(r.RequestURI, "/")
	filename := ""
	if i <= 0 {
		http.Error(w, "require id", http.StatusBadRequest)
		return
	}
	filename = r.RequestURI[i+1:]
	rs, err := s.CloudService.Delete(r.Context(), filename)
	if err != nil {
		log(s.LogError, r, err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respond(w, r, http.StatusOK, rs)
}

func (s UploadFileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// 32 << 20 is the maxMemory param for ParseMultipartForm, equals to 32MB (1MB = 1024 * 1024 bytes = 2^20 bytes)
	// After you call ParseMultipartForm, the file will be saved in the server memory with maxMemory size.
	// If the file size is larger than maxMemory, the rest of the data will be saved in a system temporary file.
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log(s.LogError, r, err.Error())
		http.Error(w, "is not available", http.StatusInternalServerError)
		return
	}

	file, handler, err0 := r.FormFile(s.KeyFile)
	if err0 != nil {
		log(s.LogError, r, fmt.Sprintf("Can't get key %s\n", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	bufferFile := bytes.NewBuffer(nil)
	if _, err1 := io.Copy(bufferFile, file); err != nil {
		log(s.LogError, r, fmt.Sprintf("Error one post request ", err1.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	content := File{FileName: handler.Filename, ContentType: handler.Header.Get(ContentTypeHeader), BytesData: bufferFile.Bytes()}

	rs, err2 := s.CloudService.Upload(r.Context(), content)
	if err2 != nil {
		s.LogError(r.Context(), err2.Error())
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
