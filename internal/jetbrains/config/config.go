package config

import (
	"license/internal/jetbrains/types"
	"time"
)

// DefaultLicenseConfig provides default configuration for license generation
var DefaultLicenseConfig = types.LicenseConfig{
	DefaultValidDays:   1095, // 3 years
	DefaultGracePeriod: 7,
	DefaultMetadata:    "0120231110PSAA003008",
	DefaultHash:        "51149839/0:-1370131430",
	EnableAI:           true,
	EnableTrial:        false,
}

// DefaultProductCodes contains built-in JetBrains product codes
var DefaultProductCodes = []string{
	"II",  // IntelliJ IDEA
	"PS",  // PhpStorm
	"AC",  // AppCode
	"DB",  // DataGrip
	"RM",  // RubyMine
	"WS",  // WebStorm
	"RD",  // Rider
	"CL",  // CLion
	"PC",  // PyCharm
	"GO",  // GoLand
	"DS",  // DataSpell
	"DC",  // dotCover
	"DPN", // dotPeek
	"DM",  // dotMemory
}

// ProductNames maps product codes to their names
var ProductNames = map[string]string{
	"II":  "IntelliJ IDEA",
	"PS":  "PhpStorm",
	"AC":  "AppCode",
	"DB":  "DataGrip",
	"RM":  "RubyMine",
	"WS":  "WebStorm",
	"RD":  "Rider",
	"CL":  "CLion",
	"PC":  "PyCharm",
	"GO":  "GoLand",
	"DS":  "DataSpell",
	"DC":  "dotCover",
	"DPN": "dotPeek",
	"DM":  "dotMemory",
}

// CalculateEffectiveDate calculates the effective date based on configuration
func CalculateEffectiveDate(validDays int) string {
	if validDays <= 0 {
		validDays = DefaultLicenseConfig.DefaultValidDays
	}
	
	now := time.Now()
	// End of today
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	// Add valid days
	effectiveDate := endOfToday.AddDate(0, 0, validDays)
	
	return effectiveDate.Format("2006-01-02")
}

// FormatDateTimeString formats a date for license use
func FormatDateTimeString(date time.Time) string {
	return date.Format("2006-01-02")
}

// GetProductName returns the product name for a given code
func GetProductName(code string) string {
	if name, ok := ProductNames[code]; ok {
		return name
	}
	return code
}

// IsValidProductCode checks if a product code is valid
func IsValidProductCode(code string) bool {
	_, exists := ProductNames[code]
	return exists
}

// MergeProductCodes merges multiple product code slices and removes duplicates
func MergeProductCodes(codeSets ...[]string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, codes := range codeSets {
		for _, code := range codes {
			if !seen[code] {
				seen[code] = true
				result = append(result, code)
			}
		}
	}
	
	return result
}