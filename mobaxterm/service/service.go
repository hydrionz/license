package service

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"license/config"
	"license/logger"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"os"
)

// List of common User-Agent strings
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:99.0) Gecko/20100101 Firefox/99.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36 Edg/101.0.1210.39",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 OPR/86.0.4363.59",
}

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// getRandomUserAgent returns a randomly selected user agent from the list
func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// FetchVersions fetches MobaXterm version numbers from the website
func FetchVersions() ([]string, error) {
	// Create HTTP client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Create request
	req, err := http.NewRequest("GET", "https://mobaxterm.mobatek.net/download-home-edition.html", nil)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to create request: %v", err))
		return nil, err
	}

	// Set random User-Agent
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to fetch MobaXterm website: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("", fmt.Errorf("bad status code: %d", resp.StatusCode))
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to parse HTML: %v", err))
		return nil, err
	}

	// Create a set to store unique versions
	versionsMap := make(map[string]bool)

	// Strategy 1: Try to get current version from the download links
	// Use a more specific regex to match "Version X.Y (date)" pattern
	versionRegex := regexp.MustCompile(`Version (\d+\.\d+)`)

	// Strategy 2: Look for version_titre class - this is our primary source
	doc.Find("p.version_titre").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		matches := versionRegex.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				versionsMap[match[1]] = true
			}
		}
	})

	// Convert map to sorted slice
	versions := make([]string, 0, len(versionsMap))
	for version := range versionsMap {
		versions = append(versions, version)
	}

	// If we found no versions, return default list
	if len(versions) == 0 {
		logger.Error("", fmt.Errorf("no versions found, returning default list"))
		return []string{"25.1", "25.0", "24.4", "24.3", "23.6", "23.5", "23.0"}, nil
	}

	// Sort versions in descending order (newer versions first)
	sort.Slice(versions, func(i, j int) bool {
		// Split version into major and minor components
		partsI := strings.Split(versions[i], ".")
		partsJ := strings.Split(versions[j], ".")

		// Parse major version
		majorI, _ := strconv.Atoi(partsI[0])
		majorJ, _ := strconv.Atoi(partsJ[0])

		if majorI != majorJ {
			return majorI > majorJ
		}

		// Parse minor version
		minorI, _ := strconv.Atoi(partsI[1])
		minorJ, _ := strconv.Atoi(partsJ[1])

		return minorI > minorJ
	})

	return versions, nil
}

// DebugHtmlStructure fetches the MobaXterm page and returns debugging information
// about the HTML structure to help diagnose selector issues
func DebugHtmlStructure() (map[string]interface{}, error) {
	// Create HTTP client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Create request
	req, err := http.NewRequest("GET", "https://mobaxterm.mobatek.net/download-home-edition.html", nil)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to create request: %v", err))
		return nil, err
	}

	// Set random User-Agent
	userAgent := getRandomUserAgent()
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to fetch MobaXterm website: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("", fmt.Errorf("bad status code: %d", resp.StatusCode))
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Error("", fmt.Errorf("failed to parse HTML: %v", err))
		return nil, err
	}

	result := map[string]interface{}{
		"selectors":    map[string]interface{}{},
		"alternatives": map[string]interface{}{},
		"request_info": map[string]interface{}{
			"user_agent": userAgent,
		},
	}

	// Detailed inspection of version_titre elements
	versionTitreDetails := []map[string]string{}
	doc.Find("p.version_titre").Each(func(i int, s *goquery.Selection) {
		// Create a detailed report for each element
		detail := map[string]string{
			"text": s.Text(),
			// "html":      s.Html(),
			// "outerHTML": goquery.OuterHtml(s),
			"index": fmt.Sprintf("%d", i),
		}

		// Get parent info to understand the structure
		parent := s.Parent()
		if parent.Length() > 0 {
			parentHtml, _ := parent.Html()
			parentClass, _ := parent.Attr("class")
			detail["parent_class"] = parentClass
			detail["parent_html_preview"] = truncateString(parentHtml, 100)
		}

		// Extract version using regex
		matches := regexp.MustCompile(`Version (\d+\.\d+)`).FindStringSubmatch(s.Text())
		if len(matches) >= 2 {
			detail["extracted_version"] = matches[1]
		} else {
			detail["extracted_version"] = "not found"
		}

		versionTitreDetails = append(versionTitreDetails, detail)
	})
	result["version_titre_details"] = versionTitreDetails
	result["version_titre_count"] = len(versionTitreDetails)

	// Also check simple container presence
	simpleTitresCount := doc.Find("p.version_titre").Length()
	result["simple_version_titre_count"] = simpleTitresCount

	// Check if the selectors actually point to valid elements
	doc.Find(".page-block .container .tight_paragraphe").Each(func(i int, s *goquery.Selection) {
		hasVersionTitle := s.Find("p.version_titre").Length() > 0
		if hasVersionTitle {
			result["found_container_with_version_titre"] = true
		}
	})

	// Extract the full HTML structure of containers that might contain our target
	containerHTMLs := []string{}
	doc.Find(".page-block .container").Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		containerHTMLs = append(containerHTMLs, truncateString(html, 300))
	})
	result["container_htmls"] = containerHTMLs

	return result, nil
}

