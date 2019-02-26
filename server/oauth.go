package server

import (
	"BookingServer/model"
	"net/http"
	"startkit/library/gins"
	"startkit/library/web/apis"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (m *BookingServer) CheckOauth(provider string) (exist bool) {
	if count := len(m.Booking.SupportOauths); count > 0 {
		for i := 0; i < count; i++ {
			if exist = m.Booking.SupportOauths[i] == provider; exist {
				return
			}
		}
	}
	return false
}

func (m *BookingServer) GETLoginOauth(group *gin.RouterGroup) gin.IRoutes {
	return group.Handle("GET", "/login/:provider", func(c *gin.Context) {
		if !m.CheckOauth(c.Param("provider")) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Not Supported Oauth Method")},
			)
			return
		}
		var (
			provider        = c.Param("provider")
			providerSecrets = map[string]map[string]string{
				"facebook": {
					"clientID":     m.Booking.FacebookClientID,
					"clientSecret": m.Booking.FacebookClientSecret,
					"redirectURL":  m.Server.Domain + "/login/facebook/callback",
				},
				"google": {
					"clientID":     m.Booking.GoogleClientID,
					"clientSecret": m.Booking.GoogleClientSecret,
					"redirectURL":  m.Server.Domain + "/login/google/callback",
				},
			}
			providerScopes = map[string][]string{
				"facebook": []string{},
				"google":   []string{},
			}
			providerData = providerSecrets[provider]
			actualScopes = providerScopes[provider]
			authURL, err = m.Dispatcher.New().
					Driver(provider).
					Scopes(actualScopes).
					Redirect(
					providerData["clientID"],
					providerData["clientSecret"],
					providerData["redirectURL"],
				)
		)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Error: "+err.Error())},
			)
			return
		}
		c.Redirect(http.StatusFound, authURL)
		return
	})
}

func (m *BookingServer) GETLoginOauthCallback(group *gin.RouterGroup) gin.IRoutes {
	return group.Handle("GET", "/login/:provider/callback", func(c *gin.Context) {
		if !m.CheckOauth(c.Param("provider")) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Not Supported Oauth Method")},
			)
			return
		}
		var (
			jwtToken              = ""
			user                  = model.User{}
			state                 = c.Query("state")
			code                  = c.Query("code")
			provider              = c.Param("provider")
			session               = sessions.Default(c)
			oauthUser, token, err = m.Dispatcher.Handle(state, code)
		)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"response": apis.Resp("Failed", "Error: "+err.Error())},
			)
			return
		}
		session.Clear()
		session.Save()
		if err = m.Mysql.DB.
			Debug().
			Where(map[string]interface{}{
				"username": oauthUser.Username,
				"email":    oauthUser.Email,
			}).
			Find(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"response": apis.Resp("Failed", "Internal Server Error (Database)")},
			)
			return
		}
		if !user.IsEmpty() {
			if err = m.Mysql.DB.
				Debug().
				Model(&model.User{}).
				Where(map[string]interface{}{
					"username": oauthUser.Username,
					"email":    oauthUser.Email},
				).Updates(map[string]interface{}{
				"oauth_provider":     provider,
				provider + "_id":     oauthUser.ID,
				"oauth_access_token": token.AccessToken,
				"oauth_token_expiry": token.Expiry,
			}).Error; err != nil {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{"response": apis.Resp("Failed", "Internal Server Error (Database)")},
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
			session.Set("username", oauthUser.Username)
			session.Set("email", oauthUser.Email)
			session.Set("jwt", jwtToken)
			session.Save()
			c.Writer.Header().Set("Authorization", "Bearer "+jwtToken)
			c.JSON(
				http.StatusOK,
				gin.H{"response": apis.Resp("Success", map[string]interface{}{
					"provider":    provider,
					"user":        oauthUser,
					"oauth_token": token,
					"jwt_token":   jwtToken,
				})},
			)
			return
		}
		if err = m.Mysql.DB.Debug().Create(user.SetOauth(provider, token, oauthUser)).Error; err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"response": apis.Resp("Failed", "Internal Server Error (Database)")},
			)
			return
		}
		// send email for set password

		return
	})
}
