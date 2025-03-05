package defaultz

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// StringDefaulter is a defaulter for string fields.
// The value is set as is.
type StringDefaulter struct{}

var _ Defaulter = &StringDefaulter{}

func (s *StringDefaulter) Name() string {
	return "defaultz.StringDefaulter"
}

func (s *StringDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.String}
}

//nolint:lll
func (s *StringDefaulter) HandleField(value string, _ string, _ reflect.StructField, fieldValue reflect.Value) (bool, bool, *Error) {
	// Handle pointer cases
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem())) // Allocate new string pointer
		}
		fieldValue.Elem().SetString(value) // Set the actual string value
	} else {
		fieldValue.SetString(value) // Direct string assignment
	}

	return true, true, nil
}

type BoolDefaulter struct{}

var _ Defaulter = &BoolDefaulter{}

func (b *BoolDefaulter) Name() string {
	return "defaultz.BoolDefaulter"
}

func (b *BoolDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Bool}
}

//nolint:lll
func (b *BoolDefaulter) HandleField(value string, path string, field reflect.StructField, fieldValue reflect.Value) (bool, bool, *Error) {
	var valueToSet bool
	switch {
	case value == "true":
		valueToSet = true
	case value == "false":
		valueToSet = false
	default:
		return true, false, NewError(b, ErrInvalidDefaultValue, path, field, "invalid boolean value (not 'true' nor 'false')")
	}

	// Handle pointer cases
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem())) // Allocate new bool pointer
		}
		fieldValue.Elem().SetBool(valueToSet) // Set the actual bool value
	} else {
		fieldValue.SetBool(valueToSet) // Direct bool assignment
	}

	return true, true, nil
}

type IntDefaulter struct{}

var _ Defaulter = &IntDefaulter{}

func (i *IntDefaulter) Name() string {
	return "defaultz.IntDefaulter"
}

func (i *IntDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}
}

//nolint:dupl,lll	// we have to use SetXXX() methods for different types, thus can't really get rid of this duplication.
func (i *IntDefaulter) HandleField(value string, path string, field reflect.StructField, fieldValue reflect.Value) (bool, bool, *Error) {
	kind := field.Type.Kind()
	if kind == reflect.Ptr {
		kind = field.Type.Elem().Kind()
	}

	var bitSize int

	//nolint:exhaustive			// there's a default case and we don't want to add meaningless cases
	switch kind {
	case reflect.Int:
		bitSize = 0
	case reflect.Int8:
		bitSize = 8
	case reflect.Int16:
		bitSize = 16
	case reflect.Int32:
		bitSize = 32
	case reflect.Int64:
		bitSize = 64
	default:
		panic(fmt.Sprintf("unsupported integer type: %v", kind))
	}

	intValue, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return true, false, NewError(i, ErrInvalidDefaultValue, path, field, err.Error())
	}

	// Handle pointer cases
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem())) // Allocate new int pointer
		}
		fieldValue.Elem().SetInt(intValue) // Set the actual int value
	} else {
		fieldValue.SetInt(intValue) // Direct int assignment
	}
	return true, true, nil
}

type UintDefaulter struct{}

var _ Defaulter = &UintDefaulter{}

func (u *UintDefaulter) Name() string {
	return "defaultz.UintDefaulter"
}

func (u *UintDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64}
}

//nolint:dupl,lll	// we have to use SetXXX() methods for different types, thus can't really get rid of this duplication.
func (u *UintDefaulter) HandleField(value string, path string, field reflect.StructField, fieldValue reflect.Value) (bool, bool, *Error) {
	kind := field.Type.Kind()
	if kind == reflect.Ptr {
		kind = field.Type.Elem().Kind()
	}

	var bitSize int

	//nolint:exhaustive			// there's a default case and we don't want to add meaningless cases
	switch kind {
	case reflect.Uint:
		bitSize = 0
	case reflect.Uint8:
		bitSize = 8
	case reflect.Uint16:
		bitSize = 16
	case reflect.Uint32:
		bitSize = 32
	case reflect.Uint64:
		bitSize = 64
	default:
		panic(fmt.Sprintf("unsupported unsigned integer type: %v", kind))
	}

	uintValue, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return true, false, NewError(u, ErrInvalidDefaultValue, path, field, err.Error())
	}

	// Handle pointer cases
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem())) // Allocate new uint pointer
		}
		fieldValue.Elem().SetUint(uintValue) // Set the actual uint value
	} else {
		fieldValue.SetUint(uintValue) // Direct uint assignment
	}

	return true, true, nil
}

type FloatDefaulter struct{}

var _ Defaulter = &FloatDefaulter{}

func (f *FloatDefaulter) Name() string {
	return "defaultz.FloatDefaulter"
}

func (f *FloatDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Float32, reflect.Float64}
}

