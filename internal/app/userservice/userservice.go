package userservice

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/FurmanovD/cloudtests/internal/pkg/errortranslator"
	model "github.com/FurmanovD/cloudtests/internal/pkg/model/userservice"
	"github.com/FurmanovD/cloudtests/internal/pkg/repository"
)

const (
	UserIDRegexp = "[-a-zA-Z0-9]*"
)

// service implements the service.
type service struct {
	repository    repository.Repository
	errTranslator errortranslator.RepositoryErrorTranslatorFn
}

// New creates and returns a new service instance.
func New(
	repo repository.Repository,
	translator errortranslator.RepositoryErrorTranslatorFn,
) Service {
	return &service{
		repository:    repo,
		errTranslator: translator,
	}
}

// CreateUser creates a new user.
func (s *service) CreateUser(ctx context.Context, user model.User) (string, error) {
	log.Trace("userservice.CreateUser(%v) called", user)

	userID, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		log.Errorf("CreateUser error: %v", err)
	}
	return userID, s.errTranslator(err)
}

// GetUser returns a user by it's ID.
func (s *service) GetUser(ctx context.Context, userID string) (model.User, error) {
	log.Trace("userservice.GetUser(%v) called", userID)

	user, err := s.repository.GetUser(ctx, userID)
	if err != nil {
		log.Errorf("GetUser error: %v", err)
	}
	return user, s.errTranslator(err)
}

// UpdateUser updates a user by its ID.
func (s *service) UpdateUser(ctx context.Context, user model.User, optimisticLock bool) error {
	log.Trace("userservice.UpdateUser(%v) called", user)

	err := s.repository.UpdateUser(ctx, user, optimisticLock)
	if err != nil {
		log.Errorf("UpdateUser error: %v", err)
	}
	return s.errTranslator(err)
}

// DeleteUser deletes a user by its ID.
func (s *service) DeleteUser(ctx context.Context, userID string) error {
	log.Trace("userservice.DeleteUser(%v) called", userID)

	err := s.repository.DeleteUser(ctx, userID)
	if err != nil {
		log.Errorf("DeleteUser error: %v", err)
	}
	return s.errTranslator(err)
}
