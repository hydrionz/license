package api

import (
	"github.com/gin-gonic/gin"
	"license/mobaxterm/service"
	v1 "license/v1"
	"strconv"
)

// Controller provides routing for MobaXterm-related actions.
type Controller struct {
}

// NewMobaXtermController creates a new instance of MobaXtermController with dependencies injected.
func NewMobaXtermController() *Controller {
	return &Controller{}
}

// FetchVersions retrieves the available MobaXterm versions from the official website.
func (ctrl *Controller) FetchVersions(c *gin.Context) {
	versions, err := service.FetchVersions()
	if err != nil {
		v1.HandleError(c, 500, "Failed to fetch versions")
		return
	}

	// Simplified response structure - no nesting to avoid confusion
	v1.HandleSuccess(c, versions)
}

// GenerateLicense generate handles the license generation process.
func (ctrl *Controller) GenerateLicense(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = c.PostForm("name")
	}
	version := c.Query("version")
	if version == "" {
		version = c.PostForm("version")
	}
	countStr := c.Query("count")
	if countStr == "" {
		countStr = c.PostForm("count")
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		v1.HandleError(c, 400, "Invalid count value")
		return
	}
	service.GenerateLicense(count, name, version, c)
}
