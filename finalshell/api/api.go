package api

import (
	"github.com/gin-gonic/gin"
	"license/finalshell/service"
	"net/http"
	"strings"
)

// Controller defines the controller structure
type Controller struct {
}

// NewController creates a new controller instance
func NewController() *Controller {
	return &Controller{}
}

// GenerateLicense generates a license handling function
func (controller *Controller) GenerateLicense(c *gin.Context) {
	machineCode := c.PostForm("machineCode")
	if strings.TrimSpace(machineCode) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Machine code cannot be empty"})
		return
	}
	licenses := service.GenerateLicense(machineCode)
	c.JSON(http.StatusOK, licenses)
}
