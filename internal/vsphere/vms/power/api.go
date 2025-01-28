package power

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sjdaws/vsphere-bridge/pkg/errors"
)

// virtualMachine api representation of a virtual machine.
type virtualMachine struct {
	CpuCount   int    `json:"cpu_count"`
	ID         string `json:"vm"`
	MemorySize int    `json:"memory_size_MiB"`
	Name       string `json:"name"`
	PowerState string `json:"power_state"`
}

// Cycle power cycle a virtual machine.
func (p *Power) Cycle(ctx echo.Context) error {
	err := p.performPowerAction(ctx, "stop", virtualMachine{Name: ctx.Param("vm")})
	if err != nil {
		return errors.Wrap(err, "unable to power off virtual machine power")
	}

	err = p.performPowerAction(ctx, "start", virtualMachine{Name: ctx.Param("vm")})
	if err != nil {
		return errors.Wrap(err, "unable to power on virtual machine power")
	}

	return ctx.JSON(http.StatusOK, map[string]any{"result": "ok"})
}

// Get power state of a virtual machine.
func (p *Power) Get(ctx echo.Context) error {
	vm, err := p.getVirtualMachineByName(ctx, ctx.Param("vm"))
	if err != nil {
		return errors.Wrap(err, "unable to get virtual machine")
	}

	return ctx.JSON(http.StatusOK, map[string]any{"result": "ok", "state": vm.PowerState})
}

// Off power down a virtual machine.
func (p *Power) Off(ctx echo.Context) error {
	err := p.performPowerAction(ctx, "stop", virtualMachine{Name: ctx.Param("vm")})
	if err != nil {
		return errors.Wrap(err, "unable to power off virtual machine power")
	}

	return ctx.JSON(http.StatusOK, map[string]any{"result": "ok"})
}

// On power up a virtual machine.
func (p *Power) On(ctx echo.Context) error {
	err := p.performPowerAction(ctx, "start", virtualMachine{Name: ctx.Param("vm")})
	if err != nil {
		return errors.Wrap(err, "unable to power on virtual machine power")
	}

	return ctx.JSON(http.StatusOK, map[string]any{"result": "ok"})
}

// Reset a virtual machine.
func (p *Power) Reset(ctx echo.Context) error {
	err := p.performPowerAction(ctx, "reset", virtualMachine{Name: ctx.Param("vm")})
	if err != nil {
		return errors.Wrap(err, "unable to reset virtual machine")
	}

	return ctx.JSON(http.StatusOK, map[string]any{"result": "ok"})
}

// Suspend a virtual machine.
func (p *Power) Suspend(ctx echo.Context) error {
	err := p.performPowerAction(ctx, "suspend", virtualMachine{Name: ctx.Param("vm")})
	if err != nil {
		return errors.Wrap(err, "unable to suspend virtual machine")
	}

	return ctx.JSON(http.StatusOK, map[string]any{"result": "ok"})
}

// getVirtualMachineByName get a virtualMachine by name.
func (p *Power) getVirtualMachineByName(ctx echo.Context, name string) (virtualMachine, error) {
	response, err := p.vsphere.Request(ctx, http.MethodGet, "/vcenter/vm", nil)
	if err != nil {
		return virtualMachine{}, errors.Wrap(err, "unable to fetch list of virtual machines")
	}

	vms := make([]virtualMachine, 0)
	err = json.Unmarshal(response, &vms)
	if err != nil {
		return virtualMachine{}, errors.Wrap(err, "unable to unmarshal virtual machine list")
	}

	for _, vm := range vms {
		if vm.Name == name {
			return vm, nil
		}
	}

	return virtualMachine{}, errors.New("virtual machine %s not found", name)
}

// performPowerAction perform a power action on a virtualMachine.
func (p *Power) performPowerAction(ctx echo.Context, action string, vm virtualMachine) error {
	var err error
	if vm.ID == "" {
		vm, err = p.getVirtualMachineByName(ctx, vm.Name)
		if err != nil {
			return errors.Wrap(err, "unable to find virtual machine")
		}
	}

	_, err = p.vsphere.Request(ctx, http.MethodPost, fmt.Sprintf("/vcenter/vm/%s/power?action=%s", vm.ID, action), nil)
	if err != nil {
		return errors.Wrap(err, "unable to perform virtual machine power action")
	}

	return nil
}
