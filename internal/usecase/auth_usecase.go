package usecase

import (
	"context"
	"fmt"
	"lemon_be/internal/entity"
	"lemon_be/internal/util"
	"lemon_be/internal/util/jwt"
	"time"
)

type AuthUseCase struct {
	userRepo      UserRepo
	jwtTokenMaker jwt.JwtTokenMaker
	sessionRepo   SessionRepo
}

func NewAuthUseCase(r UserRepo, j jwt.JwtTokenMaker, s SessionRepo) *AuthUseCase {
	return &AuthUseCase{
		userRepo:      r,
		jwtTokenMaker: j,
		sessionRepo:   s,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, c entity.CreateCaregiverRequest) (entity.Caregiver, error) {
	hashedPassword, err := util.HashPassword(c.Password)
	if err != nil {
		// internal server error
		return entity.Caregiver{}, fmt.Errorf("AuthUseCase - Register -  util.HashPassword: %w", err)
	}
	c.Password = hashedPassword
	createdUser, err := uc.userRepo.CreateUser(ctx, c)
	if err != nil {
		// internal server error/ bad request
		return entity.Caregiver{}, fmt.Errorf("AuthUseCase - Register - uc.userRepo.CreateUser: %w", err)
	}

	return createdUser, nil
}

// Login: logic login use case
func (uc *AuthUseCase) Login(ctx context.Context, l entity.LoginUserRequest) (entity.LoginUserResponse, error) {
	user, err := uc.userRepo.GetUser(ctx, l.Email)
	if err != nil {
		// Bad request User with email not found in DB
		return entity.LoginUserResponse{}, fmt.Errorf("AuthUseCase - Login - uc.userRepo.GetUser: %w", err)
	}

	err = util.CheckPassword(l.Password, user.HashedPassword)
	if err != nil {
		// unauthorized
		return entity.LoginUserResponse{}, fmt.Errorf("AuthUseCase - Login - util.CheckPassword: %w", err)
	}

	accessToken, accessPayload, err := uc.jwtTokenMaker.CreateToken(
		user.Name,
		7*time.Hour,
	)

	if err != nil {
		// internal server error
		return entity.LoginUserResponse{}, fmt.Errorf("AuthUseCase - Login - uc.jwtTokenMaker.CreateToken: %w", err)
	}

	refreshToken, refreshPayload, err := uc.jwtTokenMaker.CreateToken(
		user.Name,
		168*time.Hour,
	)
	if err != nil {
		// internal server error
		return entity.LoginUserResponse{}, fmt.Errorf("AuthUseCase - Login - uc.jwtTokenMaker.CreateToken: %w", err)
	}

	createSessionReq := entity.CreateSessionRequest{
		ID:           refreshPayload.ID.String(),
		Username:     user.Name,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	session, err := uc.sessionRepo.CreateSession(
		ctx,
		createSessionReq,
	)

	if err != nil {
		return entity.LoginUserResponse{}, fmt.Errorf("AuthUseCase - Login - uc.sessionRepo.CreateSession: %w", err)
	}

	userRes := entity.UserResponse{
		Id:       user.Id,
		Username: user.Name,
		Email:    user.Email,
	}

	res := entity.LoginUserResponse{
		SessionId:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  userRes,
	}
	return res, nil
}

// func (uc *AuthUseCase) RenewAccessToken(ctx context.Context, r entity.RenewAccessTokenRequest) (entity.RenewAccessTokenResponse, error) {
// 	refreshPayload, err := uc.jwtTokenMaker.VerifyToken(r.RefreshToken)
// 	if err != nil {
// 		// Unauthorized , token yg dikrim tidak sama dg yg ada di database
// 		return entity.RenewAccessTokenResponse{}, fmt.Errorf("AuthUseCase - RenewAccessToken - uc.jwtTokenMaker.VerifyToken: %w", err)
// 	}

// 	session, err := uc.sessionRepo.GetSession(ctx, refreshPayload.ID)
// 	if err != nil {
// 		return entity.RenewAccessTokenResponse{}, fmt.Errorf("AuthUseCase - RenewAccessToken - uc.sessionRepo.GetSession: %w", err)
// 	}

// 	if session.Username != refreshPayload.Username {
// 		return entity.RenewAccessTokenResponse{}, fmt.Errorf("Invalid session")
// 	}
// 	if session.RefreshToken != r.RefreshToken {
// 		return entity.RenewAccessTokenResponse{}, fmt.Errorf("Invalid session")
// 	}

// 	if time.Now().After(session.ExpiresAt) {
// 		return entity.RenewAccessTokenResponse{}, fmt.Errorf("Invalid session")
// 	}

// 	accessToken, accessTokenPayload, err := uc.jwtTokenMaker.CreateToken(
// 		refreshPayload.Username,
// 		7*time.Hour,
// 	)
// 	if err != nil {
// 		return entity.RenewAccessTokenResponse{}, fmt.Errorf("AuthUseCase - RenewAccessToken - uc.jwtTokenMaker.CreateToken: %w", err)
// 	}

// 	res := entity.RenewAccessTokenResponse{
// 		AccessToken:          accessToken,
// 		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
// 	}
// 	return res, nil
// }

func (uc *AuthUseCase) DeleteRefreshToken(ctx context.Context, d entity.DeleteRefreshTokenRequest) error {

	refreshPayload, err := uc.jwtTokenMaker.VerifyToken(d.RefreshToken)
	if err != nil {
		// Unauthorized , token yg dikrim tidak sama dg yg ada di database
		return fmt.Errorf("AuthUseCase - DeleteRefreshToken - uc.jwtTokenMaker.VerifyToken: %w", err)
	}

	session, err := uc.sessionRepo.GetSession(ctx, refreshPayload.ID.String())
	if err != nil {
		return fmt.Errorf("AuthUseCase - DeleteRefreshToken - uc.sessionRepo.GetSession: %w", err)
	}
	if session.Username != refreshPayload.Username {
		return fmt.Errorf("Invalid session")
	}
	if session.RefreshToken != d.RefreshToken {
		return fmt.Errorf("Invalid session")
	}

	if time.Now().After(session.ExpiresAt) {
		return fmt.Errorf("Invalid session")
	}
	err = uc.sessionRepo.DeleteSession(ctx, refreshPayload.ID.String())
	if err != nil {
		return fmt.Errorf("AuthUseCase - DeleteRefreshToken - uc.sessionRepo.DeleteSession: %w", err)
	}

	return nil

}
