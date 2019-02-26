package model

import (
	"startkit"
	"startkit/library/random"
	"startkit/library/times"
	"time"

	uuid "gopub/src/github.com/satori/go.uuid"
)

func Models() []interface{} {
	var (
		user               = &User{}
		oauthToken         = &OauthToken{}
		resetPasswordToken = &ResetPasswordToken{}
		asset              = &Asset{}
		requestRecord      = &RequestRecord{}
		orderRecord        = &OrderRecord{}
		paymentRecord      = &PaymentRecord{}
		invoiceRecord      = &InvoiceRecord{}
	)
	return []interface{}{
		user,
		oauthToken,
		resetPasswordToken,
		asset,
		requestRecord,
		orderRecord,
		paymentRecord,
		invoiceRecord,
	}
}

func TimeZero() time.Time {
	return times.Zero()
}

/*--------------------------------Testing Data------------------------------------*/

func TestingData(m *startkit.Context, b bool) {
	if b {
		defer m.Mysql.Connector()()
		var (
			count = 10
			users = []User{}
			// id    = 1
			roles = []string{"admin", "staff", "boat_owner", "client"}
		)
		for i := 0; i < count; i++ {
			var (
				user = User{
					IsActive:    false,
					Username:    "testing user" + uuid.NewV4().String(),
					Role:        RoleType(random.ShuffleStringArray(roles)[0]),
					Firstname:   uuid.NewV4().String(),
					Lastname:    "user" + uuid.NewV4().String(),
					Password:    "" + uuid.NewV4().String(),
					Email:       uuid.NewV4().String() + "@gmail.com",
					ActivatedAt: times.Zero(),
				}
			)
			// user.ID = id
			users = append(users, user)
			// id++
		}
		for i := range users {
			m.Mysql.DB.Debug().Create(&users[i])
		}
	}
}

/*--------------------------------------------------------------------------------*/
