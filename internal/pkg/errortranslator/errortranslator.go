package errortranslator

import (
	svcerrors "github.com/FurmanovD/cloudtests/internal/pkg/errors"
	"github.com/go-redis/redis/v8"
)

// RepositoryErrorTranslatorFn translates a repository error to a service one.
type RepositoryErrorTranslatorFn func(error) error

// NewRepositoryErrorTranslator creates a translation function
func NewRepositoryErrorTranslator() RepositoryErrorTranslatorFn {
	return func(err error) error {
		switch err {
		case redis.Nil:
			return svcerrors.ErrNotFound
			// TODO add here specific errors if we need to redirect them to a user
		default:
			return err
		}
	}
}
