package main

import (
	"BookingServer/model"
	"BookingServer/server"
	"startkit"
	"startkit/library/files"

	"github.com/danilopolani/gocialite"
	broadcast "github.com/dustin/go-broadcast"
)

func main() {
	var (
		m = server.BookingServer{}
	)
	files.BindFileToObj("bookingserver.ini", &m)
	m.Context = startkit.New("setting.ini")
	for _, model := range model.Models() {
		m.Mysql.AutoMigrateByAddr(model)
	}
	model.TestingData(m.Context, m.Booking.IsTesting)
	m.Dispatcher = gocialite.NewDispatcher()
	m.RoomChannels = make(map[string]broadcast.Broadcaster)
	m.Run(StartBookingServer(&m))
}

func StartBookingServer(m *server.BookingServer) func(c *startkit.Context) {
	return func(c *startkit.Context) {
		var (
			group = c.Server.Engine.Group("/booking")
		)
		m.Standard(c, group)
		m.Version1(c, group)
		c.Server.Start()
	}
}
