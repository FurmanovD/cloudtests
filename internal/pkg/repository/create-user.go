package repository

import (
	"context"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"

	"github.com/FurmanovD/cloudtests/internal/pkg/jsontime"
	model "github.com/FurmanovD/cloudtests/internal/pkg/model/userservice"
)

// CreateUser creates a new user with a generated ID that is returned to the caller.
func (repo *repository) CreateUser(ctx context.Context, user model.User) (string, error) {
	log.Trace("repository.CreateUser(%v) called", user)

	key := tableKeyPrefixUser + uuid.New().String()
	user.ID = key
	user.CreatedAt = jsontime.Now()

	err := repo.setValueAsJSON(ctx, key, user, defDuration200Years)
	if err != nil {
		return "", err
	}

	return key, nil
}
