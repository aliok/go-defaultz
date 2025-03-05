package defaultz

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrCannotSetField is returned when the field cannot be set.
// See [reflect.Value.CanSet] for more information.
var ErrCannotSetField = errors.New("cannot set field")

// ErrCannotExtractDefault is returned when the default value cannot be extracted.
var ErrCannotExtractDefault = errors.New("cannot extract default value")

// ErrInvalidDefaultValue is returned when the default value defined in the struct tag is invalid.
var ErrInvalidDefaultValue = errors.New("invalid default value")

// ErrInvalidDefaultValueItem is returned when one of the items in the default value is invalid.
// This is used when the default value is a slice or a map.
var ErrInvalidDefaultValueItem = errors.New("invalid default value item")

// ErrInvalidDefaultValueKey is returned when the key of a map default value is invalid.
var ErrInvalidDefaultValueKey = errors.New("invalid default value key")

// ErrNotSupported is returned when the operation is not supported.
var ErrNotSupported = errors.New("not supported")

type Error struct {
	Defaulter Defaulter
	Err       error
	FieldPath string
	Field     reflect.StructField
	Msg       string
}

func NewError(defaulter Defaulter, err error, fieldPath string, field reflect.StructField, msg string) *Error {
	if err == nil {
		err = errors.New("unknown error")
	}

	return &Error{
		Defaulter: defaulter,
		Err:       err,
		FieldPath: fieldPath,
		Field:     field,
		Msg:       msg,
	}
}

func (e *Error) Error() string {
	typeName := e.Field.Type.String()

	if e.Defaulter == nil {
		return fmt.Sprintf("%s - %s, "+
			"path:'%s.%s`, "+
			"field:'%s %s `%s`'",
			e.Err.Error(), e.Msg, e.FieldPath, e.Field.Name, e.Field.Name, typeName, e.Field.Tag)
	}
	return fmt.Sprintf("(%s): %s - %s, "+
		"path:'%s.%s`, "+
		"field:'%s %s `%s`'",
		e.Defaulter.Name(), e.Err.Error(), e.Msg, e.FieldPath, e.Field.Name, e.Field.Name, typeName, e.Field.Tag)
}

func (e *Error) Unwrap() error { return e.Err }
