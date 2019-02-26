package oauth

type OauthFacebook struct {
	FacebookID              string `json:"facebook_id,omitempty" gorm:"column:facebook_id"`
	FacebookUserInformation JSON   `json:"facebook_user_information,omitempty" gorm:"column:facebook_user_information; sql:type:json"`
}
