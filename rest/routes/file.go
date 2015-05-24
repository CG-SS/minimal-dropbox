package routes

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mininal-dropbox/storage"
	"net/http"
)

type uploadFilesResponse struct {
	NumUploadedFiles int `json:"num_uploaded_files"`
}

func UploadFiles(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("get form err: %w", err))
			return
		}

		files := form.File["files"]
		numUploadedFiles := 0
		for _, file := range files {
			fileHandle, err := file.Open()
			if err != nil {
				log.Printf("failed to open file header: %v", err)
				continue
			}

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, fileHandle); err != nil {
				log.Printf("failed copying file to buffer: %v", err)
				continue
			}

			err = store.StoreFile(file.Filename, buf.Bytes())
			if err != nil {
				log.Printf("failed saving file: %v", err)
				continue
			}

			numUploadedFiles++
		}

		c.JSON(http.StatusOK, uploadFilesResponse{NumUploadedFiles: numUploadedFiles})
	}
}

func GetFile(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")

		fileBytes, err := store.LoadFile(filename)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("get file error: %w", err))
			return
		}

		c.Data(http.StatusOK, "application/octet-stream", fileBytes)
	}
}

type listFilesResponse struct {
	Filenames []string `json:"filenames"`
}

func ListFiles(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		filenames, err := store.ListFiles()
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("list files error: %w", err))
			return
		}

		c.JSON(http.StatusOK, listFilesResponse{Filenames: filenames})
	}
}
