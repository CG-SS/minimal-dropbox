package routes

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

type healthResponse struct {
	Message string `json:"message"`
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{Message: "up and running!"})
}
