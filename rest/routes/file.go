package routes

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"mininal-dropbox/storage"
	"net/http"
)

type UploadFilesResponse struct {
	NumUploadedFiles int `json:"num_uploaded_files"`
}

func UploadFiles(store storage.Storage, logging zerolog.Logger) gin.HandlerFunc {
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
				logging.Warn().Err(err).Msg("failed to open file header")
				continue
			}

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, fileHandle); err != nil {
				logging.Warn().Err(err).Msg("failed copying file to buffer")
				continue
			}

			err = store.StoreFile(file.Filename, buf.Bytes())
			if err != nil {
				logging.Warn().Err(err).Msg("failed saving file")
				continue
			}

			numUploadedFiles++
		}

		c.JSON(http.StatusOK, UploadFilesResponse{NumUploadedFiles: numUploadedFiles})
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

type ListFilesResponse struct {
	Filenames []string `json:"filenames"`
}

func ListFiles(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		filenames, err := store.ListFiles()
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("list files error: %w", err))
			return
		}

		c.JSON(http.StatusOK, ListFilesResponse{Filenames: filenames})
	}
}

func DeleteFile(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")

		err := store.DeleteFile(filename)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("delete file error: %w", err))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
