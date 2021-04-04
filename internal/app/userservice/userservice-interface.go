package userservice

import (
	"context"

	model "github.com/FurmanovD/cloudtests/internal/pkg/model/userservice"
)

// Service describes the service functions.
type Service interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	GetUser(ctx context.Context, userID string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User, optimisticLock bool) error
	DeleteUser(ctx context.Context, userID string) error
}
