package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"github.com/sjdaws/vsphere-bridge/internal/configuration"
	"github.com/sjdaws/vsphere-bridge/internal/vsphere"
	"github.com/sjdaws/vsphere-bridge/internal/vsphere/vms/power"
	"github.com/sjdaws/vsphere-bridge/pkg/logging"
)

const usageText = `
usage: %s [OPTIONS]

Perform vsphere REST actions in a single call

Options:

  --fqdn string      The fqdn of the target vsphere instance including scheme, e.g. http://vsphere.local
  --insecure bool    Allow insecure SSL connections to vsphere instance
  --port int         The port to run the bridge on, defaults to 8000

Environment variables:

  ALLOW_INSECURE string      Allow insecure SSL connections to vsphere instance
  BRIDGE_PORT int			 The port to run the bridge on, defaults to 8000
  VSPHERE_FQDN string        The fqdn of the target vsphere instance including scheme, e.g. http://vsphere.local
  VSPHERE_PASSWORD string    Password for vsphere account with API access
  VSPHERE_USERNAME string    Username for vsphere account with API access

FQDN is mandatory, the rest of the parameters are optional.

Option will be used if both option and environment variable are passed for the same parameter.

`

func main() {
	// Check if help is requested
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			fmt.Printf(usageText, os.Args[0])
			os.Exit(0)
		}
	}

	logger := logging.Default()

	config, err := configuration.Resolve()
	if err != nil {
		logger.Error(err)

		fmt.Printf(usageText, os.Args[0])
		os.Exit(0)
	}

	server := echo.New()
	server.HTTPErrorHandler = func(err error, ctx echo.Context) {
		logger.Error(err)
		_ = ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	server.GET("/health", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api := vsphere.New(config, logger)
	power.New(*api, server)

	err = server.Start(":" + config.Port)
	if err != nil {
		logger.Fatal(err)
	}
}
