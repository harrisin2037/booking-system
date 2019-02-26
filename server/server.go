package server

import (
	"BookingServer/model"
	"errors"
	"startkit"
	"startkit/library/web/apis"

	"github.com/gin-gonic/gin"

	"github.com/danilopolani/gocialite"
	broadcast "github.com/dustin/go-broadcast"
)

type (
	BookingServer struct {
		*startkit.Context
		*gocialite.Dispatcher
		RoomChannels map[string]broadcast.Broadcaster
		Booking      BookingInstanse
	}

	BookingInstanse struct {
		IsTesting                  bool
		PayPalKey                  string
		MasterKey                  string
		VisaKey                    string
		AEKey                      string
		OfficialEmail              string
		OfficialPhone              string
		ResetPasswordTokenDuration int
		SupportOauths              []string
		FacebookClientID           string
		FacebookClientSecret       string
		GoogleClientID             string
		GoogleClientSecret         string
		RoomChannelsBufferLength   int
	}
)

var (
	ErrorResponseFunc = map[int]func(value interface{}) gin.H{
		10000: func(value interface{}) gin.H {
			return gin.H{"response": apis.Resp("Internal Server Error For Binding The Request", value)}
		},
		10001: func(value interface{}) gin.H {
			return gin.H{"response": apis.Resp("Invalid Request For Fields Validation", value)}
		},
		10002: func(value interface{}) gin.H {
			return gin.H{"response": apis.Resp("Request Field(s) Value Missing", value)}
		},
		10003: func(value interface{}) gin.H {
			return gin.H{"response": apis.Resp("Internal Server Error (Database)", value)}
		},
		10004: func(value interface{}) gin.H {
			return gin.H{"response": apis.Resp("Record Is Already Existed", value)}
		},
		10005: func(value interface{}) gin.H {
			return gin.H{"response": apis.Resp("Not Supported Oauth Method", value)}
		},
	}
)

func (m *BookingServer) Standard(c *startkit.Context, group *gin.RouterGroup) (routes []gin.IRoutes) {
	routes = []gin.IRoutes{
		m.GETLogin(group, c),
		m.POSTLogin(group),
		m.GETLoginOauth(group),
		m.GETLoginOauthCallback(group),
		m.POSTRegister(group),
		// TestSession(group),
	}
	return
}

func (m *BookingServer) Version1(c *startkit.Context, group *gin.RouterGroup) (routes []gin.IRoutes) {
	group = group.Group("/version1")
	var (
		obj      = model.User{}
		checkers = []func(obj interface{}) (error, bool){
			func(obj interface{}) (err error, ok bool) {
				var (
					user = obj.(*model.User)
				)
				if !user.IsActive {
					err = errors.New("User Is Not Activated Yet")
				}
				return err, !user.IsActive
			},
		}
	)
	group.Use(m.Server.SessionVarification("jwt", m.Mysql.DB, &obj, checkers))
	routes = []gin.IRoutes{
		m.GETAssets(group, c),
		m.GETUsers(group, c),
		m.GETOrderRecord(group, c),
		// TestSession(group),
	}
	return
}

// func TestSession(group *gin.RouterGroup) gin.IRoutes {
// 	return group.Handle("GET", "/session", func(c *gin.Context) {
// 		session := sessions.Default(c)
// 		var count int
// 		v := session.Get("count")
// 		if v == nil {
// 			count = 0
// 		} else {
// 			count = v.(int)
// 			count++
// 		}
// 		type Resp struct {
// 			Message interface{} `json:"message,omitempty"`
// 			Data    interface{} `json:"data,omitempty"`
// 		}
// 		session.Set("count", count)
// 		session.Save()
// 		c.JSON(200, gin.H{"error_message": Resp{
// 			Message: "Token Claim Is Not Able To Get",
// 			Data:    "Error: ",
// 		}})
// 	})
// }
