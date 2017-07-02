package routes

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

type HealthResponse struct {
	Message string `json:"message"`
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Message: "up and running!"})
}
