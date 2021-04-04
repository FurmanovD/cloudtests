package webserver

import (
	"fmt"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/FurmanovD/cloudtests/internal/app/userservice"
	"github.com/FurmanovD/cloudtests/internal/app/webserver/constants/api"
	"github.com/FurmanovD/cloudtests/internal/app/webserver/handlers"
)

//
type endpointDescription struct {
	Method           string
	Path             string
	CreateEndpointFn func(userservice.Service) endpoint.Endpoint
	CreateDecodeFn   func() kithttp.DecodeRequestFunc
}

var (
	webserverEndpoints []endpointDescription = []endpointDescription{
		{
			// HTTP Post - CreateUser /api/v1/users
			Method:           "POST",
			Path:             api.PathPrefix + "/users",
			CreateEndpointFn: handlers.NewCreateUserEndpoint,
			CreateDecodeFn:   handlers.NewCreateUserRequestDecoder,
		},
		{
			// HTTP Get - GetUser /api/v1/users/{userID}
			Method: "GET",
			Path: fmt.Sprintf(
				"%s/users/{%s:%s}",
				api.PathPrefix,
				api.ParamUserID,
				userservice.UserIDRegexp,
			),
			CreateEndpointFn: handlers.NewGetUserEndpoint,
			CreateDecodeFn:   handlers.NewGetUserRequestDecoder,
		},
		{
			// HTTP PUT - UpdateUser /api/v1/users/{userID}?optimisticlock=0|1|true|false
			//TODO(DF) "%s/users/{%s:%s}?{%s}",
			Method: "PUT",
			Path: fmt.Sprintf(
				"%s/users/{%s:%s}",
				api.PathPrefix,
				api.ParamUserID,
				userservice.UserIDRegexp,
			),
			CreateEndpointFn: handlers.NewUpdateUserEndpoint,
			CreateDecodeFn:   handlers.NewUpdateUserRequestDecoder,
		},
		{
			// HTTP Delete - DeleteUser /api/v1/users/{userID}
			Method: "DELETE",
			Path: fmt.Sprintf(
				"%s/users/{%s:%s}",
				api.PathPrefix,
				api.ParamUserID,
				userservice.UserIDRegexp,
			),
			CreateEndpointFn: handlers.NewDeleteUserEndpoint,
			CreateDecodeFn:   handlers.NewDeleteUserRequestDecoder,
		},
	}
)
