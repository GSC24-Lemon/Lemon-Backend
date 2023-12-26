package redisrepo

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"lemon_be/pkg/redispkg"
)

type GeoRedisRepo struct {
	rds *redispkg.Redis
}

func NewGeoRedisRepo(rds *redispkg.Redis) *GeoRedisRepo {
	return &GeoRedisRepo{rds}
}

// this function called every 1 second to save curent user(visually impaired) and caregiver location
// userkey = device id
// caregiverKey= tokenFcm
func (r *GeoRedisRepo) GeoAddVisuallyImpair(deviceId string, long float64,
	lat float64) {

	_, err := r.rds.Client.GeoPos(context.Background(), deviceId, "currentLocation").Result()
	if err == nil {
		// geoset with key deviceId, name currentLocation already exists
		return
	}
	r.rds.Client.GeoAdd(context.Background(), deviceId, &redis.GeoLocation{
		Name:      "currentLocation",
		Longitude: long,
		Latitude:  lat,
	})
	return
}

// return the geohash from the coordinate
func (r *GeoRedisRepo) Geohash(key string) (string, error) {
	geohash, err := r.rds.Client.GeoHash(context.Background(), key, "currentLocation").Result()
	if err != nil {
		return "", fmt.Errorf("BadRequest - GeoRedisRepo - Geohash : %w", err)
	}

	return geohash[0], nil
}

func (r *GeoRedisRepo) GeoAddCaregiver(tokenFcm string, long float64,
	lat float64) {
	_, err := r.rds.Client.GeoPos(context.Background(), tokenFcm, "currentLocation").Result()
	if err == nil {
		// geoset with key deviceId, name currentLocation already exists
		return
	}
	r.rds.Client.GeoAdd(context.Background(), tokenFcm, &redis.GeoLocation{
		Name:      "currentLocation",
		Longitude: long,
		Latitude:  lat,
	})

	geohash, _ := r.Geohash(tokenFcm)
	r.rds.Client.SAdd(context.Background(), geohash, tokenFcm)

}

func (r *GeoRedisRepo) GetCaregiverTokens(areaGeohash []string) ([]string, error) {
	var caregiverTokens []string
	for _, currGeohash := range areaGeohash {
		currTokenFcms, err := r.rds.Client.SMembers(context.Background(), currGeohash).Result()
		if err != nil {
			// there are no set members
			continue
		}

		for _, currTokenFcm := range currTokenFcms {
			caregiverTokens = append(caregiverTokens, currTokenFcm)
		}

	}
	return caregiverTokens, nil

}
