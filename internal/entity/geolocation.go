package entity

type MsgGeolocationWs struct {
	Type                    MessageType          `json:"type"`
	MsgGeolocationUser      GeoLocationUser      `json:"msg_geolocation_user,omitempty"`
	MsgGeolocationCaregiver GeoLocationCaregiver `json:"msg_geolocation_caregiver,omitempty"`
}

type GeoLocationUser struct {
	DeviceId string  `json:"deviceId"`
	Long     float64 `json:"longitude"`
	Lat      float64 `json:"latitude"`
}

type GeoLocationCaregiver struct {
	TokenFcm string  `json:"tokenFcm"`
	Long     float64 `json:"longitude"`
	Lat      float64 `json:"latitude"`
}

type (
	MessageType string 
)

const (
	MessageTypeUserLocation      MessageType = "user_location"
	MessageTypeCaregiverLocation MessageType = "caregiver_location"
)
