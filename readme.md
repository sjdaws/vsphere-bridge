# vSphere Bridge

Perform vSphere REST actions in a single call.

## Background

vSphere requires an API token when performing actions. This makes it difficult to perform actions as a webhook since each action requires multiple API calls to authenticate, act, and logout.

This application acts as a bridge for webhooks by exposing an API which will authenticate with vSphere prior to performing an action.

## Configuration

This application can be configured with command line options of with environment variables. If both an command line option and environment variable is set, the command line option will be preferred.

### Command line options

| Flag         | Type    | Description                                                                                  | Mandatory |
|--------------|---------|----------------------------------------------------------------------------------------------|-----------|
| `--fqdn`     | string  | The fully qualified domain name of the API server include scheme, e.g. https://vsphere.local | Y         |
| `--insecure` | boolean | If set to true the SSL certificate presented by the API server will not be verified          | N         |
| `--port`     | int     | The port to run the bridge on, defaults to 8000                                              | N         |

### Environment variables

| Key              | Description                                                                                  | Mandatory     |
|------------------|----------------------------------------------------------------------------------------------|---------------|
| ALLOW_INSECURE   | If set to true the SSL certificate presented by the API server will not be verified          | N             |
| BRIDGE_PORT      | The port to run the bridge on, defaults to 8000                                              | N             |
| VSPHERE_FQDN     | The fully qualified domain name of the API server include scheme, e.g. https://vsphere.local | Y             |
| VSPHERE_PASSWORD | The password for the account which has access to the API server                              | N<sup>1</sup> |
| VSPHERE_USERNAME | The username for the account which has access to the API server                              | N<sup>1</sup> |

<sup>1</sup> Credentials are mandatory but can be sent with the webhook rather than setting them as an environment variable. See <a href="#authentication">authentication</a>.

## Usage

Currently only power management for virtual machines is supported.

### Authentication

If credentials are set as an evironment variable the bridge will accept and process unauthenticated requests. To provide a layer of protection, credentials can be sent as a <a href="https://en.wikipedia.org/wiki/Basic_access_authentication#Client_side" target="_blank">basic authentication header</a>. The credentials in the header will be passed through to the vSphere API.

### Endpoints

| Endpoint             | Description                                                                             |
|----------------------|-----------------------------------------------------------------------------------------|
| `/power/:vm`         | Get power state for a virtual machine. `:vm` must be a valid name of a virtual machine. |
| `/power/cycle/:vm`   | Power a virtual machine off and on. `:vm` must be a valid name of a virtual machine.    |
| `/power/on/:vm`      | Power on a virtual machine. `:vm` is the friendly name of a virtual machine.            |
| `/power/off/:vm`     | Power off a virtual machine. `:vm` is the friendly name of a virtual machine.           |
| `/power/reset/:vm`   | Reset a virtual machine. `:vm` is the friendly name of a virtual machine.               |
| `/power/suspend/:vm` | Suspend a virtual machine. `:vm` is the friendly name of a virtual machine.             |
