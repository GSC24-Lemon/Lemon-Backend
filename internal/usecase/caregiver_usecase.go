package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mmcloughlin/geohash"
	"golang.org/x/net/context"
	"io/ioutil"
	"lemon_be/internal/entity"
	"net/http"
	"os"
)

type CaregiverUseCase struct {
	repo        GeoRedisRepo
	userRdsRepo UserRedisRepo
}

func NewCaregiverUseCase(r GeoRedisRepo, userRdsRepo UserRedisRepo) *CaregiverUseCase {
	return &CaregiverUseCase{
		repo:        r,
		userRdsRepo: userRdsRepo,
	}
}

func (uc *CaregiverUseCase) NotifyNearestCaregiver(ctx context.Context, e entity.UserLocation) {
	userGeohash, err := uc.repo.Geohash(e.DeviceId)
	userGeohash = userGeohash[0:6]
	if err != nil {
		return
	}
	userGeohashNeighbors := geohash.Neighbors(userGeohash)

	nearestGeohash := userGeohashNeighbors
	nearestGeohash = append(nearestGeohash, userGeohash)

	caregiverTokenFcms, err := uc.repo.GetCaregiverTokens(nearestGeohash)
	if err != nil {
		return
	}

	pushNotificationToCaregivers(caregiverTokenFcms, e.Long, e.Lat, uc.userRdsRepo.GetUsernameFromDeviceId(e.DeviceId))

}

func pushNotificationToCaregivers(tokenFcms []string, longitude float64, latitude float64, username string) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s", latitude, longitude, os.Getenv("GEOCODING_API_KEY"))
	client, err := http.DefaultClient.Get(url)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(client.Body)

	var jsonData map[string][]map[string]interface{}

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return
	}

	address := jsonData["results"][2]["formatted_address"]

	postUrl := "https://fcm.googleapis.com/v1/projects/lemon-df113/messages:send"

	for token := range tokenFcms {
		notifBody := []byte(fmt.Sprintf(`
			{
				"message": {
					"token": %s,
					"notification": {
						"body": "%s needs your help right now. are you willing to help him?  His location is on %s",
						"title": "%s needs your help right now"
					},
					"data": {
						"uLatitude": %f,
						"uLongitude": %f,
						"username": %s
					}
				}
			}
			`, token, username, address, username, latitude, longitude, username))
		r, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(notifBody))
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("Authorization", "Bearer "+os.Getenv("FCM_KEY"))

		client := &http.Client{}
		_, err = client.Do(r)
		if err != nil {
			return
		}
	}
}
