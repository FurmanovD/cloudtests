package userservice

import "github.com/FurmanovD/cloudtests/internal/pkg/jsontime"

// User contains user details
type User struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	MiddleName string `json:"middle-name"`
	LastName   string `json:"last-name"`
	Address    string `json:"address"`

	CreatedAt jsontime.JSONTime `json:"created"`
	UpdatedAt jsontime.JSONTime `json:"updated"`
	// TODO(DF) DeletedAt
}
