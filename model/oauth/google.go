package oauth

type OauthGoogle struct {
	GoogleID              string `json:"google_id,omitempty" gorm:"column:google_id"`
	GoogleUserInformation JSON   `json:"google_user_information,omitempty" gorm:"column:google_user_information; sql:type:json"`
}
