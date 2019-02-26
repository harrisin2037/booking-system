package server

import (
	"BookingServer/model"
	"errors"
	"net/http"
	"startkit/library/gorms"
	"startkit/library/web/apis"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
)

func (m *BookingServer) POSTOauthRegister(group *gin.RouterGroup) gin.IRoutes {
	return group.Handle("POST", "/register/:provider", func(c *gin.Context) {
		if !m.CheckOauth(c.Param("provider")) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Not Supported Oauth Method")},
			)
			return
		}
		var (
			// user                  = model.User{}
			state                 = c.Query("state")
			code                  = c.Query("code")
			provider              = c.Param("provider")
			oauthUser, token, err = m.Dispatcher.Handle(state, code)
		)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Error: "+err.Error())},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			gin.H{"response": apis.Resp("Success", map[string]interface{}{
				"provider": provider,
				"user":     oauthUser,
				"token":    token,
			})},
		)
		return
		// user = model.User{
		// 	Username:  oauthUser.Username,
		// 	Firstname: oauthUser.FirstName,
		// 	Lastname:  oauthUser.LastName,
		// 	Email:     oauthUser.Email,
		// }
	})
}

// oauth later
func (m *BookingServer) POSTRegister(group *gin.RouterGroup) gin.IRoutes {
	return group.Handle("POST", "/register", func(c *gin.Context) {
		defer m.Mysql.Connector()()
		var (
			requestTime = time.Now()
			err         = errors.New("")
			user        = model.User{}
			// oauthProvider = model.OauthProviderType("null")
			req = struct {
				Email     string `json:"email" binding:"required"`
				Username  string `json:"username" binding:"required"`
				Password  string `json:"password" binding:"required"`
				Age       int    `json:"age,omitempty"`
				Firstname string `json:"firstname" binding:"required"`
				Lastname  string `json:"lastname" binding:"required"`
				Role      string `json:"role" binding:"required"`
				Address   string `json:"address" binding:"required"`
				Phone     string `json:"phone" binding:"required"`
				Region    string `json:"region" binding:"required"`
			}{}
		)
		if err = c.BindJSON(&req); err != nil {
			_, ok := err.(validator.ValidationErrors)
			if !ok {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"response": apis.Resp("Failed", "Error For Binding The Request")},
				)
			} else {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"response": apis.Resp("Failed", "Invalid Request")},
				)
			}
			return
		}
		if req.Email == "" ||
			req.Username == "" ||
			req.Password == "" ||
			req.Phone == "" ||
			req.Region == "" ||
			req.Age == 0 ||
			req.Firstname == "" ||
			req.Lastname == "" ||
			req.Role == "" {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Request Field(s) Value Missing")},
			)
			return
		}
		if err = m.Mysql.DB.
			Debug().
			Where(map[string]interface{}{"username": req.Username, "email": req.Email}).
			Find(&user).
			Error; err != nil && err != gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"response": apis.Resp("Failed", "Internal Server Error (Database)")},
			)
			return
		}
		if !user.IsEmpty() {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Record Is Already Existed")},
			)
			return
		}
		user = model.User{
			Email:     req.Email,
			Username:  req.Username,
			Password:  req.Password,
			Age:       req.Age,
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			Role:      model.RoleType(req.Role).CorrectTo(model.Client),
			// OauthProviderType: oauthProvider,
		}
		if err = m.Mysql.DB.Debug().Create(&user).Error; err != nil && !gorms.IsDuplicateError(err) {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"response": apis.Resp("Failed", "Internal Server Error (Database)")},
			)
			return
		} else if gorms.IsDuplicateError(err) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Record Already Existed, Error: "+err.Error())},
			)
			return
		}
		// send email for active the acount
		user.SetResetPasswordToken(requestTime, m.Booking.ResetPasswordTokenDuration)
		return
	})
}
