package server

import (
	"BookingServer/model"
	"startkit"
	"startkit/library/web/apis"

	"github.com/gin-gonic/gin"
)

func (m *BookingServer) GETOrderRecords(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return apis.
		New(c).
		ReqPath("GET", "/order_record").
		// Param([]string{"uuid"}).
		Pagination("", "").
		Model(model.OrderRecord{}).
		Table(&model.OrderRecord{}).
		Find(&[]model.OrderRecord{}).
		Handle(group)
}

func (m *BookingServer) GETOrderRecord(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return apis.
		New(c).
		ReqPath("GET", "/order_record").
		ID("external_id").
		Model(model.OrderRecord{}).
		Table(&model.OrderRecord{}).
		Find(&[]model.OrderRecord{}).
		Handle(group)
}
