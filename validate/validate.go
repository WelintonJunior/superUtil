package validate

import (
	"fmt"
	"reflect"
	"strings"
)

func ValidateStruct[T any](item *T) error {
	if item == nil {
		return fmt.Errorf("request body cannot be nil")
	}

	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() != reflect.Pointer || itemValue.IsNil() {
		return fmt.Errorf("request body must be a non-nil pointer")
	}

	itemValue = itemValue.Elem()
	if itemValue.Kind() != reflect.Struct {
		return fmt.Errorf("request body must be a struct")
	}

	itemType := itemValue.Type()

	for i := 0; i < itemValue.NumField(); i++ {
		field := itemValue.Field(i)
		fieldType := itemType.Field(i)

		validateTag := fieldType.Tag.Get("validate")
		isRequired := strings.Contains(validateTag, "required")

		if isRequired && field.IsZero() {
			return fmt.Errorf("field '%s' is required", fieldType.Name)
		}
	}
	return nil
}
