package vsphere

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/sjdaws/vsphere-bridge/pkg/errors"
)

// authenticate to the vsphere API.
func (v *Vsphere) authenticate(ctx echo.Context) error {
	// Only create an auth header if one doesn't exist
	credentials := strings.TrimSpace(ctx.Request().Header.Get("Authorization"))
	if credentials == "" && v.config.Username != "" && v.config.Password != "" {
		credentials = "Basic " + base64.StdEncoding.EncodeToString([]byte(v.config.Username+":"+v.config.Password))
	}

	// If no credentials, there isn't much we can do
	if credentials == "" {
		return errors.New("one of: vsphere username and password, basic authorization header are required")
	}

	response, err := v.request(http.MethodPost, "/session", nil, header{key: "Authorization", value: credentials})
	if err != nil {
		return errors.Wrap(err, "unable to fetch session token")
	}

	// Body will contain api token, but it is also quoted for some wierd reason so trim off quotes
	v.token = strings.Trim(string(response), `"`)

	return nil
}

// logout from the vsphere API.
func (v *Vsphere) logout(ctx echo.Context) {
	if v.token == "" {
		return
	}

	_, err := v.Request(ctx, http.MethodDelete, "/session", nil)
	if err != nil {
		v.logger.Error("unable to logout of session: %v", err)
	}

	v.token = ""
}
