package initialize

import (
	"license/internal/gitlab/service"
)

func InitGitLabCert() error {
	return service.LoadKeys()
}
