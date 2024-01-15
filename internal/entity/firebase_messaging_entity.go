package entity

type SpecificHelpNotificationRequest struct {
	Token       string `json:"token"`
	Username    string `json:"username"`
	Address     string `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Telephone   string `json:"telephone"`
	Destination string `json:"destination"`
}

type NotificationMessagingPayload struct {
	Data         map[string]string `json:"data"`
	Token        string            `json:"token"`
	Notification map[string]string `json:"notification"`
}

type NotificationBody struct {
	Message NotificationMessagingPayload `json:"message"`
}
