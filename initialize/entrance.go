package initialize

func ExecuteInitialize() {
	// Initialize certificates
	InitCert()
	// Initialize GitLab
	InitGitLabCert()
	// Initialize JetBrains
	InitJetbrains()
}
