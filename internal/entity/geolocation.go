package entity

type MsgGeolocationWs struct {
	Type                    MessageType          `json:"type"`
	MsgGeolocationUser      GeoLocationUser      `json:"msg_geolocation_user,omitempty"`
	MsgGeolocationCaregiver GeoLocationCaregiver `json:"msg_geolocation_caregiver,omitempty"`
}

type GeoLocationUser struct {
	DeviceId string
	Long     float64
	Lat      float64
}

type GeoLocationCaregiver struct {
	TokenFcm string
	Long     float64
	Lat      float64
}

type (
	MessageType string
)

const (
	MessageTypeUserLocation      MessageType = "user_location"
	MessageTypeCaregiverLocation MessageType = "caregiver_location"
)
