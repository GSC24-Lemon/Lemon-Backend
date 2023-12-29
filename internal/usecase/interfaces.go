package usecase

import (
	"context"
	"lemon_be/internal/entity"
	"net/http"

	"github.com/gorilla/websocket"
)

type (
	Caregiver interface {
		NotifyNearestCaregiver(context.Context, entity.UserLocation)
		TestGeoAdd(context.Context, entity.UserLocation)
	}

	UserRepo interface {
		CreateUser(context.Context, entity.CreateCaregiverRequest) (entity.Caregiver, error)
		GetUser(context.Context, string) (entity.Caregiver, error)
	}

	// Auth Use Case
	Auth interface {
		Register(ctx context.Context, c entity.CreateCaregiverRequest) (entity.Caregiver, error)
		Login(context.Context, entity.LoginUserRequest) (entity.LoginUserResponse, error)
		//RenewAccessToken(context.Context, entity.RenewAccessTokenRequest) (entity.RenewAccessTokenResponse, error)
		DeleteRefreshToken(context.Context, entity.DeleteRefreshTokenRequest) error
	}

	// SessionRepo
	SessionRepo interface {
		CreateSession(ctx context.Context, c entity.CreateSessionRequest) (entity.Session, error)
		GetSession(context.Context, string) (entity.Session, error)
		DeleteSession(context.Context, string) error
	}

	GeoRedisRepo interface {
		GeoAddVisuallyImpair(context.Context, string, float64, float64)
		Geohash(context.Context, string) (string, error)
		GeoAddCaregiver(context.Context, string, float64, float64)
		GetCaregiverTokens(context.Context, []string) ([]string, error)
	}

	// Websocket usecase
	Websocket interface {
		WebsocketHandler(http.ResponseWriter, *http.Request, context.Context) error
	}

	UserRedisRepo interface {
		SaveUsernameAndDeviceId(context.Context, string, string, string)
		GetUsernameFromDeviceId(context.Context, string) ([]string, error)
	}

	UserUseCaseI interface {
		SaveUsernameAndDeviceId(context.Context, entity.SaveUsername)
	}

	ChatHubI interface {
		Register(context.Context, *websocket.Conn, string) *User
		Run()
	}

	HelpRepo interface {
		InsertHelp(context.Context, entity.UserLocation, string) error
	}
)
