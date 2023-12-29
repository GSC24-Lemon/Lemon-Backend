package usecase

import (
	"fmt"
	"lemon_be/internal/entity"

	"golang.org/x/net/context"
)

type UserUseCase struct {
	userRedisRepo UserRedisRepo
}

func NewUserUseCase(userRedisRepo UserRedisRepo) *UserUseCase {
	return &UserUseCase{userRedisRepo}
}

func (uc *UserUseCase) SaveUsernameAndDeviceId(ctx context.Context, e entity.SaveUsername) {
	uc.userRedisRepo.SaveUsernameAndDeviceId(ctx, e.DeviceId, e.Username, e.Telephone)
	fmt.Printf("registering user: %s %s %s", e.DeviceId, e.Username, e.Telephone)
}
