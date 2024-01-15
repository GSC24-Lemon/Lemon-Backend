package firestorerepo

import (
	"context"
	"lemon_be/internal/controller/http/errorWrapper"
	"lemon_be/internal/entity"
	"lemon_be/pkg/firestore"
)

type HelpRepo struct {
	firestore *firestore.Firestore
}

func NewHelpRepo(client *firestore.Firestore) *HelpRepo {
	return &HelpRepo{client}
}

func (r *HelpRepo) InsertHelp(ctx context.Context, e entity.UserLocation, uGeohash string) error {
	_, _, err := r.firestore.Client.Collection("NeedHelp").Add(ctx, map[string]interface{}{
		"name":        e.Username,
		"deviceId":    e.DeviceId,
		"latitude":    e.Lat,
		"longitude":   e.Long,
		"destination": e.Destination,
		"userGeohash": uGeohash,
	})

	if err != nil {
		// return  fmt.Errorf("HelpRepo - r.db.Collection.Add: %w", err)
		return errorWrapper.NewHTTPError(err, 400, "Error when adding new data to NeedHelp firestore collection ")
	}

	return nil

}
