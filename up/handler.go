package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const contentTypeHeader = "Content-Type"

type UploadTransport interface {
	GetFile(w http.ResponseWriter, r *http.Request)
	DeleteFile(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service UploadManager, logError func(context.Context, string, ...map[string]interface{}),
	keyFile string, generate func(ctx context.Context) (string, error), allowedExtensions string, opts ...int,
) *Handler {
	idIndex := 1
	return &Handler{Service: service, LogError: logError,
		KeyFile: keyFile, generateId: generate, AllowedExtensions: allowedExtensions, idIndex: idIndex,
	}
}

type Handler struct {
	Service           UploadManager
	LogError          func(context.Context, string, ...map[string]interface{})
	KeyFile           string
	AllowedExtensions string
	generateId        func(ctx context.Context) (string, error)
	idIndex           int
}

func (u *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredParam(w, r, u.idIndex)
	if len(id) > 0 {
		err := r.ParseMultipartForm(200000)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		fd := r.MultipartForm

		files := fd.File[u.KeyFile]
		_, handler, _ := r.FormFile(u.KeyFile)
		contentType := handler.Header.Get(contentTypeHeader)
		if len(contentType) == 0 {
			contentType = getExt(handler.Filename)
		}
		generateStr, err := u.generateId(r.Context())
		if err != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		var list []Request
		for i, _ := range files { // loop through the files one by one
			file, err := files[i].Open()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			defer file.Close()
			out := bytes.NewBuffer(nil)

			_, err = io.Copy(out, file)

			if err != nil {
				if u.LogError != nil {
					u.LogError(r.Context(), err.Error())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			bytes := out.Bytes()
			name := generateStr + "_" + files[i].Filename
			if valid, err := u.validateFileName(files[i].Filename, u.AllowedExtensions); !valid {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			list = append(list, Request{files[i].Filename, name, contentType, bytes})
		}
		if len(list) == 0 {
			http.Error(w, "require input file", http.StatusBadRequest)
		}
		rs, err := u.Service.Upload(r.Context(), id, list[0])
		if err != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		respond(w, http.StatusOK, Response{Attachments: *rs})
	}
}

func (u *Handler) validateFileName(filename string, allowedExtensions string) (bool, error) {
	if len(allowedExtensions) == 0 {
		return true, nil
	}
	// Create the regular expression object
	regex := regexp.MustCompile(allowedExtensions)

	// Check if the filename matches the regular expression pattern
	if regex.MatchString(filename) {
		return false, errors.New("valid file type")
	} else {
		return false, errors.New("invalid file type")
	}
}

func (u *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	id := GetRequiredParam(w, r, u.idIndex)
	if len(id) > 0 {
		result, err4 := u.Service.Delete(r.Context(), id, url)
		if err4 != nil {
			if u.LogError != nil {
				u.LogError(r.Context(), err4.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				http.Error(w, err4.Error(), http.StatusInternalServerError)
			}
			return
		}
		respond(w, http.StatusOK, result)
	}
}

func respond(w http.ResponseWriter, code int, result interface{}) {
	res, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
func GetRequiredParam(w http.ResponseWriter, r *http.Request, options ...int) string {
	p := GetParam(r, options...)
	if len(p) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return ""
	}
	return p
}

func GetParam(r *http.Request, options ...int) string {
	offset := 0
	if len(options) > 0 && options[0] > 0 {
		offset = options[0]
	}
	s := r.URL.Path
	params := strings.Split(s, "/")
	i := len(params) - 1 - offset
	if i >= 0 {
		return params[i]
	} else {
		return ""
	}
}
