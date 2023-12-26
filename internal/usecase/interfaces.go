package usecase

import (
	"context"
	"lemon_be/internal/entity"
)

type (
	Caregiver interface {
		getNearestCaregiver(context.Context)
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
)
