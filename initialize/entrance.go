package initialize

// ExecuteInitialize initializes various components
func ExecuteInitialize() error {
	// Initialize certificates
	if err := InitCert(); err != nil {
		return err
	}

	// Initialize GitLab
	if err := InitGitLabCert(); err != nil {
		return err
	}

	// Initialize JetBrains
	if err := InitJetbrains(); err != nil {
		return err
	}

	return nil
}
