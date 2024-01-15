package redisrepo

import (
	"context"
	"fmt"
	"lemon_be/internal/controller/http/errorWrapper"
	"lemon_be/pkg/redispkg"
	"math"

	"github.com/redis/go-redis/v9"
)

type GeoRedisRepo struct {
	rds *redispkg.Redis
}

func NewGeoRedisRepo(rds *redispkg.Redis) *GeoRedisRepo {
	return &GeoRedisRepo{rds}
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

const geoKey = "currentLocation"

// this function called every 1 second to save curent user(visually impaired) and caregiver location
// userkey = device id
// caregiverKey= tokenFcm
func (r *GeoRedisRepo) GeoAddVisuallyImpair(ctx context.Context, deviceId string, long float64,
	lat float64) {
	res := r.rds.Client.GeoPos(ctx, geoKey, deviceId).Val()
	err := r.rds.Client.GeoPos(ctx, geoKey, deviceId).Err()
	if err != nil {
		fmt.Println("errorr geopos!!: ", err)
	}
	if err == nil && res[0] != nil {
		// geoset with deviceId geoKey, name curr already exists
		loc := res[0]
		err = r.rds.Client.ZRem(ctx, geoKey, deviceId).Err()
		if err != nil {
			fmt.Println("error delete coord: ", err)
		}

		if loc.Latitude == lat && loc.Longitude == long {
			fmt.Println("same coordinate just like  prev coord: ", err)
			return
		}

	}

	fmt.Println("deviceId geoadd: ", deviceId, " long: ", long, " lat: ", lat)
	err = r.rds.Client.GeoAdd(ctx, geoKey,
		&redis.GeoLocation{
			Longitude: long,
			Latitude:  lat,
			Name:      deviceId,
		}).Err()

	if err != nil {
		fmt.Println("errorr geoadd!!: ", err)
	}

	return
}

// return the geohash from the coordinate
func (r *GeoRedisRepo) Geohash(ctx context.Context, key string) (string, error) {
	geohash, err := r.rds.Client.GeoHash(ctx, geoKey, key).Result()
	if err != nil {
		// return "", fmt.Errorf("BadRequest - GeoRedisRepo - Geohash : %w", err)
		return "", errorWrapper.NewHTTPError(err, 400, "Cannot geohash key: "+key)
	}

	return geohash[0], nil
}

func (r *GeoRedisRepo) GeoAddCaregiver(ctx context.Context, tokenFcm string, long float64,
	lat float64) {
	res, err := r.rds.Client.GeoPos(ctx, geoKey, tokenFcm).Result()
	if err == nil && res[0] != nil {
		// geoset with key tokenFcm, name curr already exists
		err := r.rds.Client.ZRem(ctx, geoKey, tokenFcm).Err()
		if err != nil {
			fmt.Println("error delete coord: ", err)
		}
		loc := res[0]
		if loc.Latitude == lat && loc.Longitude == long {
			return
		}

	}

	r.rds.Client.GeoAdd(ctx, geoKey, &redis.GeoLocation{
		Name:      tokenFcm,
		Longitude: long,
		Latitude:  lat,
	})

	geohash, _ := r.Geohash(ctx, tokenFcm)
	geohash = geohash[0:7]
	r.rds.Client.SAdd(ctx, geohash, tokenFcm)

}

func (r *GeoRedisRepo) GetCaregiverTokens(ctx context.Context, areaGeohash []string) ([]string, error) {
	var caregiverTokens []string
	for _, currGeohash := range areaGeohash {
		currTokenFcms, err := r.rds.Client.SMembers(ctx, currGeohash).Result()
		if err != nil {
			// there are no set members
			continue
		}

		for _, currTokenFcm := range currTokenFcms {
			caregiverTokens = append(caregiverTokens, currTokenFcm)
		}

	}
	if len(caregiverTokens) == 0 {
		// return nil, fmt.Errorf("no caregiver found in this area!! %w", errors.New("no caregiver found in this area"));
		return nil, errorWrapper.NewHTTPError(nil, 404, "no caregiver found in this area!!")
	}
	return caregiverTokens, nil

}
