package vsphere

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/sjdaws/vsphere-bridge/pkg/errors"
)

// header http header.
type header struct {
	key   string
	value string
}

// Request send an authenticated http request.
func (v *Vsphere) Request(ctx echo.Context, method string, path string, payload io.Reader) ([]byte, error) {
	if v.token == "" {
		err := v.authenticate(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to authenticate with vsphere api")
		}
		defer v.logout(ctx)
	}

	return v.request(method, path, payload)
}

// request send an http request.
func (v *Vsphere) request(method string, path string, payload io.Reader, headers ...header) ([]byte, error) {
	request, err := http.NewRequest(method, fmt.Sprintf("%s/api/%s", v.config.Server, strings.TrimPrefix(path, "/")), payload)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create http request")
	}

	for _, reqHeader := range headers {
		request.Header.Add(reqHeader.key, reqHeader.value)
	}

	if v.token != "" {
		request.Header.Set("vmware-api-session-id", v.token)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if v.config.Insecure {
		client = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "error sending http request")
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("unexpected response received from server (%s): %s", response.Status, string(body))
	}

	return body, nil
}
