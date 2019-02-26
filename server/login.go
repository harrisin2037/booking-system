package server

import (
	"BookingServer/model"
	"errors"
	"net/http"
	"startkit"
	"startkit/library/gins"
	"startkit/library/web/apis"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v8"
)

func (m *BookingServer) GETLogin(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return apis.
		New(c).
		ReqPath("GET", "/login").
		Static("login.html").
		Handle(group)
}

func (m *BookingServer) POSTLogin(group *gin.RouterGroup) gin.IRoutes {
	return group.Handle("POST", "/login", func(c *gin.Context) {
		defer m.Mysql.Connector()()
		var (
			err      = errors.New("")
			session  = sessions.Default(c)
			user     = model.User{}
			jwtToken = ""
			req      = struct {
				Email    string `json:"email" binding:"required"`
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}{}
		)
		session.Clear()
		session.Save()
		if err = c.BindJSON(&req); err != nil {
			_, ok := err.(validator.ValidationErrors)
			if !ok {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"response": apis.Resp("Failed", "Internal Server Error For Binding The Request")},
				)
			} else {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"response": apis.Resp("Failed", "Invalid Request")},
				)
			}
			return
		}
		if req.Email == "" || req.Username == "" || req.Password == "" {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Login Request Field(s) Value Missing")},
			)
			return
		}
		if err = m.Mysql.DB.Debug().Where(map[string]interface{}{"username": req.Username, "email": req.Email}).Find(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.AbortWithStatusJSON(
					http.StatusNotFound,
					gin.H{"response": apis.Resp("Failed", "User Record Not Found")},
				)
			} else {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{"response": apis.Resp("Failed", "Internal Server Error (Database)")},
				)
			}
			return
		}
		if !user.IsActive {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"response": apis.Resp("Failed", "User Account Is Not Yet Activated")},
			)
			return
		}
		if req.Password != "" {
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					gin.H{"response": apis.Resp("Failed", "Invalid Password")},
				)
				return
			}
		} else {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"response": apis.Resp("Failed", "Password Is Required But Missing")},
			)
			return
		}
		jwtToken, err = gins.JWT(user.ID, user.ExternalID, m.Server.JWTIssuer, m.Server.JWTSignedString)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"response": apis.Resp("Failed", "JWT Generation Failed")},
			)
			return
		}
		session.Set("username", req.Username)
		session.Set("email", req.Email)
		session.Set("jwt", jwtToken)
		session.Save()
		c.Writer.Header().Set("Authorization", "Bearer "+jwtToken)
		c.JSON(
			http.StatusOK,
			gin.H{"response": apis.Resp("Success", map[string]string{
				"jwt_token": jwtToken,
			})},
		)
		return
	})
}
