package power

import (
	"github.com/labstack/echo/v4"

	"github.com/sjdaws/vsphere-bridge/internal/vsphere"
)

type Power struct {
	vsphere vsphere.Vsphere
}

// New create a new power instance.
func New(vsphere vsphere.Vsphere, server *echo.Echo) *Power {
	api := &Power{
		vsphere: vsphere,
	}

	group := server.Group("/power")
	group.GET("/:vm", api.Get)
	group.POST("/cycle/:vm", api.Cycle)
	group.POST("/off/:vm", api.Off)
	group.POST("/on/:vm", api.On)
	group.POST("/reset/:vm", api.Reset)
	group.POST("/suspend/:vm", api.Suspend)

	return api
}
