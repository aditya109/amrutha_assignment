package models

const (
	BAD_REQUEST_NAME                       = "BAD_REQUEST"
	INTERNAL_SERVER_ERROR_NAME             = "INTERNAL_SERVER_ERROR"
	URL_PARAMETER_VALIDATION_FAILED_NAME   = "URL_PARAMETER_VALIDATION_FAILED"
	QUERY_PARAMETER_VALIDATION_FAILED_NAME = "QUERY_PARAMETER_VALIDATION_FAILED"
)

const (
	BAD_REQUEST_MESSAGE_MESSAGE               = "Bad Request"
	INTERNAL_SERVER_ERROR_MESSAGE             = "Something Went Wrong"
	URL_PARAMETER_VALIDATION_FAILED_MESSAGE   = "URL Parameter Validation Failed"
	QUERY_PARAMETER_VALIDATION_FAILED_MESSAGE = " Query Parameter Validation Failed"
)

type Response struct {
	Message string
	Data    interface{}
	Status  int
	Error   interface{}
}

// SuccessResponseDTO is structure for success response for all APIs
type SuccessResponseDTO struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Data    *interface{} `json:"data"`
}

// Error is struct for error data in failure response
type Error struct {
	Code              string `json:"code"`
	Message           string `json:"message"`
	ResolutionMessage string `json:"resolutionMessage"`
	Data              any    `json:"data"`
}

// FailureResponse is structure for failure response for all APIs
type FailureResponse struct {
	Success bool  `json:"success"`
	Error   Error `json:"error"`
}

func SuccessResponse(data interface{}, message string) SuccessResponseDTO {
	return SuccessResponseDTO{
		Success: true,
		Data:    &data,
		Message: message,
	}
}

func ErrorResponse(error Error) FailureResponse {
	return FailureResponse{
		Success: false,
		Error:   error,
	}
}

func InternalServerError(e *Error) FailureResponse {
	MESSAGE := INTERNAL_SERVER_ERROR_MESSAGE

	{
		if e.Message != "" {
			MESSAGE = e.Message
		}
	}

	return ErrorResponse(Error{
		Code:              "SERVICE_ERROR_101",
		Message:           MESSAGE,
		ResolutionMessage: "Please try again",
		Data:              &e.Data,
	})
}
