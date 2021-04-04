package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/FurmanovD/cloudtests/internal/pkg/jsontime"
	model "github.com/FurmanovD/cloudtests/internal/pkg/model/userservice"
	"github.com/FurmanovD/cloudtests/internal/pkg/redislock"
)

// UpdateUser updates a user record by its ID.
func (repo *repository) UpdateUser(ctx context.Context, user model.User, optimisticLock bool) error {
	log.Trace("repository.UpdateUser(%v, %v) called", user, optimisticLock)

	if !strings.HasPrefix(user.ID, tableKeyPrefixUser) { // also covers an empty case
		return redis.Nil // not found
	}

	if optimisticLock {

		lastUpdatedAt := user.UpdatedAt
		updated := false
		for !updated {

			currDBUser, err := repo.GetUser(ctx, user.ID)
			if err != nil {
				return err
			}

			if lastUpdatedAt == currDBUser.UpdatedAt {
				// patch CreatedAt
				user.CreatedAt = currDBUser.CreatedAt
				// the user has not been updated since the last read, so the user can be updated now
				user.UpdatedAt = jsontime.Now()
				if err := repo.setValueAsJSON(ctx, user.ID, user, defDuration200Years); err != nil {
					return err
				}
				updated = true
			} else {
				lastUpdatedAt = currDBUser.UpdatedAt
			}
		}
		// we must only leave the loop is the user was updated:
		if updated {
			return nil
		} else { // just to be aware a repo logic is broken after some update
			return errors.New("repository internal logic error")
		}

	} else { // just a standard pessimistic lock

		lockRes := repo.lockRecord(ctx, user.ID)
		defer repo.unlockRecord()

		switch lockRes {
		case nil, redislock.Ok: // the user record has been successfully locked
			// patch CreatedAt
			currDBUser, err := repo.GetUser(ctx, user.ID)
			if err != nil {
				return err
			}
			user.CreatedAt = currDBUser.CreatedAt
			user.UpdatedAt = jsontime.Now()
			if err := repo.setValueAsJSON(ctx, user.ID, user, defDuration200Years); err != nil {
				return err
			}
		default:
			return lockRes
		}
	}

	return nil
}
