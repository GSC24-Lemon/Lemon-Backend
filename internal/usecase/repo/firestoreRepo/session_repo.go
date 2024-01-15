package firestorerepo

import (
	"context"
	"encoding/json"
	"lemon_be/internal/controller/http/errorWrapper"
	"lemon_be/internal/entity"
	"lemon_be/pkg/firestore"
	"time"

	"google.golang.org/api/iterator"
)

type SessionRepo struct {
	firestore *firestore.Firestore
}

func NewSessionRepo(firestore *firestore.Firestore) *SessionRepo {
	return &SessionRepo{firestore}
}

type Session struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	RefreshToken string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

// CreateSession insert session /refrsh token baru ke database
func (r *SessionRepo) CreateSession(ctx context.Context, c entity.CreateSessionRequest) (entity.Session, error) {

	createSession := Session{
		ID:           c.ID,
		Name:         c.Username,
		RefreshToken: c.RefreshToken,
		ExpiresAt:    c.ExpiresAt,
	}

	_, _, err := r.firestore.Client.Collection("Session").Add(ctx, map[string]interface{}{
		"id":           c.ID,
		"name":         c.Username,
		"refreshToken": c.RefreshToken,
		"expiresAt":    c.ExpiresAt,
	})

	if err != nil {
		// return entity.Session{}, fmt.Errorf("SessionRepo - r.db.Create: %w", err)
		return entity.Session{}, errorWrapper.NewHTTPError(err, 400, "Error when adding new data to Session firestore collection ")

	}

	session := entity.Session{
		ID:           c.ID,
		Username:     createSession.Name,
		RefreshToken: createSession.RefreshToken,
		ExpiresAt:    createSession.ExpiresAt,
	}
	return session, nil
}

func (r *SessionRepo) GetSession(ctx context.Context, refreshTokkenId string) (entity.Session, error) {
	var sessionDb Session

	// result := r.db.Where(&Session{ID: refreshTokkenId}).First(&sessionDb)
	iter := r.firestore.Client.Collection("Session").Where("id", "==", refreshTokkenId).Documents(ctx)
	var objectMap map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// return entity.Session{}, fmt.Errorf("SessionRepo - GetSession- r.db.Where(&Session{ID: refreshTokkenId}).First(&sessionDb): %w", err)
			return entity.Session{}, errorWrapper.NewHTTPError(err, 401, "Session not found in database")

		}

		objectMap = doc.Data()
	}

	jsonObjectMap, _ := json.Marshal(objectMap)

	json.Unmarshal(jsonObjectMap, &sessionDb)

	session := entity.Session{
		ID:           sessionDb.ID,
		Username:     sessionDb.Name,
		RefreshToken: sessionDb.RefreshToken,
		ExpiresAt:    sessionDb.ExpiresAt,
		CreatedAt:    sessionDb.CreatedAt,
	}

	return session, nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, uuid string) error {
	var sessionDb Session

	// result := r.db.Where(&Session{ID: refreshTokkenId}).First(&sessionDb)
	iter := r.firestore.Client.Collection("Session").Where("id", "==", uuid).Documents(ctx)
	var objectMap map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// return fmt.Errorf("SessionRepo - GetSession- r.db.Where(&Session{ID: refreshTokkenId}).First(&sessionDb): %w", err)
			return errorWrapper.NewHTTPError(err, 404, "Session not found in database")
		}

		objectMap = doc.Data()
	}

	jsonObjectMap, _ := json.Marshal(objectMap)

	json.Unmarshal(jsonObjectMap, &sessionDb)

	session := entity.Session{
		ID:           sessionDb.ID,
		Username:     sessionDb.Name,
		RefreshToken: sessionDb.RefreshToken,
		ExpiresAt:    sessionDb.ExpiresAt,
		CreatedAt:    sessionDb.CreatedAt,
	}

	// result := r.firestore.Client.Where(&Session{ID: uuid}).Delete(&sessionDb)

	_, err := r.firestore.Client.Collection("Session").Doc(session.ID).Delete(ctx)
	if err != nil {
		// return fmt.Errorf("Sessionrepo - r.db.Where(&Session{ID: uuid}).Delete(&sessionDb) - %w", err)
		return errorWrapper.NewHTTPError(err, 401, "Session not found in database")
	}

	return nil
}
