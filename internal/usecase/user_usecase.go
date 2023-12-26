package usecase

import (
	"golang.org/x/net/context"
	"lemon_be/internal/entity"
	"lemon_be/internal/usecase/repo/redisrepo"
)

type UserUseCase struct {
	userRedisRepo *redisrepo.UserRedisRepo
}

func NewUserUseCase(userRedisRepo *redisrepo.UserRedisRepo) *UserUseCase {
	return &UserUseCase{userRedisRepo}
}

func (uc *UserUseCase) SaveUsernameAndDeviceId(ctx context.Context, e entity.SaveUsername) {
	uc.userRedisRepo.SaveUsernameAndDeviceId(e.DeviceId, e.Username)

}
