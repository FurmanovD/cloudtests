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
	model "github.com/FurmanovD/cloudtests/internal/pkg/model/userservice"
)

// GetUserRequest contains a UserID to be retrieved.
type GetUserRequest struct {
	UserID string
}

// GetUserResponse describes a User.
type GetUserResponse struct {
	model.User
}

// NewGetUserEndpoint creates an endpoint for GetUserRequest processing
func NewGetUserEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req, ok := request.(GetUserRequest)
		if !ok {
			return nil, svcerror.ErrInvalidRequest
		}
		user, err := s.GetUser(ctx, req.UserID)
		if err != nil {
			return nil, err
		}

		return GetUserResponse{
			User: user,
		}, nil
	}
}

func NewGetUserRequestDecoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {

		userID := mux.Vars(r)[api.ParamUserID]

		if userID == "" {
			return nil, svcerror.ErrInvalidRequest
		}
		return GetUserRequest{
			UserID: userID,
		}, nil
	}
}
