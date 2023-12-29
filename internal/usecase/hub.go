package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"lemon_be/internal/entity"
	"lemon_be/pkg/redispkg"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	io   sync.Mutex
	Conn *websocket.Conn

	Id       uint
	DeviceId string
	Hub      *Hub

	inbox chan *entity.MsgGeolocationWs
}

type Hub struct {
	mu       sync.RWMutex
	seq      uint
	Rds      *redispkg.Redis
	geoRedis GeoRedisRepo
	us       []*User

	register chan *User

	unregister chan *User
}

func NewHub(rds *redispkg.Redis, geoRedis GeoRedisRepo,
) *Hub {
	return &Hub{
		Rds:        rds,
		geoRedis:   geoRedis,
		unregister: make(chan *User),
		register:   make(chan *User),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case user := <-h.register:
			h.mu.Lock()
			user.Id = h.seq
			//user.DeviceId = user.DeviceId

			h.us = append(h.us, user)
			h.seq++
			h.mu.Unlock()
		case user := <-h.unregister:
			h.mu.Lock()
			// binary search utk cari index user di array us
			i := sort.Search(len(h.us), func(i int) bool {
				return h.us[i].Id >= user.Id
			})

			// hapus client dari array chat.us
			without := make([]*User, len(h.us)-1) // us = nil
			copy(without[:i], h.us[:i])
			copy(without[i:], h.us[i+1:])
			h.us = without
			h.mu.Unlock()
		}
	}
}

var (
	// pongWait : berapa lama server menunggu message pong dari client (30 detik)
	pongWait = 30 * time.Second
	// pingInterval : setiap 5 detik server mengirim ping message ke client.
	// pingInterval haruslah lebih kecil dari pongWait
	pingInterval = (pongWait * 2) / 10
)

// Receive membaca next message websocket dari client
// It blocks until full message received.
func (u *User) Recieve() error {
	defer func() {
		u.Hub.unregister <- u
		u.Conn.Close()
	}()
	// Set Max Size of Messages in Bytes
	u.Conn.SetReadLimit(1024)

	if err := u.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		u.Hub.unregister <- u
		u.Conn.Close()
		return err
	}

	u.Conn.SetPongHandler(u.pongHandler)

	for {
		// ReadMessage dari client websocket
		_, msg, err := u.Conn.ReadMessage()

		if err != nil {
			fmt.Println("errr break loopp:", err)
			fmt.Println("break loop")
			break
		}

		msgWs := &entity.MsgGeolocationWs{}
		if err = json.Unmarshal(msg, msgWs); err != nil {
			log.Println("json.Unmarshal: ", err)
			break
		}

		switch msgWs.Type {
		case entity.MessageTypeUserLocation:
			clientMsg := msgWs.MsgGeolocationUser
			u.Hub.geoRedis.GeoAddVisuallyImpair(
				context.Background(),
				clientMsg.DeviceId,
				clientMsg.Long,
				clientMsg.Lat,
			)
		case entity.MessageTypeCaregiverLocation:
			clientMsg := msgWs.MsgGeolocationCaregiver
			u.Hub.geoRedis.GeoAddCaregiver(context.Background(), clientMsg.TokenFcm,
				clientMsg.Long, clientMsg.Lat)
		}

	}
	return nil
}

// Register registers new connection as a User.
func (h *Hub) Register(ctx context.Context, conn *websocket.Conn, deviceId string) *User {
	user := &User{
		Hub:      h,
		Conn:     conn,
		DeviceId: deviceId,
		inbox:    make(chan *entity.MsgGeolocationWs),
	}
	user.Hub.register <- user

	go user.Recieve()
	go user.writePump()

	return user
}

// pongHandler handle message pong yang dikirim oleh client
// -> mereset durasi readDeadline (tambah 30 detik lagi) & set user online di dalam redis
// dan juga mengirim status online user ke semua kontaknya
func (u *User) pongHandler(pongMsg string) error {

	return u.Conn.SetReadDeadline(time.Now().Add(pongWait))
}

// writePump mengirim message websocket ke user/client/frontend
// 1 goroutine yg menjalankan writePump dijalankan di setiap koneksi client websocket.
func (u *User) writePump() {

	// Create a ticker that triggers a ping at given interval
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
		u.Conn.Close()
	}()
	for {

		select {

		case <-ticker.C:

			err := u.Conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("wsutil.WriteServerMessage", err)
				return
			}

		}
	}
}
