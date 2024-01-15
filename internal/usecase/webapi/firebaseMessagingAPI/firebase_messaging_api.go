package firebasemessagingapi

import (
	"context"
	"fmt"
	"lemon_be/internal/entity"
	"lemon_be/pkg/fcm"

	"firebase.google.com/go/messaging"
)

type FirebaseMessagingAPI struct {
	fbMessaging *fcm.FirebaseMessaging
}

func NewFirebaseMessagingAPI(fbMessaging *fcm.FirebaseMessaging) *FirebaseMessagingAPI {
	return &FirebaseMessagingAPI{fbMessaging}
}

func (api *FirebaseMessagingAPI) SendNotifToSpecificDevice(ctx context.Context, np entity.SpecificHelpNotificationRequest) error {
	registrationToken := np.Token
	message := &messaging.Message{
		Data: map[string]string{
			"uLatitude":   fmt.Sprintf("%f", np.Latitude),
			"uLongitude":  fmt.Sprintf("%f", np.Longitude),
			"username":    np.Username,
			"telephone":   np.Telephone,
			"destination": np.Destination,
			"type":        "help_request",
		},
		Token: registrationToken,
		Notification: &messaging.Notification{
			Body:  np.Username + " needs your help right now. are you willing to help him?  His location is on " + np.Address,
			Title: np.Username + " needs your help right now",
		},
	}

	_, err := api.fbMessaging.Client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("FirebaseMessagingAPI - api.fbMessaging.Client.Send: %w", err)
	}

	return nil
}
