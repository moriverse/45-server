package context

// contextKey is an unexported type for context keys to prevent collisions.
type contextKey string

// Defines the keys used to store and retrieve values from the request context.
const (
	// UserIDKey is the key for the user ID in the context.
	UserIDKey = contextKey("userID")
	// LoggerKey is the key for the logger in the context.
	LoggerKey = contextKey("logger")
)