// Helper function to truncate long strings for debug output
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func GenerateLicense(count int, username, version string, c *gin.Context) {
	// 检查参数是否为空或不符合要求
	if username == "" || version == "" || count <= 0 {
		logger.Error("", fmt.Errorf("参数错误: %s, %s, %d", username, version, count))
	}

	// 拆分版本号
	versionArr := strings.Split(version, ".")
	if len(versionArr) != 2 {
		logger.Error("", fmt.Errorf("版本号格式错误: %s", version))
	}

	// 检查版本号是否为数字
	if _, err := strconv.Atoi(versionArr[0]); err != nil {
		// panic("版本号格式错误")
		log.Printf("版本号格式错误: %s", version)
	}
	if _, err := strconv.Atoi(versionArr[1]); err != nil {
		// panic("版本号格式错误")
		log.Printf("版本号格式错误: %s", version)
	}

	// 提取主次版本号
	major, _ := strconv.ParseInt(versionArr[0], 10, 64)
	minor, _ := strconv.ParseInt(versionArr[1], 10, 64)

	license := generateLicense(1, count, username, major, minor)
	// 先写入文件，解决直接输出压缩包文件大小不对导致无法使用的问题
	toFile(license)
	// 读取文件, 输出到浏览器
	c.FileAttachment(config.GetConfig().DataDir+"/Custom.mxtpro", "Custom.mxtpro")
	// 删除文件
	err := os.Remove(config.GetConfig().DataDir + "/Custom.mxtpro")
	if err != nil {
		logger.Error("", fmt.Errorf("删除文件失败: %v", err))
	}
}

func toFile(license []byte) {
	fileName := config.GetConfig().DataDir + "/Custom.mxtpro"
	_ = os.Remove(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		logger.Error("", fmt.Errorf("创建文件失败: %v", err))
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Error("", fmt.Errorf("关闭文件失败: %v", err))
		}
	}(f)

	zipFile := zip.NewWriter(f)
	defer func(zipFile *zip.Writer) {
		err := zipFile.Close()
		if err != nil {
			logger.Error("", fmt.Errorf("关闭ZIP写入器失败: %v", err))
		}
	}(zipFile)
	header := &zip.FileHeader{
		Name:               "Pro.key",
		Method:             zip.Store,
		CompressedSize64:   38,
		UncompressedSize64: 38,
	}
	// TODO 关键点，放到FileHeader会导致文件大小不对
	header.SetModTime(time.Now())
	proFile, err := zipFile.CreateRaw(header)
	if err != nil {
		logger.Error("", fmt.Errorf("创建ZIP文件失败: %v", err))
		return
	}
	_, err = proFile.Write(license)
	if err != nil {
		logger.Error("", fmt.Errorf("写入ZIP文件失败: %v", err))
		return
	}
	return
}

func generateLicense(userType, count int, username string, major, minor int64) []byte {
	licenseString := fmt.Sprintf("%d#%s|%d%d#%d#%d3%d6%d#%d#%d#%d#", userType, username, major, minor, count, major, minor, minor, 0, 0, 0)
	return variantBase64Encode(encryptBytes(0x787, []byte(licenseString)))
}

func encryptBytes(key int, bs []byte) []byte {
	var result []byte
	for _, b := range bs {
		encryptedByte := b ^ byte((key>>8)&0xff)
		result = append(result, encryptedByte)
		key = (int(result[len(result)-1]) & key) | 0x482D
	}
	return result
}

var (
	variantBase64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	variantBase64Map   = func() map[int]byte {
		result := make(map[int]byte)
		for i, v := range variantBase64Table {
			result[i] = byte(v)
		}
		return result
	}()
)

func variantBase64Encode(bs []byte) []byte {
	var result []byte
	blocksCount := len(bs) / 3
	leftBytes := len(bs) % 3
	for i := 0; i < blocksCount; i++ {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*i : 3*i+3])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>12)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>18)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
	}
	if leftBytes == 0 {
		return result
	} else if leftBytes == 1 {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*blocksCount:])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
		return result
	} else {
		var blocks []byte
		codingInt := littleEndianBytes(bs[3*blocksCount:])
		block := variantBase64Map[codingInt&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>6)&0x3f]
		blocks = append(blocks, block)
		block = variantBase64Map[(codingInt>>12)&0x3f]
		blocks = append(blocks, block)
		result = append(result, blocks...)
		return result
	}
}

func littleEndianBytes(bs []byte) int {
	var result = int(bs[0])
	for i := 1; i < len(bs); i++ {
		result = result | int(bs[i])<<(8*i)
	}
	return result
}