//nolint:lll
func (f *FloatDefaulter) HandleField(value string, path string, field reflect.StructField, fieldValue reflect.Value) (bool, bool, *Error) {
	kind := field.Type.Kind()
	if kind == reflect.Ptr {
		kind = field.Type.Elem().Kind()
	}

	var bitSize int

	//nolint:exhaustive			// there's a default case and we don't want to add meaningless cases
	switch kind {
	case reflect.Float32:
		bitSize = 32
	case reflect.Float64:
		bitSize = 64
	default:
		panic(fmt.Sprintf("unsupported float type: %v", kind))
	}

	floatValue, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return true, false, NewError(f, ErrInvalidDefaultValue, path, field, err.Error())
	}

	// Handle pointer cases
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem())) // Allocate new float pointer
		}
		fieldValue.Elem().SetFloat(floatValue) // Set the actual float value
	} else {
		fieldValue.SetFloat(floatValue) // Direct float assignment
	}

	return true, true, nil
}

type SliceDefaulter struct{}

var _ Defaulter = &SliceDefaulter{}

func (s *SliceDefaulter) Name() string {
	return "defaultz.SliceDefaulter"
}

func (s *SliceDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Slice}
}

//nolint:lll
func (s *SliceDefaulter) HandleField(value string, path string, field reflect.StructField, fieldValue reflect.Value) (bool, bool, *Error) {
	parts := strings.Fields(value) // Split by space
	sliceType := field.Type
	elemType := sliceType.Elem()
	slice := reflect.MakeSlice(sliceType, len(parts), len(parts))
	for j, part := range parts {
		v, err := convertValue(part, elemType)
		if err != nil {
			return true, false, NewError(s, ErrInvalidDefaultValueItem, path, field, err.Error())
		}
		slice.Index(j).Set(v)
	}
	fieldValue.Set(slice)
	return true, true, nil
}

type MapDefaulter struct{}

var _ Defaulter = &MapDefaulter{}

func (m *MapDefaulter) Name() string {
	return "defaultz.MapDefaulter"
}

func (m *MapDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Map}
}

//nolint:lll
func (m *MapDefaulter) HandleField(value string, path string, field reflect.StructField, fieldValue reflect.Value) (bool, bool, *Error) {
	mapInstance := reflect.MakeMap(field.Type)
	pairs := strings.Fields(value) // Split by space

	for _, pair := range pairs {
		//nolint:mnd	// well... pairs have 2 parts
		kv := strings.SplitN(pair, ":", 2)
		//nolint:mnd	// well... pairs have 2 parts
		if len(kv) == 2 {
			// Convert the key to the appropriate type
			keyType := field.Type.Key() // The map's key type
			key, err := convertValue(kv[0], keyType)
			if err != nil {
				return true, false, NewError(m, ErrInvalidDefaultValueKey, path, field, err.Error())
			}

			// Convert the value to the appropriate type
			valueType := field.Type.Elem() // The map's value type
			value, err := convertValue(kv[1], valueType)
			if err != nil {
				return true, false, NewError(m, ErrInvalidDefaultValueItem, path, field, err.Error())
			}

			// Set the key-value pair in the map
			mapInstance.SetMapIndex(key, value)
		}
	}
	fieldValue.Set(mapInstance)
	return true, true, nil
}

type DurationDefaulter struct{}

var _ Defaulter = &DurationDefaulter{}

func (d *DurationDefaulter) Name() string {
	return "defaultz.DurationDefaulter"
}

func (d *DurationDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Int64}
}

//nolint:lll
func (d *DurationDefaulter) HandleField(value string, path string, field reflect.StructField, fieldValue reflect.Value) (bool, bool, *Error) {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return true, false, NewError(d, ErrInvalidDefaultValue, path, field, fmt.Sprintf("invalid duration value: %v", err))
	}

	// Handle pointer cases
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem())) // Allocate new duration pointer
		}
		fieldValue.Elem().Set(reflect.ValueOf(duration).Convert(fieldValue.Type().Elem())) // Set the actual duration value
	} else {
		fieldValue.Set(reflect.ValueOf(duration).Convert(field.Type)) // Direct duration assignment
	}

	return true, true, nil
}

// Converts string values to correct type.
func convertValue(value string, fieldType reflect.Type) (reflect.Value, error) {
	//nolint:exhaustive			// there's a default case and we don't want to add meaningless cases
	switch fieldType.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		return reflect.ValueOf(b).Convert(fieldType), err

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i).Convert(fieldType), err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(u).Convert(fieldType), err

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(f).Convert(fieldType), err

	case reflect.String:
		return reflect.ValueOf(value), nil

	// we don't support:
	// - slice of slices
	// - slice of structs
	// - map of structs
	// - etc.

	default:
		return reflect.Zero(fieldType), fmt.Errorf("unsupported type: %v", fieldType)
	}
}
