package repository

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	model "github.com/FurmanovD/cloudtests/internal/pkg/model/userservice"
)

// GetUser returns a user by its ID.
func (repo *repository) GetUser(ctx context.Context, userID string) (model.User, error) {
	log.Trace("repository.GetUser(%v) called", userID)

	val, err := repo.client.Get(ctx, userID).Result()
	if err != nil {
		return model.User{}, err
	}
	user := model.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
