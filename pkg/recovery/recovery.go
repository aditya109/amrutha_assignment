package recovery

import (
	"github.com/aditya109/amrutha_assignment/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiPanicRecovery(c *gin.Context, recovered interface{}) {
	if recovered != nil {
		c.JSON(http.StatusInternalServerError, models.InternalServerError(&models.Error{Data: recovered}))
		c.Abort()
		return
	}
}

func Responder(c *gin.Context) {
	__response, _ := c.Get("Response")
	__errorResponse, _ := c.Get("Error")

	if __errorResponse != nil {
		errorResponse := __errorResponse.(models.Response)
		c.JSON(errorResponse.Status, errorResponse.Error)
	} else if __response != nil {
		response := __response.(models.Response)
		c.JSON(response.Status, models.SuccessResponse(response.Data, response.Message))
	}

	c.Abort()
	return
}
