package usecase

import (
	"context"
	"lemon_be/internal/entity"
	"net/http"
)

type (
	Caregiver interface {
		NotifyNearestCaregiver(ctx context.Context, e entity.UserLocation)
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
		GeoAddVisuallyImpair(deviceId string, long float64, lat float64)
		Geohash(key string) (string, error)
		GeoAddCaregiver(tokenFcm string, long float64, lat float64)
		GetCaregiverTokens(areaGeohash []string) ([]string, error)
	}

	// Websocket usecase
	Websocket interface {
		WebsocketHandler(http.ResponseWriter, *http.Request, context.Context) error
	}

	UserRedisRepo interface {
		SaveUsernameAndDeviceId(deviceId string, username string)
		GetUsernameFromDeviceId(deviceId string) string
	}

	UserUseCaseI interface {
		SaveUsernameAndDeviceId(ctx context.Context, e entity.SaveUsername)
	}
)
