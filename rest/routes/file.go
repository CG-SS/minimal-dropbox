package routes

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"mininal-dropbox/storage"
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

			err = store.StoreFile(file.Filename, fileHandle)
			if err != nil {
				logging.Warn().Err(err).Msg("failed saving file")
				continue
			}

			err = fileHandle.Close()
			if err != nil {
				logging.Warn().Err(err).Msg("failed closing file")
				continue
			}

			numUploadedFiles++
		}

		c.JSON(http.StatusOK, UploadFilesResponse{NumUploadedFiles: numUploadedFiles})
	}
}

func GetFile(store storage.Storage, bufferSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")

		fileReader, err := store.LoadFile(filename)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("get file error: %w", err))
			return
		}
		defer fileReader.Close()

		c.Stream(func(w io.Writer) bool {
			buffer := make([]byte, bufferSize)

			for {
				_, err := fileReader.Read(buffer)
				if err != nil {
					return false
				}
				_, err = w.Write(buffer)
				if err != nil {
					return false
				}
			}
		})
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
