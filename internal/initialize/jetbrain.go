package initialize

import (
	"license/internal/jetbrains/util"
	"license/internal/logger"
)

// InitJetbrains initialize JetBrains components
func InitJetbrains() error {
	logger.Info("Initializing JetBrains certificates")

	// Initialize certificate paths for JetBrains
	if err := util.InitCertPaths(GetCertManager()); err != nil {
		return err
	}

	fakeCert := util.GetFake()

	// Load or generate keys
	if err := fakeCert.LoadOrGenerate(); err != nil {
		logger.Error("Failed to load or generate keys: %v", err)
		return err
	}

	// Load root certificates
	if err := fakeCert.LoadRootCert(); err != nil {
		logger.Error("Failed to load root certificates: %v", err)
		return err
	}

	// Generate certificates if needed
	if err := fakeCert.GenerateRootCert(); err != nil {
		logger.Error("Failed to generate certificates: %v", err)
		return err
	}

	// Load certificates
	if err := fakeCert.LoadCert(); err != nil {
		logger.Error("Failed to load certificates: %v", err)
		return err
	}

	logger.Info("JetBrains certificates initialized successfully")
	return nil
}
