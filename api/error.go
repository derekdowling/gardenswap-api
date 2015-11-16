package api

// Error is a JSON Spec error format
type Error struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

// ISE is a convenience function that returns a pre-packaged 500 APIError
// for times when we don't want to expose internal errors
func ISE(err error) *Error {
	return &Error{
		Status: 500,
		Detail: "Whoops! Something went wrong.",
	}
}
