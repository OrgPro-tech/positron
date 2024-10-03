package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field  string `json:"field"`
	Error  string `json:"error"`
	Nested string `json:"nested,omitempty"`
}

// ValidateJSONBody is a generic function to validate JSON body in Fiber APIs
func ValidateJSONBody[T any](c *fiber.Ctx) (T, []ValidationError, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return body, nil, fmt.Errorf("error parsing JSON body: %v", err)
	}

	errors := validate(reflect.ValueOf(body), "")
	return body, errors, nil
}

// validate performs validation on the given struct based on struct tags
func validate(v reflect.Value, prefix string) []ValidationError {
	var errors []ValidationError

	// If pointer, get the underlying element
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name
		tag := fieldType.Tag.Get("validate")

		// If the field is a nested struct, recursively validate it
		if field.Kind() == reflect.Struct {
			nestedPrefix := prefix
			if nestedPrefix != "" {
				nestedPrefix += "."
			}
			nestedPrefix += fieldName

			// Check if the nested struct itself is required
			if tag == "required" && isZeroValue(field) {
				errors = append(errors, ValidationError{
					Field:  nestedPrefix,
					Error:  "nested struct is required",
					Nested: "",
				})
			} else {
				nestedErrors := validate(field, nestedPrefix)
				errors = append(errors, nestedErrors...)
			}
			continue
		}

		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")

		for _, rule := range rules {
			parts := strings.Split(rule, "=")
			var err error

			switch parts[0] {
			case "required":
				err = validateRequired(field)
			case "min":
				err = validateMin(field, parts[1])
			case "max":
				err = validateMax(field, parts[1])
			case "regex":
				err = validateRegex(field, parts[1])
			}

			if err != nil {
				fullFieldName := prefix
				if fullFieldName != "" {
					fullFieldName += "."
				}
				fullFieldName += fieldName

				errors = append(errors, ValidationError{
					Field:  fullFieldName,
					Error:  err.Error(),
					Nested: prefix,
				})
			}
		}
	}

	return errors
}

func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func validateRequired(field reflect.Value) error {
	if isZeroValue(field) {
		return fmt.Errorf("field is required")
	}
	return nil
}

func validateMin(field reflect.Value, param string) error {
	num, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return fmt.Errorf("invalid min value: %s", param)
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() < int64(num) {
			return fmt.Errorf("value must be at least %v", num)
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() < num {
			return fmt.Errorf("value must be at least %v", num)
		}
	case reflect.String:
		if float64(len(field.String())) < num {
			return fmt.Errorf("length must be at least %v", num)
		}
	default:
		return fmt.Errorf("min validation not supported for this type")
	}

	return nil
}

func validateMax(field reflect.Value, param string) error {
	num, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return fmt.Errorf("invalid max value: %s", param)
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() > int64(num) {
			return fmt.Errorf("value must be at most %v", num)
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() > num {
			return fmt.Errorf("value must be at most %v", num)
		}
	case reflect.String:
		if float64(len(field.String())) > num {
			return fmt.Errorf("length must be at most %v", num)
		}
	default:
		return fmt.Errorf("max validation not supported for this type")
	}

	return nil
}

func validateRegex(field reflect.Value, param string) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("regex validation only supports string fields")
	}

	re, err := regexp.Compile(param)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %s", param)
	}

	if !re.MatchString(field.String()) {
		return fmt.Errorf("value does not match pattern %s", param)
	}

	return nil
}
