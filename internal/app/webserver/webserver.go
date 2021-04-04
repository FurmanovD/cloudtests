package webserver

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/FurmanovD/cloudtests/internal/app/userservice"
	svcerrors "github.com/FurmanovD/cloudtests/internal/pkg/errors"
	"github.com/FurmanovD/cloudtests/internal/pkg/webserver/response"
)

// NewService creates HTTP endpoints.
func NewHTTPServer(usrSvc userservice.Service) http.Handler {

	// set-up router and initialize http endpoints
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(wrapErrorResponse),
	}

	for _, wse := range webserverEndpoints {
		r.Methods(wse.Method).Path(wse.Path).
			Handler(kithttp.NewServer(
				wse.CreateEndpointFn(usrSvc),
				wse.CreateDecodeFn(),
				encodeResponse,
				options...,
			))
	}

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set correct content type
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	type errorer interface {
		error() error
	}

	if e, ok := response.(errorer); ok && e.error() != nil {
		wrapErrorResponse(ctx, e.error(), w)
		return nil
	}
	return kithttp.EncodeJSONResponse(ctx, w, response)
}

func wrapErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("wrapError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCodeFromError(err))

	json.NewEncoder(w).Encode(response.Error{
		Error: err.Error(),
	})
}

func httpCodeFromError(err error) int {
	//TODO(DF): add error converting to HTTP status here
	switch err {
	case svcerrors.ErrInvalidRequest:
		return http.StatusBadRequest
	case svcerrors.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
