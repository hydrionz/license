package api

import (
	"github.com/gin-gonic/gin"
	"license/mobaxterm/service"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": "Failed to fetch versions",
			"details": err.Error(),
		})
		return
	}
	
	// Simplified response structure - no nesting to avoid confusion
	c.JSON(http.StatusOK, versions)
}

// generate handles the license generation process.
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count value"})
		return
	}
	service.GenerateLicense(count, name, version, c)
}

// DebugHtmlStructure provides debugging information for the HTML structure.
func (ctrl *Controller) DebugHtmlStructure(c *gin.Context) {
	debugInfo, err := service.DebugHtmlStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to debug HTML structure", "details": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, debugInfo)
}
