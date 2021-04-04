package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"github.com/FurmanovD/cloudtests/internal/pkg/redislock"
)

const (
	// 200 years duration to save records.
	defDuration200Years = time.Hour * time.Duration(200*365*24)
	defDB               = 0
	tableKeyPrefixUser  = "user-"
	tableKeyPrefixLock  = "lock-"

	// timeouts and periods:
	lockExpiration = time.Minute * 2
	lockTimeout    = time.Second * 5
	lockRetry      = time.Millisecond * 100
)

type repository struct {
	client    *redis.Client
	keyLocker redislock.RedisLock
}

// New creates a repository instance.
func New(addr, pass string) (Repository, error) {
	log.Trace("repository.New(%v, %v) called", addr, pass)

	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       defDB,
	})

	if err := conn.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &repository{
			client:    conn,
			keyLocker: redislock.NewRedisLocker(conn),
		},
		nil
}

// ================== Utility methods ======================

// setValueAsJSON sets any value marshalled as JSON to a key.
func (repo *repository) setValueAsJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	log.Trace("repository.setValueAsJSON(%v,%v,%v) called", key, value, ttl)

	strValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = repo.client.Set(ctx, key, strValue, ttl).Err()
	if err != nil {
		log.Errorf("redis Set(%v, %v, %v) error: %v", key, strValue, ttl, err)
		return err
	}

	return nil
}

// lockRecord locks a record with default retry parameters.
func (repo *repository) lockRecord(ctx context.Context, key string) redislock.RedisLockError {
	log.Trace("repository.lockRecord(%v) called", key)

	return repo.keyLocker.ObtainLock(
		ctx,
		key,
		lockExpiration,
		lockTimeout,
		lockRetry,
	)
}

// unlockRecord just to be consistent with lockRecord(...) method.
func (repo *repository) unlockRecord() {
	log.Trace("repository.ulockRecord() called")

	repo.keyLocker.Unlock()
}
