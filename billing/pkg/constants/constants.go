package constants

const (
	ServiceIdentifier   = "billing"
	ApplicationTraceKey = "x-app-trace-id"

	GenericErrorMessage         = "Uh-Oh! Something went wrong."
	RequestBodyValidationFailed = "validation for input payload"
	ResourceNotFoundMessage     = "requested resource is not found"
	MethodNotAllowedMessage     = "used request method is not supported for resource"
	ValidationFailed            = "validation failed"
)

//goland:noinspection SpellCheckingInspection,SpellCheckingInspection
const (
	CustomerPrefix    = "CUST"
	LoanPrefix        = "LN"
	LoanAccountPrefix = "LNAC"
	PaymentPrefix     = "PAY"
)
