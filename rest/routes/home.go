package routes

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"mininal-dropbox/storage"
)

//go:embed static/index.html
var indexPage []byte

type pageData struct {
	Files []string
}

func Home(store storage.Storage) (gin.HandlerFunc, error) {
	tmpl, err := template.New("index").Parse(string(indexPage))
	if err != nil {
		return nil, fmt.Errorf("failed parsing template: %w", err)
	}

	return func(c *gin.Context) {
		files, err := store.ListFiles()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("list files error: %w", err))
			return
		}

		var b bytes.Buffer
		err = tmpl.Execute(&b, pageData{Files: files})
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("template execution error: %w", err))
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", b.Bytes())
	}, nil
}
