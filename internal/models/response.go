package models

// Response represents a standard API response
// @Description Standard response format for all API endpoints
type Response struct {
	// @Description Indicates if the request was successful
	Success bool `json:"success"`
	// @Description Human-readable message about the response
	Message string `json:"message"`
	// @Description HTTP status code
	StatusCode int `json:"statusCode"`
	// @Description Response payload data
	Data interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, statusCode int, err error) *ErrorResponse {
	errorResponse := &ErrorResponse{
		Success:    false,
		Message:    message,
		StatusCode: statusCode,
	}
	if err != nil {
		errorResponse.Error = err.Error()
	}
	return errorResponse
}
