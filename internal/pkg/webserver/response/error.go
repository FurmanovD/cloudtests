package response

// Error used to return an error description.
type Error struct {
	Error string `json:"error,omitempty"`
}
