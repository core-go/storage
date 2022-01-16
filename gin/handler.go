package gin

import (
	"bytes"
	"context"
	st "github.com/core-go/storage"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

const contentTypeHeader = "Content-Type"

type FileHandler struct {
	Service   st.StorageService
	Directory string
	KeyFile   string
	Error     func(context.Context, string)
}

func NewFileHandler(service st.StorageService, directory string, keyFile string, options ...func(context.Context, string)) *FileHandler {
	var logError func(context.Context, string)
	if len(options) > 0 && options[0] != nil {
		logError = options[0]
	}
	return &FileHandler{Service: service, Directory: directory, KeyFile: keyFile, Error: logError}
}

func (s FileHandler) DeleteFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		id := ctx.Param("id")
		if len(id) == 0 {
			ctx.String(http.StatusBadRequest, "require id")
			return
		}
		rs, err := s.Service.Delete(r.Context(), id)
		if err != nil {
			log(s.Error, r, err.Error())
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		ctx.JSON(http.StatusOK, rs)
	}
}

func (s FileHandler) UploadFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		// 32 << 20 is the maxMemory param for ParseMultipartForm, equals to 32MB (1MB = 1024 * 1024 bytes = 2^20 bytes)
		// After you call ParseMultipartForm, the file will be saved in the server memory with maxMemory size.
		// If the file size is larger than maxMemory, the rest of the data will be saved in a system temporary file.
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log(s.Error, r, err.Error())
			ctx.String(http.StatusInternalServerError, "not available")
			return
		}

		file, handler, err0 := r.FormFile(s.KeyFile)
		if err0 != nil {
			log(s.Error, r, "Cannot get key "+err.Error())
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		bufferFile := bytes.NewBuffer(nil)
		if _, err1 := io.Copy(bufferFile, file); err != nil {
			log(s.Error, r, "error one post request "+err1.Error())
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer file.Close()
		bytes := bufferFile.Bytes()
		contentTye := handler.Header.Get(contentTypeHeader)
		if len(contentTye) == 0 {
			contentTye = getExt(handler.Filename)
		}
		rs, err2 := s.Service.Upload(r.Context(), s.Directory, handler.Filename, bytes, contentTye)
		if err2 != nil {
			s.Error(r.Context(), err2.Error())
			ctx.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		ctx.JSON(http.StatusOK, rs)
	}
}

func log(logError func(context.Context, string), r *http.Request, err string) {
	if logError != nil {
		logError(r.Context(), err)
	}
}
func getExt(file string) string {
	ext := filepath.Ext(file)
	if strings.HasPrefix(ext, ":") {
		ext = ext[1:]
		return ext
	}
	return ext
}