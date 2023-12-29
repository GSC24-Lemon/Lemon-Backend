package usecase

import (
	"context"
	"errors"
	"net/http"

	websocket "github.com/gorilla/websocket"
)

var (
	WebsocketConnectionError   = errors.New("websocket connection error")
	WebsocketUnauthorizedError = errors.New("websocket unauthorized error")
)

// WebsocketUseCase bussines logic websocketc
type WebsocketUseCase struct {
	userPg       UserRepo
	geoRedisRepo GeoRedisRepo
	hub          ChatHubI
}

// NewWebsocket Create new websocketUseCase
func NewWebsocketUseCase(hub ChatHubI, userPg UserRepo, geoRedisRepo GeoRedisRepo) *WebsocketUseCase {
	return &WebsocketUseCase{userPg, geoRedisRepo, hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (uc *WebsocketUseCase) WebsocketHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) error {

	deviceId := r.URL.Query().Get("deviceId")
	if deviceId == "" {
		// Tell the user its not authorized
		return WebsocketUnauthorizedError
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return WebsocketConnectionError
	}

	_ = uc.hub.Register(ctx, conn, deviceId)

	return nil
}
