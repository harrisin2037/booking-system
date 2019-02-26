package server

import (
	"BookingServer/model"
	"startkit"
	"startkit/library/web/apis"

	"github.com/gin-gonic/gin"
)

func (m *BookingServer) GETUsers(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return apis.
		New(c).
		ReqPath("GET", "/user").
		Pagination("", "").
		Model(model.User{}).
		Table(&model.User{}).
		Find(&[]model.User{}).
		Handle(group)
}

func (m *BookingServer) GETUser(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return apis.
		New(c).
		ReqPath("GET", "/user").
		ID("external_id").
		Model(model.User{}).
		Table(&model.User{}).
		Find(&[]model.User{}).
		Handle(group)
}
