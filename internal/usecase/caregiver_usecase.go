package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lemon_be/internal/entity"
	"net/http"
	"os"

	"github.com/mmcloughlin/geohash"
	"golang.org/x/net/context"
)

type CaregiverUseCase struct {
	repo        GeoRedisRepo
	userRdsRepo UserRedisRepo
	helpRepo    HelpRepo
	fcmAPI      FirebaseMessaging
}

func NewCaregiverUseCase(r GeoRedisRepo,
	userRdsRepo UserRedisRepo,
	helpRepo HelpRepo,
	fcm FirebaseMessaging) *CaregiverUseCase {
	return &CaregiverUseCase{
		repo:        r,
		userRdsRepo: userRdsRepo,
		helpRepo:    helpRepo,
		fcmAPI:      fcm,
	}
}

func (uc *CaregiverUseCase) NotifyNearestCaregiver(ctx context.Context, e entity.UserLocation) error {
	userGeohash, err := uc.repo.Geohash(ctx, e.DeviceId)
	userGeohash = userGeohash[0:7]
	if err != nil {
		// return fmt.Errorf("CaregiverUseCase - NotifyNearestCaregiver -  uc.repo.Geohash(ctx, e.DeviceId): %w", err)
		return err
	}
	userGeohashNeighbors := geohash.Neighbors(userGeohash)

	nearestGeohash := userGeohashNeighbors
	nearestGeohash = append(nearestGeohash, userGeohash)

	caregiverTokenFcms, err := uc.repo.GetCaregiverTokens(ctx, nearestGeohash)
	if err != nil {
		// return fmt.Errorf("CaregiverUseCase - NotifyNearestCaregiver - uc.repo.GetCaregiverTokens(ctx, nearestGeohash): %w", err)
		return err

	}
	res, err := uc.userRdsRepo.GetUsernameFromDeviceId(ctx, e.DeviceId)
	if err != nil {
		// return fmt.Errorf("CaregiverUseCase - NotifyNearestCaregiver - uc.userRdsRepo.GetUsernameFromDeviceId(ctx, e.DeviceId): %w", err)
		return err
	}
	username := res[0]
	telephone := res[1]

	uc.helpRepo.InsertHelp(ctx, e, userGeohash)
	err = uc.PushNotificationToCaregivers(caregiverTokenFcms, e.Long, e.Lat, username,
		telephone, e.Destination)
	if err != nil {
		return fmt.Errorf("CaregiverUseCase - NotifyNearestCaregiver - uc.PushNotificationToCaregivers(caregiverTokenFcms, e.Long, e.Lat, username telephone, e.Destination):")

	}

	return nil
}

func (uc *CaregiverUseCase) TestGeoAdd(ctx context.Context, e entity.UserLocation) {
	uc.repo.GeoAddVisuallyImpair(ctx, e.DeviceId, e.Long, e.Lat)

	//res, err := uc.userRdsRepo.GetUsernameFromDeviceId(ctx, e.DeviceId)
	//if err != nil {
	//	return
	//}
	//username := res[0]
	//telephone := res[1]
	//
	//fmt.Println("username : " + username + " telephone: " + telephone)

}

func (uc *CaregiverUseCase) PushNotificationToCaregivers(tokenFcms []string, longitude float64, latitude float64,
	username string, telephone string, destination string) error {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s", latitude, longitude, os.Getenv("GEOCODING_API_KEY"))
	client, err := http.DefaultClient.Get(url)
	if err != nil {
		return fmt.Errorf("CaregiverUseCase - PushNotificationToCaregivers - http.DefaultClient.Get: ")
		// return errorWrapper.NewHTTPError(err, 500, "cannot do GET request to url https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s")
	}
	data, err := ioutil.ReadAll(client.Body)

	//var jsonData map[string][]map[string]interface{}
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		return fmt.Errorf("CaregiverUseCase - PushNotificationToCaregivers - json.Unmarshal(data, &jsonData): ")

	}
	results := jsonData["results"].([]interface{})

	addressMap := results[2].(map[string]interface{})
	address := addressMap["formatted_address"]

	// postUrl := "https://fcm.googleapis.com/v1/projects/lemon-df113/messages:send"

	for _, token := range tokenFcms {

		err := uc.fcmAPI.SendNotifToSpecificDevice(context.Background(), entity.SpecificHelpNotificationRequest{
			Token:       token,
			Username:    username,
			Address:     fmt.Sprintf("%v", address),
			Latitude:    latitude,
			Longitude:   longitude,
			Telephone:   telephone,
			Destination: destination,
		})

		if err != nil {
			return fmt.Errorf("CaregiverUseCase - PushNotificationToCaregivers - uc.fcmAPI.SendNotifToSpecificDevice: ")

		}

		// notifBody := []byte(fmt.Sprintf(`
		// 	{
		// 		"message": {
		// 			"token": %s,
		// 			"notification": {
		// 				"body": "%s needs your help right now. are you willing to help him?  His location is on %s",
		// 				"title": "%s needs your help right now"
		// 			},
		// 			"data": {
		// 				"uLatitude": %f,
		// 				"uLongitude": %f,
		// 				"username": %s,
		// 				"telephone": %s,
		// 				"destination": %s
		// 			}
		// 		}
		// 	}
		// 	`, token, username, address, username, latitude, longitude, username, telephone, destination))
		// r, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(notifBody))
		// if err != nil {
		// 	return
		// }
		// r.Header.Add("Content-Type", "application/json")
		// r.Header.Add("Authorization", "Bearer "+os.Getenv("FCM_KEY"))

		// client := &http.Client{}
		// _, err = client.Do(r)
		// if err != nil {
		// 	fmt.Errorf("error request: %w", err)
		// 	return
		// }
	}
	return nil
}
