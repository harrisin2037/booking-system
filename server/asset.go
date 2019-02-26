package server

import (
	"BookingServer/model"
	"startkit"
	"startkit/library/web/apis"

	"github.com/gin-gonic/gin"
)

/*
	TODO:
	GET: asset/:id
	GET: asset
	POST: asset
	POST: asset/:id
	DELETE: asset
	DELETE: asset/:id
*/

func (m *BookingServer) GETAssets(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return apis.
		New(c).
		ReqPath("GET", "/asset").
		Pagination("", "").
		Model(model.Asset{}).
		Table(&model.Asset{}).
		Find(&[]model.Asset{}).
		Handle(group)
}

func (m *BookingServer) GETAsset(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return apis.
		New(c).
		ReqPath("GET", "/asset").
		ID("external_id").
		Model(model.Asset{}).
		Table(&model.Asset{}).
		Find(&model.Asset{}).
		Handle(group)
}
