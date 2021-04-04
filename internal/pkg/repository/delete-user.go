package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/FurmanovD/cloudtests/internal/pkg/redislock"
)

// DeleteUser deletes a user record by its ID.
func (repo *repository) DeleteUser(ctx context.Context, userID string) error {
	log.Trace("repository.DeleteUser(%v) called", userID)

	lockRes := repo.lockRecord(ctx, userID)
	if lockRes != redislock.Ok {
		return lockRes
	}
	defer repo.unlockRecord()

	deleted, err := repo.client.Del(ctx, userID).Result()
	if err != nil {
		log.Errorf("redis Del(%v) error: %v", userID, err)
	} else {
		if deleted == 0 {
			err = redis.Nil // Not Found
		}
	}

	return err
}
