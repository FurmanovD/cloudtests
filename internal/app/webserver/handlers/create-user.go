package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/FurmanovD/cloudtests/internal/app/userservice"
	svcerror "github.com/FurmanovD/cloudtests/internal/pkg/errors"
	model "github.com/FurmanovD/cloudtests/internal/pkg/model/userservice"
)

// CreateUserRequest describes a User to be created.
type CreateUserRequest struct {
	model.User
}

// CreateUserResponse describes a created User ID.
type CreateUserResponse struct {
	UserID string `json:"userId,omitempty"`
}

// NewCreateUserEndpoint creates an endpoint for CreateUserRequest processing
func NewCreateUserEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req, ok := request.(CreateUserRequest)
		if !ok {
			return nil, svcerror.ErrInvalidRequest
		}
		userID, err := s.CreateUser(ctx, req.User)
		if err != nil {
			return nil, err
		}

		return CreateUserResponse{
			UserID: userID,
		}, nil
	}
}

func NewCreateUserRequestDecoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		request := CreateUserRequest{}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, fmt.Errorf("failed to parse request JSON: %w", err)
		}
		return request, nil
	}
}
