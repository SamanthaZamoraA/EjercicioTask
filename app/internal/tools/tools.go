package tools

import (
	"fmt"
)

type FieldError struct {
	Field string
	Msg   string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

// CheckFieldExistance es una función que recibe un mapa de strings y un slice de strings
func CheckFieldExistance(fields map[string]any, requiredFields ...string) (err error) {

	// Iterar sobre el slice de strings
	for _, field := range requiredFields {
		// Comprobar que el campo exista en el mapa de strings
		if _, ok := fields[field]; !ok {
			err = &FieldError{
				Field: field,
				Msg:   "field is required",
			}
			return
		}
	}
	return
}
