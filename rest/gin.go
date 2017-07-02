package rest

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"mininal-dropbox/rest/routes"
	"mininal-dropbox/storage"
	"net/http"
	"time"
)

type ginServer struct {
	logging    zerolog.Logger
	httpServer http.Server
	errChan    chan error
}

const homeRoutePath = "/"
const healthRoutePath = "/health"
const fileAllPath = "/file/all"
const getFilePath = "/file/:filename"
const uploadFilesPath = "/file/upload"

func newGinServer(cfg Config, store storage.Storage, logging zerolog.Logger) (Server, error) {
	router, err := createRouter(cfg, store, logging)
	if err != nil {
		return nil, fmt.Errorf("failed creating gin server: %w", err)
	}

	return &ginServer{
		httpServer: http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: router,
		},
		errChan: make(chan error),
	}, nil
}

func createRouter(cfg Config, store storage.Storage, logging zerolog.Logger) (*gin.Engine, error) {
	router := gin.New()

	router.Use(zerologGinMiddleware(logging))
	router.Use(errorHandler(logging))
	if cfg.Cors.Enabled {
		router.Use(cors.New(cors.Config{
			AllowOrigins: cfg.Cors.AllowOrigins,
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		}))
	}

	homeRoute, err := routes.Home(store)
	if err != nil {
		return nil, fmt.Errorf("failed creating home route: %w", err)
	}

	if cfg.HomeRouteEnabled {
		router.GET(homeRoutePath, homeRoute)
	}
	router.GET(healthRoutePath, routes.Health)

	router.GET(fileAllPath, routes.ListFiles(store))
	router.GET(getFilePath, routes.GetFile(store))
	router.DELETE(getFilePath, routes.DeleteFile(store))
	router.POST(uploadFilesPath, routes.UploadFiles(store, logging))

	return router, nil
}

func zerologGinMiddleware(logging zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		timestamp := time.Now()
		latency := timestamp.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		logging.
			Info().
			Dur("latency", latency).
			Str("clientIP", clientIP).
			Str("method", method).
			Int("statusCode", statusCode).
			Str("errorMessage", errorMessage).
			Int("bodySize", bodySize).
			Str("path", path).
			Msg("")
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

func errorHandler(logging zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				logging.Error().Err(ginErr).Send()
			}

			// status -1 doesn't overwrite existing status code
			c.JSON(-1, errorResponse{Message: "failed to handle request"})
		}
	}
}

func (s *ginServer) Start() {
	defer close(s.errChan)

	err := s.httpServer.ListenAndServe()

	if err != nil {
		s.errChan <- err
	}
}

func (s *ginServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *ginServer) ErrChan() <-chan error {
	return s.errChan
}
