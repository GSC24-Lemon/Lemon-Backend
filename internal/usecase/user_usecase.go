package usecase

import (
	"lemon_be/internal/entity"

	"golang.org/x/net/context"
)

type UserUseCase struct {
	userRedisRepo UserRedisRepo
}

func NewUserUseCase(userRedisRepo UserRedisRepo) *UserUseCase {
	return &UserUseCase{userRedisRepo}
}

func (uc *UserUseCase) SaveUsernameAndDeviceId(ctx context.Context, e entity.SaveUsername) error {
	err := uc.userRedisRepo.SaveUsernameAndDeviceId(ctx, e.DeviceId, e.Username, e.Telephone)
	if err != nil {
		return err
	}

	return nil
}
