package configuration

import (
	"flag"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/carlmjohnson/truthy"

	"github.com/sjdaws/vsphere-bridge/pkg/errors"
)

// Configuration resolved configuration from os.Getenv and os.Args.
type Configuration struct {
	Insecure  bool
	NotifyURL string
	Password  string
	Port      string
	Server    *url.URL
	Username  string
	fqdn      string
}

type resolved map[string]string

// Resolve configuration from environment and command arguments.
func Resolve() (*Configuration, error) {
	env := resolveEnv()
	flags := resolveFlags()

	port := strings.TrimSpace(preferFlags(flags, env, "port"))

	// Prefer flags where possible
	config := &Configuration{
		Insecure:  truthy.Value(strings.TrimSpace(preferFlags(flags, env, "insecure"))),
		NotifyURL: strings.TrimSpace(preferFlags(flags, env, "notify_url")),
		Password:  strings.TrimSpace(env["password"]),
		Port:      truthy.Cond(port != "", port, "8000"),
		Username:  strings.TrimSpace(env["username"]),
		fqdn:      strings.TrimSpace(strings.TrimSuffix(preferFlags(flags, env, "fqdn"), "/")),
	}

	err := validate(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// preferFlags will return the flag key if it exists
func preferFlags(flags resolved, env resolved, key string) string {
	return truthy.Cond(flags[key] != "", flags[key], env[key])
}

// resolveEnv from environment variables.
func resolveEnv() resolved {
	return resolved{
		"fqdn":       os.Getenv("VSPHERE_FQDN"),
		"insecure":   truthy.Cond(truthy.Value(os.Getenv("ALLOW_INSECURE")), "true", "false"),
		"notify_url": os.Getenv("NOTIFY_URL"),
		"password":   os.Getenv("VSPHERE_PASSWORD"),
		"port":       os.Getenv("BRIDGE_PORT"),
		"username":   os.Getenv("VSPHERE_USERNAME"),
	}
}

// resolveFlags passed at run time.
func resolveFlags() resolved {
	var port = flag.Int("port", -1, "port to run bridge on")
	var insecure = flag.Bool("insecure", false, "disable tls certificate verification")
	var fqdn = flag.String("fqdn", "", "vsphere server fqdn")
	var notifyURL = flag.String("notify-url", "", "shoutrrr compatible notify url")
	flag.Parse()

	return resolved{
		"fqdn":       *fqdn,
		"insecure":   truthy.Cond(*insecure, "true", "false"),
		"notify_url": *notifyURL,
		"port":       truthy.Cond(*port >= 0, strconv.Itoa(*port), ""),
	}
}

// validate configuration.
func validate(config *Configuration) error {
	if config.fqdn == "" {
		return errors.New("vsphere fqdn is required, run %s --help for more information.", os.Args[0])
	}

	server, err := url.Parse(config.fqdn)
	if err != nil {
		return errors.Wrap(err, "invalid server URL")
	}

	_, err = strconv.Atoi(config.Port)
	if err != nil {
		return errors.Wrap(err, "invalid port number")
	}

	config.Server = server
	config.Port = config.Port

	return nil
}
