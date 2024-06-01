package helpers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aditya109/amrutha_assignment/billing/pkg/constants"
	"github.com/aditya109/atomic"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
	"log"
	"strings"
	"time"
)

func CreatePointerForValue[T any](v T) *T {
	return &v
}

func FormatCurrency(value decimal.Decimal) string {
	tag := language.Make("en-IN")
	printer := message.NewPrinter(tag)
	frac := 2
	format := "\u20b9" + "%m"
	opts := []number.Option{
		number.IncrementString(fmt.Sprintf("0.%0[1]*d", frac+1, 5)),
		number.MaxFractionDigits(frac),
		number.MinFractionDigits(frac),
	}
	num := printer.Sprint(number.Decimal(value.InexactFloat64(), opts...))
	return printer.Sprintf(format, num)
}

func GetDateAsTimeFromString(input string, formats ...string) time.Time {
	var incomingDateFormat = time.RFC3339Nano
	switch len(formats) {
	case 1:
		if formats[0] != "" {
			incomingDateFormat = formats[0]
		}
	case 2:
		if formats[0] != "" {
			incomingDateFormat = formats[0]
		}
	}
	t, _ := time.Parse(incomingDateFormat, input)

	return t
}

func ConvertStructIntoHashString(ob any) (string, error) {
	// Encode the struct into JSON
	jsonData, err := ConvertStructIntoBytesArray(ob)
	if err != nil {
		return "", err
	}

	// Hash the JSON data using SHA-256
	hash := sha256.Sum256(jsonData)

	// Print the hash as a hexadecimal string
	return fmt.Sprintf("%x", hash), nil
}

func ConvertStructIntoBytesArray(a any) ([]byte, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %v", err)
	}
	return bytes, nil
}

func PrintLandingRequestCurl(c *gin.Context) {
	curl, err := atomic.Boom(c.Request)
	if err != nil {
		log.Println(err)
	}
	log.Printf("trace_id: %s, incoming request cURL: %s", getTraceId(c), curl)
}

func getTraceId(c *gin.Context) string {
	if c != nil {
		return c.Request.Header.Get(constants.ApplicationTraceKey)
	}
	return ""
}

func CreateUniqueDisplayId(ob interface{}, prefix string) (string, error) {
	if uniqueHash, err := ConvertStructIntoHashString(ob); err != nil {
		return "", fmt.Errorf("error while create unique hashable id for customer: %w", err)
	} else {
		if len(uniqueHash) > 16 {
			uniqueHash = uniqueHash[:16]
		}
		return fmt.Sprintf("%s%s", prefix, strings.ToUpper(uniqueHash)), nil
	}
}
