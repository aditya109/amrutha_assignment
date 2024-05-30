package helpers

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aditya109/amrutha_assignment/pkg/logger"
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
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

func GetDateFromString(dateString string, incomingDateTimeFormat string, outputDateTimeFormat string) string {
	layout := "2006-01-02"
	if incomingDateTimeFormat == "YYYY-MM-DD" {
		myDate, _ := time.Parse(layout, dateString)
		return myDate.Format(outputDateTimeFormat)
	}
	return ""
}

func BytesToMapStringInterface(bytes []byte) (map[string]interface{}, error) {
	var target map[string]interface{}

	err := json.Unmarshal(bytes, &target)

	if err != nil {
		return nil, err
	}

	return target, err
}

func PayloadToMapStringInterface(payload interface{}) (map[string]interface{}, error) {

	var finalPayload map[string]interface{}

	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return map[string]interface{}{}, err
	}

	err = json.Unmarshal(payloadBytes, &finalPayload)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return finalPayload, nil
}

func BytesToStruct[T any](data []byte) (T, error) {
	var realData T

	// unmarshal bytes to real data
	err := json.Unmarshal(data, &realData)
	if err != nil {
		return realData, fmt.Errorf("error while processing bytes to struct: %v", err)
	}

	return realData, nil
}
func InterfaceToStruct[T any](data interface{}) (T, error) {
	var realData T

	bytes, err := ConvertStructIntoBytesArray(data)
	if err != nil {
		return realData, err
	}

	return BytesToStruct[T](bytes)
}

func GetPrettyPrintJsonStringFromStruct(ob any) string {
	b, err := json.MarshalIndent(ob, "", "  ")
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func GetJsonStringFromStruct(ob any) (string, error) {
	b, err := json.Marshal(ob)
	if err != nil {
		return "", fmt.Errorf("error encoding JSON: %v", err)
	}
	return string(b), nil
}

func StructToBytes(ob any) []byte {
	return []byte(fmt.Sprintf("%v", ob))
}

func StringConvertToInt(val string) (int, error) {
	if i, err := strconv.Atoi(strings.Split(val, ".")[0]); err == nil {
		return int(i), nil
	} else {
		return -1, fmt.Errorf("error while converting to int: %s", val)
	}
}

func StringConvertToFloat64(val string) float64 {
	log := logger.GetInternalContextLogger(StringConvertToFloat64)
	if i, err := strconv.ParseFloat(val, 64); err == nil {
		return i
	} else {
		log.Fatalf("error while converting to int: %s", val)
	}
	return -1
}

func BytesToJson(data []byte) (map[string]interface{}, error) {
	var finalData map[string]interface{}
	err := json.Unmarshal(data, &finalData)
	return finalData, err
}

func PrintBytesBufferToString(buf *bytes.Buffer) {
	log := logger.GetInternalContextLogger(PrintBytesBufferToString)
	log.Println(buf.String())
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

func ConvertToMap(obj interface{}, enableSnakeCasingForKey bool) (map[string][]string, error) {
	// Initialize an empty map to hold the converted values
	convertedMap := make(map[string][]string)

	// Get the type of the object
	objType := reflect.TypeOf(obj)

	// Ensure the object is a struct
	if objType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input object is not a struct")
	}

	// Get the value of the object
	objValue := reflect.ValueOf(obj)

	// Iterate through the fields of the struct
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)

		// Convert field value to string
		var stringValue string
		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			stringValue = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			stringValue = strconv.FormatUint(fieldValue.Uint(), 10)
		case reflect.Float32:
			stringValue = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 32)
		case reflect.Float64:
			stringValue = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.Bool:
			stringValue = strconv.FormatBool(fieldValue.Bool())
		case reflect.String:
			stringValue = fieldValue.String()
		// Handle other types
		default:
			stringValue = fmt.Sprintf("%v", fieldValue.Interface())
		}

		// Store the field name and converted value in the map
		if enableSnakeCasingForKey {
			convertedMap[ToLowerSnakeCase(field.Name)] = []string{stringValue}
		} else {
			convertedMap[field.Name] = []string{stringValue}
		}
	}

	return convertedMap, nil
}

func ToLowerSnakeCase(input string) string {
	var buf bytes.Buffer
	buf.Grow(len(input) * 2)

	// Iterate over each character in the input string
	for i, r := range input {
		if unicode.IsUpper(r) {
			// If the character is uppercase and not the first character,
			// insert an underscore before it
			if i != 0 {
				buf.WriteByte('_')
			}
			// Convert the uppercase character to lowercase
			buf.WriteRune(unicode.ToLower(r))
		} else {
			// If the character is not uppercase, just append it
			buf.WriteRune(r)
		}
	}

	return buf.String()
}
