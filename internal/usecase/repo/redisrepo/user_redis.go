package redisrepo

import (
	"golang.org/x/net/context"
	"lemon_be/pkg/redispkg"
)

type UserRedisRepo struct {
	rds *redispkg.Redis
}

func NewUserRedisRepo(rds *redispkg.Redis) *UserRedisRepo {
	return &UserRedisRepo{rds}
}

func (r *UserRedisRepo) SaveUsernameAndDeviceId(deviceId string, username string) {
	r.rds.Client.Set(context.Background(), deviceId, username, 0)
	return

}

func (r *UserRedisRepo) GetUsernameFromDeviceId(deviceId string) string {
	res, _ := r.rds.Client.SMembers(context.Background(), deviceId).Result()
	username := res[0]
	return username
}
