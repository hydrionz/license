package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"license/config"
	finalshell "license/finalshell/api"
	gitlab "license/gitlab/api"
	jetbrainCode "license/jetbrains/code/api"
	jetbrainServer "license/jetbrains/server/api"
	jrebel "license/jrebel/api"
	"license/logger"
	mobaxterm "license/mobaxterm/api"
	rpc "license/rpc/controller"
	"license/utils/useragent"
	"net/http"
	"strings"
	"time"
)

// VersionResponse defines the version response structure
type VersionResponse struct {
	Version       string `json:"version"`
	NeedUpdate    bool   `json:"needUpdate"`
	LatestVersion string `json:"latestVersion,omitempty"`
}

// GitHubRelease GitHub API release response structure
type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// Cache for GitHub latest version information
var (
	cachedLatestVersion string
	lastFetchTime       time.Time
	cacheExpiration     = 30 * time.Minute
)

// List of API path prefixes
var apiPrefixes = []string{
	"/server/",
	"/final-shell/",
	"/gitlab/",
	"/rpc/",
	"/jrebel/",
	"/agent/",
	"/mobaxterm/",
	"/jetbrains/",
}

// IsAPIPath determines if the given path is an API path
func IsAPIPath(path string) bool {
	for _, prefix := range apiPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// HandleAPIRequest handles API requests
func HandleAPIRequest(c *gin.Context) {
	// Create a temporary routing engine to handle the request
	tmpEngine := gin.New()
	tmpGroup := tmpEngine.Group("/")

	// Set up the router
	SetupRouter(tmpGroup)

	// Handle the request
	tmpEngine.HandleContext(c)
}

// getLatestVersionFromGitHub fetches the latest version from GitHub(removing the "v" prefix)
func getLatestVersionFromGitHub() string {
	// If the cache has not expired, use the cached value
	if !lastFetchTime.IsZero() && time.Since(lastFetchTime) < cacheExpiration && cachedLatestVersion != "" {
		return cachedLatestVersion
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/nannanStrawberry314/license/releases/latest", nil)
	if err != nil {
		logger.Error("Failed to create GitHub API request", err)
		return ""
	}

	// Use random User-Agent
	req.Header.Set("User-Agent", useragent.GetRandom())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to request GitHub API", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("GitHub API returned non-200 status code", nil)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read GitHub API response", err)
		return ""
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		logger.Error("Failed to parse GitHub API response", err)
		return ""
	}

	// Remove "v" prefix from version number
	version := release.TagName
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	// Update cache
	cachedLatestVersion = version
	lastFetchTime = time.Now()

	return version
}

// compareVersions compares version numbers
func compareVersions(current, latest string) bool {
	// Split both versions by dot
	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	// Compare each part of the version numbers
	for i := 0; i < len(currentParts) && i < len(latestParts); i++ {
		if currentParts[i] < latestParts[i] {
			return true // Update needed
		} else if currentParts[i] > latestParts[i] {
			return false // No update needed
		}
	}

	// If all previous parts are equal, but latest has more parts, update is needed
	return len(latestParts) > len(currentParts)
}

func SetupRouter(r *gin.RouterGroup) {
	serverGroup := r.Group("/server")
	{
		serverGroup.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": true,
			})
		})

		serverGroup.GET("/version", func(c *gin.Context) {
			currentVersion := config.Version
			latestVersion := getLatestVersionFromGitHub()
			needUpdate := false

			if latestVersion != "" {
				needUpdate = compareVersions(currentVersion, latestVersion)
			}

			c.JSON(200, VersionResponse{
				Version:       currentVersion,
				NeedUpdate:    needUpdate,
				LatestVersion: latestVersion,
			})
		})
	}

	// final-shell
	finalShellApi := finalshell.NewController()
	finalShellGroup := r.Group("/final-shell")
	{
		finalShellGroup.POST("/generateLicense", finalShellApi.GenerateLicense)
	}

	// gitlab
	gitlabApi := gitlab.NewController()
	gitlabGroup := r.Group("/gitlab")
	{
		gitlabGroup.POST("/generate", gitlabApi.Generate)
	}

	// rpc
	rpcApi := rpc.NewRpcController()
	rpcGroup := r.Group("/rpc")
	{
		rpcGroup.GET("/ping.action", rpcApi.Ping)
		rpcGroup.GET("/obtainTicket.action", rpcApi.ObtainTicket)
		rpcGroup.GET("/releaseTicket.action", rpcApi.ReleaseTicket)
	}

	// jrebel
	jrebelLeasesApi := jrebel.NewLeasesController()
	jrebelIndexApi := jrebel.NewIndexController()
	jrebelGroup := r.Group("/jrebel")
	{
		jrebelGroup.GET("/", jrebelIndexApi.IndexHandler)
		jrebelGroup.DELETE("/leases/1", jrebelLeasesApi.Leases1Handler)
		jrebelGroup.POST("/leases/1", func(c *gin.Context) {
			c.Status(405)
		})
		jrebelGroup.POST("/leases", jrebelLeasesApi.LeasesHandler)
		jrebelGroup.POST("/validate-connection", jrebelLeasesApi.ValidateHandler)
		jrebelGroup.POST("/features", jrebelLeasesApi.ValidateHandler)
		jrebelGroup.GET("/features", jrebelLeasesApi.ValidateHandler)
	}
	jrebelAgentGroup := r.Group("/agent")
	{
		jrebelAgentGroup.DELETE("/leases/1", jrebelLeasesApi.Leases1Handler)
		jrebelAgentGroup.POST("/leases/1", func(c *gin.Context) {
			c.Status(405)
		})
		jrebelAgentGroup.POST("/leases", jrebelLeasesApi.LeasesHandler)
		jrebelAgentGroup.POST("/validate-connection", jrebelLeasesApi.ValidateHandler)
		jrebelAgentGroup.POST("/features", jrebelLeasesApi.ValidateHandler)
		jrebelAgentGroup.GET("/features", jrebelLeasesApi.ValidateHandler)
	}

	// mobaxterm
	mobaxtermApi := mobaxterm.NewMobaXtermController()
	mobaxtermGroup := r.Group("/mobaxterm")
	{
		mobaxtermGroup.POST("/generate", mobaxtermApi.GenerateLicense)
		mobaxtermGroup.GET("/versions", func(c *gin.Context) {
			// Add no-cache headers
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")

			mobaxtermApi.FetchVersions(c)
		})
	}

	// jetbrains
	jetbrainsServerApi := jetbrainServer.NewLicenseServerController()
	jetbrainsCodeApi := jetbrainCode.NewController()

	jetbrainsGroup := r.Group("/jetbrains")
	{
		jetbrainsGroup.GET("/licenseServerRule", jetbrainsServerApi.LicenseServerRule)
		jetbrainsGroup.GET("/product/fetchLatest", jetbrainsCodeApi.FetchProduceLatest)
		jetbrainsGroup.GET("/plugin/fetchLatest", jetbrainsCodeApi.FetchPluginLatest)
		jetbrainsGroup.GET("/generate", jetbrainsCodeApi.Generate)
	}
}
