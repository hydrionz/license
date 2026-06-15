package task

import (
	v2 "license/internal/jetbrains/code/service/v2"
	"license/internal/logger"
)

// FetchProductLatest refreshes products and plugins from the upstream catalog.
func FetchProductLatest() {
	if err := v2.FetchLatestProducts(); err != nil {
		logger.Error("Failed to fetch latest product:", err)
	}
	if err := v2.FetchLatestPlugins(); err != nil {
		logger.Error("Failed to fetch latest plugin:", err)
	}
}
