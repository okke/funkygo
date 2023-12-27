package fu

import (
	"fmt"
	"reflect"
)

func Map2Struct[T any](data map[string]any, into T) (T, error) {
	destValue := reflect.ValueOf(into)

	// Check if destStruct is a pointer to a struct
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
		return *new(T), fmt.Errorf("destination is not a pointer to a struct")
	}

	destValue = destValue.Elem()

	for key, value := range data {
		field := destValue.FieldByName(key)

		if !field.IsValid() || !field.CanSet() {
			continue
		}

		fieldType := field.Type()
		mapValueType := reflect.TypeOf(value)

		if !mapValueType.AssignableTo(fieldType) {
			continue
		}

		field.Set(reflect.ValueOf(value))
	}

	return into, nil
}
