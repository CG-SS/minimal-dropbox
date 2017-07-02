package rest

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"mininal-dropbox/rest/routes"
	"mininal-dropbox/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.ReleaseMode)

	nopLogging := zerolog.Nop()

	var storeCfg storage.Config
	storeCfg.System = storage.Nop
	store, err := storage.NewStorage(storeCfg, nopLogging)
	assert.NoError(t, err)

	var cfg Config
	cfg.HomeRouteEnabled = true
	cfg.Cors.Enabled = false

	router, err := createRouter(cfg, store, nopLogging)
	assert.NoError(t, err)

	return router
}

func TestHealthRoute(t *testing.T) {
	router := setupRouter(t)
	assert.NotNil(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, healthRoutePath, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var healthResponse routes.HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &healthResponse)
	assert.NoError(t, err)
}

func TestHomeRoute(t *testing.T) {
	router := setupRouter(t)
	assert.NotNil(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, homeRoutePath, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, len(w.Body.Bytes()) > 0)
}

func TestGetFileRoute(t *testing.T) {
	router := setupRouter(t)
	assert.NotNil(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, getFilePath, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteFileRoute(t *testing.T) {
	router := setupRouter(t)
	assert.NotNil(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, getFilePath, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestListFilesRoute(t *testing.T) {
	router := setupRouter(t)
	assert.NotNil(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, fileAllPath, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var listFilesResponse routes.ListFilesResponse
	err = json.Unmarshal(w.Body.Bytes(), &listFilesResponse)
	assert.NoError(t, err)
}

func TestUploadFilesRoute(t *testing.T) {
	router := setupRouter(t)
	assert.NotNil(t, router)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, uploadFilesPath, nil)
	assert.NoError(t, err)

	multipartForm := multipart.Form{
		File: make(map[string][]*multipart.FileHeader),
	}
	multipartForm.File["files"] = []*multipart.FileHeader{}
	req.MultipartForm = &multipart.Form{}

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

type fileTest struct {
	name    string
	content string
}

func writeTestFile(t *testing.T, file fileTest, formWriter *multipart.Writer) {
	t.Helper()

	part, err := formWriter.CreateFormFile("files", file.name)
	assert.NoError(t, err)

	_, err = part.Write([]byte(file.content))
	assert.NoError(t, err)
}

func TestUploadFile(t *testing.T) {
	files := []fileTest{
		{
			name:    "test1.txt",
			content: "This is a test file.",
		},
		{
			name:    "test2.txt",
			content: "This is another test file.",
		},
		{
			name:    "test3.txt",
			content: "This is yet another test file.",
		},
	}

	gin.SetMode(gin.ReleaseMode)

	nopLogging := zerolog.Nop()

	var storeCfg storage.Config
	storeCfg.System = storage.Memory
	store, err := storage.NewStorage(storeCfg, nopLogging)
	assert.NoError(t, err)

	var cfg Config
	cfg.HomeRouteEnabled = false
	cfg.Cors.Enabled = false

	router, err := createRouter(cfg, store, nopLogging)
	assert.NoError(t, err)

	pr, pw := io.Pipe()
	formWriter := multipart.NewWriter(pw)

	go func() {
		for _, file := range files {
			writeTestFile(t, file, formWriter)
		}

		err = formWriter.Close()
		assert.NoError(t, err)
	}()

	res := httptest.NewRecorder()
	req, err := http.NewRequest("POST", uploadFilesPath, pr)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", formWriter.FormDataContentType())

	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusOK)

	filenames, err := store.ListFiles()
	assert.NoError(t, err)

	for _, file := range files {
		assert.Contains(t, filenames, file.name)

		fileContent, err := store.LoadFile(file.name)
		assert.NoError(t, err)

		assert.Equal(t, file.content, string(fileContent))
	}
}

func TestClosesGinServerCloses(t *testing.T) {
	var cfg Config
	cfg.Host = "localhost"
	cfg.Port = 12345
	cfg.System = Gin
	cfg.Cors.Enabled = false

	gin.SetMode(gin.ReleaseMode)

	nopLogging := zerolog.Nop()

	var storeCfg storage.Config
	storeCfg.System = storage.Memory
	store, err := storage.NewStorage(storeCfg, nopLogging)
	assert.NoError(t, err)

	server, err := newGinServer(cfg, store, nopLogging)
	assert.NoError(t, err)
	go server.Start()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = server.Stop(ctx)
	assert.NoError(t, err)

	select {
	case serverErr := <-server.ErrChan():
		assert.Equal(t, serverErr, http.ErrServerClosed)
	case <-time.After(1 * time.Second):
		t.FailNow()
	}
}

func TestGinServerNoHomeRouteIfDisabled(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	nopLogging := zerolog.Nop()

	var storeCfg storage.Config
	storeCfg.System = storage.Nop
	store, err := storage.NewStorage(storeCfg, nopLogging)
	assert.NoError(t, err)

	var cfg Config
	cfg.HomeRouteEnabled = false
	cfg.Cors.Enabled = false

	router, err := createRouter(cfg, store, nopLogging)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, homeRoutePath, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusNotFound)
}
