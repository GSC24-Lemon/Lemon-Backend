package firestorerepo

import (
	"context"
	"encoding/json"
	"fmt"
	"lemon_be/internal/entity"
	"lemon_be/pkg/firestore"

	"google.golang.org/api/iterator"
)

type UserRepo struct {
	firestore *firestore.Firestore
}

func NewUserRepo(client *firestore.Firestore) *UserRepo {
	return &UserRepo{client}
}

func (r *UserRepo) CreateUser(ctx context.Context, e entity.CreateCaregiverRequest) (entity.Caregiver, error) {
	iter := r.firestore.Client.Collection("User").Where("email", "==", e.Email).Documents(ctx)
	var objMap map[string]interface{}
	var userId string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}
		objMap = doc.Data()
		userId = doc.Ref.ID
	}
	if objMap != nil && userId != "" {
		return entity.Caregiver{}, fmt.Errorf("BadRequest - UserRepo - CreateUser - r.firestore.Client.Collection ")
	}

	_, _, err := r.firestore.Client.Collection("User").Add(ctx, map[string]interface{}{
		"name":     e.Name,
		"password": e.Password,
		"email":    e.Email,
		"gender":   e.Gender,
		"age":      int(e.Age),
	})

	if err != nil {
		return entity.Caregiver{}, fmt.Errorf("UserRepo - CreateUser - r.firestore.Client.Collection - %w", err)
	}

	caregiver := entity.Caregiver{
		Name:   e.Name,
		Email:  e.Email,
		Age:    e.Age,
		Gender: e.Gender,
		Job:    e.Job,
	}

	return caregiver, nil
}

func (r *UserRepo) GetUser(ctx context.Context, e string) (entity.Caregiver, error) {
	var userDb entity.Caregiver
	// result := r.db.Where(&User{Email: email}).First(&userDb)
	iter := r.firestore.Client.Collection("User").Where("email", "==", e).Documents(ctx)
	var objMap map[string]interface{}
	var userId string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return entity.Caregiver{}, fmt.Errorf("UserRepo - GetUser - Collection().Where %w", err)
		}
		objMap = doc.Data()
		userId = doc.Ref.ID
	}

	jsonObjectMap, _ := json.Marshal(objMap)

	json.Unmarshal(jsonObjectMap, &userDb)

	user := entity.Caregiver{
		Id:             userId,
		Name:           userDb.Name,
		Email:          userDb.Email,
		Age:            userDb.Age,
		Gender:         userDb.Gender,
		HashedPassword: userDb.HashedPassword,
		Job:            userDb.Job,
	}

	return user, nil
}
