package api

import (
	"github.com/gin-gonic/gin"
	"license/finalshell/service"
	v1 "license/v1"
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
		v1.HandleError(c, 400, "Machine code cannot be empty")
		return
	}
	licenses := service.GenerateLicense(machineCode)
	v1.HandleSuccess(c, licenses)
}
