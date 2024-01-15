package redisrepo

import (
	"context"
	"fmt"
	"lemon_be/internal/controller/http/errorWrapper"
	"lemon_be/pkg/redispkg"
)

type UserRedisRepo struct {
	rds *redispkg.Redis
}

func NewUserRedisRepo(rds *redispkg.Redis) *UserRedisRepo {
	return &UserRedisRepo{rds}
}

func (r *UserRedisRepo) SaveUsernameAndDeviceId(ctx context.Context, deviceId string, username string, telephone string) error {

	err := r.rds.Client.Set(ctx, deviceId, username, 0).Err()
	if err != nil {
		return errorWrapper.NewHTTPError(err, 400, fmt.Sprintf("cannot add hash with key and value: %s, %s", deviceId, username))
	}
	r.rds.Client.Set(ctx, username+":telephone", telephone, 0)
	if err != nil {
		return errorWrapper.NewHTTPError(err, 400, fmt.Sprintf("cannot add hash with key and value: %s, %s", deviceId, telephone))
	}
	return nil

}

func (r *UserRedisRepo) GetUsernameFromDeviceId(ctx context.Context, deviceId string) ([]string, error) {
	res, err := r.rds.Client.Get(ctx, deviceId).Result()
	if err != nil {
		// return nil, fmt.Errorf("UserRedisRepo -  GetUsernameFromDeviceId - r.rds.Client.Get(ctx, deviceId).Result(): %w", err)
		return nil, errorWrapper.NewHTTPError(err, 404, "username not found for deviceId: "+deviceId)
	}
	username := res

	resTelephone, err := r.rds.Client.Get(ctx, username+":telephone").Result()
	if err != nil {
		// return nil, fmt.Errorf("UserRedisRepo -  GetUsernameFromDeviceId -  r.rds.Client.Get(ctx, username).Result(): %w", err)
		return nil, errorWrapper.NewHTTPError(err, 404, "telepohone not found for deviceId: "+deviceId)
	}
	telephone := resTelephone
	var ans []string
	ans = append(ans, username)
	ans = append(ans, telephone)

	return ans, nil
}
