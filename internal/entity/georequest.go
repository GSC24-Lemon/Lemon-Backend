package entity

type UserLocation struct {
	DeviceId string  `json:"deviceId"`
	Long     float64 `json:"longitude"`
	Lat      float64 `json:"latitude"`
	Username string `json:"username"`
	Destination string `json:"destination"`
}
