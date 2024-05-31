package constants

const (
	SERVICE_IDENTIFIER    = "billing"
	APPLICATION_TRACE_KEY = "x-app-trace-id"

	SOMETHING_WENT_WRONG           = "something went wrong"
	GENERIC_ERROR_MESSAGE          = "Uh-Oh! Something went wrong."
	REQUEST_BODY_VALIDATION_FAILED = "validation for input payload"
	RESOURCE_NOT_FOUND_MESSAGE     = "requested resource is not found"
	METHOD_NOT_ALLOWED_MESSAGE     = "used request method is not supported for resource"
	VALIDATION_FAILED              = "validation failed"
	INTERNAL_SERVER_ERROR          = "internal server error"
)

const (
	CUSTOMER_PREFIX     = "CUST"
	LOAN_PREFIX         = "LN"
	LOAN_ACCOUNT_PREFIX = "LNAC"
)
