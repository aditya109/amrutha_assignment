package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
	"github.com/aditya109/amrutha_assignment/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strings"
)

type Backdrop interface {
	// Response Error(statusCode int, err fallback.Error)
	Response(statusCode int, data interface{})
	GetLogger(function interface{}) *logrus.Entry
	GetMetaStore() *MetaStore
	GetAllHeaders() map[string]string
	SetCustomErrorMessage(msg string)
	SetCustomResolutionMessage(msg string)
	SetCustomTraceId(traceId string)
	GetDatabaseInstance() *gorm.DB
	SetDatabaseLogger()
	GetContext() *gin.Context
	SetDatabaseInstance(*gorm.DB)
	Error(int, models.Error)
	ReadRequestPayload(dst any) error
	SetStatusCodeForResponse(code int)
}

func GetNewBackdrop(c *gin.Context, db *gorm.DB) Backdrop {
	return &callContext{
		metaStore: &MetaStore{
			TraceId: bindTraceIdToContextAsValue(),
		},
		dbInstance: db,
		ginContext: c,
	}
}

// SetStatusCodeForResponse implements models.Backdrop.
func (c *callContext) SetStatusCodeForResponse(code int) {
	c.metaStore.StatusCode = code
}

func (c *callContext) SetCustomErrorMessage(msg string) {
	c.metaStore.CustomErrorMessage = msg
}

func (c *callContext) SetCustomResolutionMessage(msg string) {
	c.metaStore.CustomResolutionMessage = msg
}

func (c *callContext) GetAllHeaders() map[string]string {
	var fHeaders = make(map[string]string)
	fHeaders[constants.ApplicationTraceKey] = c.metaStore.TraceId
	return fHeaders
}

func (c *callContext) GetMetaStore() *MetaStore {
	return c.metaStore
}

func (c *callContext) GetLogger(function interface{}) *logrus.Entry {
	return logger.GetInternalContextLogger(function, c.metaStore.TraceId)
}

func (c *callContext) Response(statusCode int, data interface{}) {
	entry := c.GetLogger(c.Response)
	entry.WithFields(logrus.Fields{"response": data}).Info()
	c.ginContext.JSON(statusCode, models.SuccessResponse(data, ""))
}

// GetDatabaseInstance implements Backdrop.
func (c *callContext) GetDatabaseInstance() *gorm.DB {
	return c.dbInstance
}
func (c *callContext) GetContext() *gin.Context {
	return c.ginContext
}
func (c *callContext) GetMode() string {
	return c.mode
}
func (c *callContext) SetCustomTraceId(traceId string) {
	c.metaStore.TraceId = traceId
	// c.SetDatabaseLogger()
}

// SetDatabaseInstance implements Backdrop.
func (c *callContext) SetDatabaseInstance(db *gorm.DB) {
	c.dbInstance = db
}
func (c *callContext) SetDatabaseLogger() {
	c.dbInstance.Logger = logger.GetCustomGormLogger(c.metaStore.TraceId)
}

// SetMode implements Backdrop.
func (c *callContext) SetMode(mode string) {
	c.mode = mode
}

func (c *callContext) Error(statusCode int, err models.Error) {
	entry := c.GetLogger(c.Error)
	if c.metaStore.StatusCode != 0 {
		statusCode = c.metaStore.StatusCode
	} else if statusCode == 0 {
		entry.Info("invalid status code")
		c.ginContext.JSON(http.StatusInternalServerError, models.ErrorResponse(err))
		return
	}

	if c.metaStore.CustomErrorMessage != "" {
		err.Message = c.metaStore.CustomErrorMessage
	} else {
		switch {
		case statusCode == http.StatusMethodNotAllowed:
			err.ResolutionMessage = constants.MethodNotAllowedMessage
			err.Message = constants.GenericErrorMessage
		case statusCode == http.StatusUnprocessableEntity:
			err.ResolutionMessage = constants.RequestBodyValidationFailed
			err.Message = constants.GenericErrorMessage
		case statusCode == http.StatusNotFound:
			err.ResolutionMessage = constants.ResourceNotFoundMessage
			err.Message = constants.GenericErrorMessage
		case statusCode == http.StatusInternalServerError:
			fallthrough
		default:
			err.Message = constants.GenericErrorMessage
		}
	}

	if c.metaStore.CustomResolutionMessage != "" {
		err.ResolutionMessage = c.metaStore.CustomResolutionMessage
	}
	err.ResolutionMessage = fmt.Sprintf("%s, trace_id: %s", err.ResolutionMessage, c.metaStore.TraceId)
	log.Printf("status_code: %d, err: %v", statusCode, err)
	c.ginContext.AbortWithStatusJSON(statusCode, models.ErrorResponse(err))
}

func bindTraceIdToContextAsValue() string {
	return uuid.New().String()
}

func GetTraceId(c *gin.Context) string {
	if c != nil {
		return c.Request.Header.Get(constants.ApplicationTraceKey)
	}
	return ""
}

func (c *callContext) ReadRequestPayload(dst any) error {
	body := c.ginContext.Request.Body
	maxBytes := 1_048_576
	body = http.MaxBytesReader(c.ginContext.Writer, body, int64(maxBytes))
	decoder := json.NewDecoder(body)
	//decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			return err
		default:
			return err
		}
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
func (c *callContext) IsJsonStructValidationSuccessful(data interface{}) error {
	// Create a new validator instance
	validatorInstance := validator.New()

	// Validate the struct
	var s = make([]string, 0)
	err := validatorInstance.Struct(data)
	if err != nil {
		// Extract detailed validation errors
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return err
		}
		for _, validationError := range err.(validator.ValidationErrors) {
			s = append(s, fmt.Sprintf("field: %s, error: %s", validationError.Field(), fmt.Sprintf("%s", validationError.Tag())))
		}

	}

	return fmt.Errorf(strings.Join(s, ";"))
}
