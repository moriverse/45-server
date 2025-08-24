package response

import "github.com/gin-gonic/gin"

// APIError defines the structure for a standard API error response.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Data sends a standard success response with a JSON payload.
// It uses the provided HTTP status code and marshals the data interface to JSON.
func Data(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// Error sends a standard error response.
// It uses the provided HTTP status code and marshals the APIError struct to JSON.
func Error(c *gin.Context, status int, err APIError) {
	c.JSON(status, gin.H{"error": err})
}
