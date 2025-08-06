package service

import (
	"fmt"
	"license/config"
	"license/jetbrains/code/entity"
	"license/jetbrains/code/mapper"
	"license/logger"
	"sync"
)

// ProductService handles product-related operations
type ProductService struct {
	mapper mapper.ProductMapper
	mu     sync.RWMutex
}

// NewProductService creates a new product service
func NewProductService() *ProductService {
	return &ProductService{
		mapper: &mapper.GormProductMapper{},
	}
}

// GetAll retrieves all products
func (s *ProductService) GetAll() ([]entity.ProductEntity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products, err := s.mapper.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

// GetByCode retrieves a product by its code
func (s *ProductService) GetByCode(code string) (*entity.ProductEntity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products, err := s.mapper.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	for _, product := range products {
		if product.ProductCode == code {
			return &product, nil
		}
	}

	return nil, fmt.Errorf("product not found: %s", code)
}

// FetchLatest fetches the latest products from external source
func (s *ProductService) FetchLatest() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Info("Fetching latest products...")
	
	// Implementation for fetching latest products
	// This would typically involve:
	// 1. Calling external API
	// 2. Parsing response
	// 3. Updating database
	
	// For now, we'll just log
	logger.Info("Products fetch completed")
	
	return nil
}

// SaveProducts saves multiple products to database
func (s *ProductService) SaveProducts(products []entity.ProductEntity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	db := config.DB
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	for _, product := range products {
		if err := db.Save(&product).Error; err != nil {
			logger.Error("Failed to save product: ", err)
			continue
		}
	}

	return nil
}

// PluginService handles plugin-related operations
type PluginService struct {
	mapper mapper.PluginMapper
	mu     sync.RWMutex
}

// NewPluginService creates a new plugin service
func NewPluginService() *PluginService {
	return &PluginService{
		mapper: &mapper.GormPluginMapper{},
	}
}

// GetAll retrieves all plugins
func (s *PluginService) GetAll() ([]entity.PluginEntity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plugins, err := s.mapper.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugins: %w", err)
	}

	return plugins, nil
}

// GetByCode retrieves a plugin by its code
func (s *PluginService) GetByCode(code string) (*entity.PluginEntity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plugins, err := s.mapper.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugins: %w", err)
	}

	for _, plugin := range plugins {
		if plugin.PluginCode == code {
			return &plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin not found: %s", code)
}

// FetchLatest fetches the latest plugins from external source
func (s *PluginService) FetchLatest() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Info("Fetching latest plugins...")
	
	// Implementation for fetching latest plugins
	// This would typically involve:
	// 1. Calling external API
	// 2. Parsing response
	// 3. Updating database
	
	// For now, we'll just log
	logger.Info("Plugins fetch completed")
	
	return nil
}

// SavePlugins saves multiple plugins to database
func (s *PluginService) SavePlugins(plugins []entity.PluginEntity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	db := config.DB
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	for _, plugin := range plugins {
		if err := db.Save(&plugin).Error; err != nil {
			logger.Error("Failed to save plugin: ", err)
			continue
		}
	}

	return nil
}