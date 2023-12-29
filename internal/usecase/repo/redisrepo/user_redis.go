package redisrepo

import (
	"context"
	"lemon_be/pkg/redispkg"
)

type UserRedisRepo struct {
	rds *redispkg.Redis
}

func NewUserRedisRepo(rds *redispkg.Redis) *UserRedisRepo {
	return &UserRedisRepo{rds}
}

func (r *UserRedisRepo) SaveUsernameAndDeviceId(ctx context.Context, deviceId string, username string, telephone string) {

	r.rds.Client.Set(ctx, deviceId, username, 0)
	r.rds.Client.Set(ctx, username+":telephone", telephone, 0)
	return

}

func (r *UserRedisRepo) GetUsernameFromDeviceId(ctx context.Context, deviceId string) ([]string, error) {
	res, err := r.rds.Client.Get(ctx, deviceId).Result()
	if err != nil {
		return nil, nil
	}
	username := res

	resTelephone, err := r.rds.Client.Get(ctx, username+":telephone").Result()
	if err != nil {
		return nil, nil
	}
	telephone := resTelephone
	var ans []string
	ans = append(ans, username)
	ans = append(ans, telephone)

	return ans, nil
}
