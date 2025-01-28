package vsphere

import (
	"github.com/sjdaws/vsphere-bridge/internal/configuration"
	"github.com/sjdaws/vsphere-bridge/pkg/logging"
)

// Vsphere instance of vsphere.
type Vsphere struct {
	config *configuration.Configuration
	logger logging.Logger
	token  string
}

// New create a new Vsphere instance.
func New(config *configuration.Configuration, logger logging.Logger) *Vsphere {
	return &Vsphere{
		config: config,
		logger: logger,
	}
}
