package models

//goland:noinspection ALL
const (
	BadRequestName                     = "BAD_REQUEST"
	InternalServerErrorName            = "INTERNAL_SERVER_ERROR"
	UrlParameterValidationFailedName   = "URL_PARAMETER_VALIDATION_FAILED"
	QueryParameterValidationFailedName = "QUERY_PARAMETER_VALIDATION_FAILED"
)

//goland:noinspection ALL
const (
	BadRequestMessageMessage              = "Bad Request"
	InternalServerErrorMessage            = "Something Went Wrong"
	UrlParameterValidationFailedMessage   = "URL Parameter Validation Failed"
	QueryParameterValidationFailedMessage = " Query Parameter Validation Failed"
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
	MESSAGE := InternalServerErrorMessage

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
