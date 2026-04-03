package validate

import (
	"fmt"
	"reflect"
)

func ValidateStruct[T any](item *T) error {
	itemValue := reflect.ValueOf(item).Elem()
	itemType := itemValue.Type()

	for i := 0; i < itemValue.NumField(); i++ {
		field := itemValue.Field(i)
		fieldType := itemType.Field(i)

		if fieldType.Tag.Get("validate") == "required" && field.IsZero() {
			return fmt.Errorf("field '%s' is required", fieldType.Name)
		}
	}
	return nil
}
