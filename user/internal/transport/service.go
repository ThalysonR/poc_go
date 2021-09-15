package transport

import "github.com/thalysonr/poc_go/user/internal/config"

type Service interface {
	ConfigChanged(config.Config) bool
	Start(config.Config) error
	Stop() error
}
