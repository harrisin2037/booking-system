package server

import (
	"io"
	"net/http"
	"startkit"
	"startkit/library/web/apis"

	broadcast "github.com/dustin/go-broadcast"
	"github.com/gin-gonic/gin"
	uuid "gopub/src/github.com/satori/go.uuid"
)

func (m *BookingServer) GETChatrooms(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("GET", "/chatroom", func(c *gin.Context) {

	})
}

func (m *BookingServer) GETChatroom(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("GET", "/chatroom/:room_id", func(c *gin.Context) {

	})
}

func (m *BookingServer) POSTChatrooms(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("POST", "/chatroom", func(c *gin.Context) {

	})
}

func (m *BookingServer) POSTChatroom(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("POST", "/chatroom/:room_id", func(c *gin.Context) {

	})
}

func (m *BookingServer) DELETEChatrooms(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("DELETE", "/chatroom", func(c *gin.Context) {

	})
}

func (m *BookingServer) DELETEChatroom(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("DELETE", "/chatroom/:room_id", func(c *gin.Context) {

	})
}

func (m *BookingServer) GETNonStreamChatrooms(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("GET", "/nonstream/chatroom", func(c *gin.Context) {

	})
}

func (m *BookingServer) GETNonStreamChatroom(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("GET", "/nonstream/chatroom/:room_id", func(c *gin.Context) {

	})
}

func (m *BookingServer) GETLiveChatroom(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("GET", "/live/chatroom/:room_id", func(c *gin.Context) {
		var (
			roomid = c.Param("room_id")
			userid = uuid.NewV4().String()
		)
		c.JSON(
			http.StatusOK,
			gin.H{"response": apis.Resp("Success", map[string]string{
				"room_id": roomid,
				"user_id": userid,
			})},
		)
	})
}

func (m *BookingServer) POSTLiveChatroom(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("POST", "/live/chatroom/:room_id", func(c *gin.Context) {
		var (
			roomid  = c.Param("room_id")
			userid  = c.PostForm("userId")
			message = c.PostForm("message")
			b, ok   = m.RoomChannels[roomid]
		)
		if !ok {
			b = broadcast.NewBroadcaster(m.Booking.RoomChannelsBufferLength)
			m.RoomChannels[roomid] = b
		}
		b.Submit(userid + ": " + message)

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": message,
		})
		c.JSON(
			http.StatusOK,
			gin.H{"response": apis.Resp("Success", map[string]string{
				"room_id": roomid,
				"user_id": userid,
				"message": message,
			})},
		)
	})
}

func (m *BookingServer) DELETELiveChatroom(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("DELETE", "/live/chatroom/:room_id", func(c *gin.Context) {
		var (
			roomid = c.Param("room_id")
			b, ok  = m.RoomChannels[roomid]
		)
		if ok {
			b.Close()
			delete(m.RoomChannels, roomid)
		}
		return
	})
}

func (m *BookingServer) GETLiveStreamChatroom(group *gin.RouterGroup, c *startkit.Context) gin.IRoutes {
	return group.Handle("GET", "/live/stream/chatroom/:room_id", func(c *gin.Context) {
		var (
			roomid   = c.Param("room_id")
			listener = make(chan interface{})
			b, ok    = m.RoomChannels[roomid]
			close    = func(b broadcast.Broadcaster, listener chan interface{}) {
				b.Unregister(listener)
				close(listener)
			}
		)
		if !ok {
			b = broadcast.NewBroadcaster(m.Booking.RoomChannelsBufferLength)
			m.RoomChannels[roomid] = b
		}
		b.Register(listener)
		defer close(b, listener)
		c.Stream(func(w io.Writer) bool {
			c.SSEvent("message", <-listener)
			return true
		})
	})
}
