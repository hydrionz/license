package initialize

import (
	"license/gitlab/service"
)

func InitGitLabCert() error {
	return service.LoadKeys()
}
