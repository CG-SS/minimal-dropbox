package rest

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"mininal-dropbox/rest/routes"
	"mininal-dropbox/storage"
	"net/http"
	"time"
)

type ginServer struct {
	httpServer http.Server
	errChan    chan error
}

func newGinServer(cfg Config, store storage.Storage) (Server, error) {
	router := gin.New()

	router.Use(logGinMiddleware())
	router.Use(errorHandler())
	router.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Cors.AllowOrigins,
		AllowMethods: []string{"GET", "POST"},
	}))

	homeRoute, err := routes.Home(store)
	if err != nil {
		return nil, fmt.Errorf("failed creating home route: %w", err)
	}

	router.GET("/", homeRoute)
	router.GET("/health", routes.Health)

	router.GET("/file/all", routes.ListFiles(store))
	router.GET("/file/:filename", routes.GetFile(store))
	router.POST("/file/upload", routes.UploadFiles(store))

	return &ginServer{
		httpServer: http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: router,
		},
		errChan: make(chan error),
	}, nil
}

func logGinMiddleware() gin.HandlerFunc {
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

		log.Printf("latency: %v clientIP: %v method: %v statusCode: %v errorMessage: %v bodySize: %v path: %v", latency, clientIP, method, statusCode, errorMessage, bodySize, path)
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				log.Printf("error: %v", ginErr)
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
