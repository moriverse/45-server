package errors

import "fmt"

// NotFoundError is a generic error for any resource that cannot be found.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// DataInconsistencyError is returned when an unexpected state is found in the database,
// for example, an auth record exists but its corresponding user record does not.
type DataInconsistencyError struct {
	Details string
}

func (e DataInconsistencyError) Error() string {
	return fmt.Sprintf("data inconsistency detected: %s", e.Details)
}