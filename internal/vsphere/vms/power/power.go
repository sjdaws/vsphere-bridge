package power

import (
	"github.com/labstack/echo/v4"

	"github.com/sjdaws/vsphere-bridge/internal/vsphere"
	"github.com/sjdaws/vsphere-bridge/pkg/notifier"
)

type Power struct {
	notify  *notifier.Notifier
	vsphere vsphere.Vsphere
}

// New create a new power instance.
func New(vsphere vsphere.Vsphere, notify *notifier.Notifier, server *echo.Echo) *Power {
	api := &Power{
		notify:  notify,
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
