package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/FurmanovD/cloudtests/internal/app/userservice"
	"github.com/FurmanovD/cloudtests/internal/app/webserver/constants/api"
	svcerror "github.com/FurmanovD/cloudtests/internal/pkg/errors"
	model "github.com/FurmanovD/cloudtests/internal/pkg/model/userservice"
	"github.com/FurmanovD/cloudtests/internal/pkg/webserver/response"
)

// UpdateUserRequest describes a User to be updated.
type UpdateUserRequest struct {
	model.User
	OptimisticLock bool
}

// UpdateUserResponse describes a created User ID.
type UpdateUserResponse response.Empty

// NewUpdateUserEndpoint creates an endpoint for UpdateUserRequest processing
func NewUpdateUserEndpoint(s userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req, ok := request.(UpdateUserRequest)
		if !ok {
			return nil, svcerror.ErrInvalidRequest
		}

		err := s.UpdateUser(ctx, req.User, req.OptimisticLock)
		if err != nil {
			return nil, err
		}

		return UpdateUserResponse{}, nil
	}
}

func NewUpdateUserRequestDecoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		request := UpdateUserRequest{}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, fmt.Errorf("failed to parse request JSON: %w", err)
		}

		// add UserID
		userID := mux.Vars(r)[api.ParamUserID]
		if userID == "" {
			return nil, svcerror.ErrInvalidRequest
		}
		request.User.ID = userID

		// add lock type
		lockTypeStr := r.URL.Query().Get(api.ParamOptimisticLock)
		if lockTypeStr == "1" || lockTypeStr == "true" {
			request.OptimisticLock = true
		}

		return request, nil
	}
}
