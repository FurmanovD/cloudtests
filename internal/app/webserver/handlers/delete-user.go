package handlers

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/FurmanovD/cloudtests/internal/app/userservice"
	"github.com/FurmanovD/cloudtests/internal/app/webserver/constants/api"
	svcerror "github.com/FurmanovD/cloudtests/internal/pkg/errors"
	"github.com/FurmanovD/cloudtests/internal/pkg/webserver/response"
)

// DeleteUserRequest contains a UserID to be deleted.
type DeleteUserRequest struct {
	UserID string
}

// DeleteUserResponse is empty.
type DeleteUserResponse response.Empty

// NewDeleteUserEndpoint creates an endpoint for DeleteUserRequest processing
func NewDeleteUserEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req, ok := request.(DeleteUserRequest)
		if !ok {
			return nil, svcerror.ErrInvalidRequest
		}
		err := s.DeleteUser(ctx, req.UserID)
		if err != nil {
			return nil, err
		}

		return DeleteUserResponse{}, nil
	}
}

func NewDeleteUserRequestDecoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {

		userID := mux.Vars(r)[api.ParamUserID]

		if userID == "" {
			return nil, svcerror.ErrInvalidRequest
		}
		return DeleteUserRequest{
			UserID: userID,
		}, nil
	}
}
