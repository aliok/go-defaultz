package defaultz

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/go-multierror"
)

// extractorInstance is the DefaultExtractor instance used by default.
var extractorInstance = NewDefaultzExtractor("default", "", ",")

var instance = NewDefaulterRegistry(
	WithBasicDefaulters(),
	WithDefaultExtractor(extractorInstance),
)

const PriorityPrimitiveDefaulter = 1000
const PriorityOtherDefaulter = 2000

// ApplyDefaults applies default values to the struct using the basic defaulters.
func ApplyDefaults(obj interface{}) error {
	return instance.ApplyDefaults(obj)
}

// Defaulter defines an interface for setting default values based on kind.
type Defaulter interface {
	// Name returns the name of the defaulter, which is used for logging and error reporting purposes.
	// It is encouraged to return the name of the defaulter's type with the package name.
	// For example, the name of the basic StringDefaulter type is "defaultz.StringDefaulter".
	Name() string

	// HandledKinds returns the kinds of fields that the defaulter can handle.
	HandledKinds() []reflect.Kind

	// HandleField sets the default value for the field.
	//
	// Parameters:
	//   value: The value to be used for setting the field.
	//   path: The path of the field in the struct. This is useful for error reporting.
	//   field: The field being set.
	//   fieldValue: The value of the field to be set.
	//
	// Returns:
	//   callNext (bool): True if the next defaulter should be called, false otherwise.
	//   set (bool): True if a value is set, false otherwise. The defaultz package will return an error
	//               if no defaulter has set a value and there are errors.
	//   err (*defaultz.Error): An error if the defaulter fails to set the value.
	// /nolint:lll
	HandleField(value string, path string, field reflect.StructField, fieldValue reflect.Value) (callNext bool, set bool, err error)
}

// DefaulterWithPriority is a wrapper for Defaulter with a priority.
type DefaulterWithPriority struct {
	// Defaulter is the defaulter to be used.
	Defaulter Defaulter

	// Priority is the priority of the defaulter. Lower values have higher priority.
	Priority int
}

// DefaulterRegistry defines an interface for managing defaulters.
type DefaulterRegistry interface {
	Register(priority int, defaulter Defaulter) DefaulterRegistry
	ApplyDefaults(obj interface{}) error
}

// defaulterRegistry manages registered defaulters for different kinds.
type defaulterRegistry struct {
	extractor  DefaultExtractor
	defaulters map[reflect.Kind][]DefaulterWithPriority

	// IgnoreCannotSet is a flag to ignore fields that cannot be set.
	ignoreCannotSet bool
}

// compile-time check for interface implementation.
var _ DefaulterRegistry = &defaulterRegistry{}

// NewDefaulterRegistry creates a new defaulterRegistry with optional configurations.
func NewDefaulterRegistry(options ...DefaulterRegistryOption) DefaulterRegistry {
	dr := &defaulterRegistry{
		extractor:  nil,
		defaulters: make(map[reflect.Kind][]DefaulterWithPriority),
	}
	for _, option := range options {
		option(dr)
	}

	// sort defaulters by priority
	for kind := range dr.defaulters {
		sortDefaulters(dr.defaulters[kind])
	}

	return dr
}

func sortDefaulters(dwps []DefaulterWithPriority) {
	// in-place sort by priority, using go's sort package
	// we use a stable sort to keep the order of defaulters with the same priority
	sort.SliceStable(dwps, func(i, j int) bool {
		return dwps[i].Priority < dwps[j].Priority
	})
}

// Register adds a defaulter to the registry with its priority.
// The defaulters are sorted by priority and called in that order.
// If a defaulter denotes that the next defaulter should not be called, the process will stop.
func (r *defaulterRegistry) Register(priority int, defaulter Defaulter) DefaulterRegistry {
	for _, kind := range defaulter.HandledKinds() {
		r.defaulters[kind] = append(r.defaulters[kind], DefaulterWithPriority{Defaulter: defaulter, Priority: priority})
		// sort
		sortDefaulters(r.defaulters[kind])
	}
	return r
}

// DefaulterRegistryOption represents functional options for configuring defaulterRegistry.
type DefaulterRegistryOption func(r *defaulterRegistry)

// - [DurationDefaulter] - priority 2000.
func WithBasicDefaulters() DefaulterRegistryOption {
	return func(r *defaulterRegistry) {
		r.Register(PriorityPrimitiveDefaulter, &BoolDefaulter{})
		r.Register(PriorityPrimitiveDefaulter, &IntDefaulter{})
		r.Register(PriorityPrimitiveDefaulter, &UintDefaulter{})
		r.Register(PriorityPrimitiveDefaulter, &FloatDefaulter{})
		r.Register(PriorityPrimitiveDefaulter, &SliceDefaulter{})
		r.Register(PriorityPrimitiveDefaulter, &MapDefaulter{})
		r.Register(PriorityPrimitiveDefaulter, &StringDefaulter{})

		// some additional defaulters for non-primitive types

		// time.Duration is an int64 under the hood, so we need to handle it separately.
		// should run after IntDefaulter as we want non-durations to be handled by IntDefaulter first
		r.Register(PriorityOtherDefaulter, &DurationDefaulter{})
	}
}

// WithDefaultExtractor sets the default extractor for the registry.
// See [DefaultExtractor] for more information.
func WithDefaultExtractor(extractor DefaultExtractor) DefaulterRegistryOption {
	return func(r *defaulterRegistry) {
		r.extractor = extractor
	}
}

// WithIgnoreCannotSet sets the flag to ignore fields that cannot be set.
// This is useful when the struct has fields that cannot be set, such as unexported fields.
// See https://golang.org/pkg/reflect/#Value.CanSet for more information.
func WithIgnoreCannotSet(ignore bool) DefaulterRegistryOption {
	return func(r *defaulterRegistry) {
		r.ignoreCannotSet = ignore
	}
}

// ApplyDefaults applies default values to the struct.
func (r *defaulterRegistry) ApplyDefaults(obj interface{}) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("object must be a pointer to a struct")
	}

	// check if the type definition allows cycles
	if detectPotentialCycles(val.Elem().Type(), make(map[reflect.Type]bool)) {
		return errors.New("type definition must not have cycles")
	}

	var path string
	typeName := val.Elem().Type().Name()
	if typeName == "" {
		path = "<root>"
	} else {
		path = fmt.Sprintf("%s.(%s)", val.Elem().Type().PkgPath(), typeName)
	}

	return r.DoApplyDefaults(val.Elem(), path)
}

// TODO: can we extract a function to call in the for loop?
//
//nolint:gocognit,funlen 	// we can't extract to a function/method because of the error handling
func (r *defaulterRegistry) DoApplyDefaults(value reflect.Value, path string) error {
	if r.extractor == nil {
		return errors.New("default extractor is not set")
	}

	if len(r.defaulters) == 0 {
		return errors.New("no defaulters are registered")
	}

	// dereference pointer
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// we only handle structs
	if value.Kind() != reflect.Struct {
		return nil
	}

	fieldType := value.Type()
	for i := range value.NumField() {
		field := fieldType.Field(i)
		fieldValue := value.Field(i)

		// Handle nested struct (including pointers to structs)
		if fieldValue.Kind() == reflect.Struct {
			if err := r.DoApplyDefaults(fieldValue, addFieldToPath(path, field)); err != nil {
				return err
			}
			continue
		} else if fieldValue.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			// Initialize pointer to struct if nil
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(field.Type.Elem()))
			}
			if err := r.DoApplyDefaults(fieldValue.Elem(), addFieldToPath(path, field)); err != nil {
				return err
			}
			continue
		}

		if fieldValue.IsValid() && !fieldValue.IsZero() {
			// we do not overwrite non-zero values
			continue
		}

		defaultStr, found, err := r.extractor.ExtractDefault(field)
		if err != nil {
			return NewError(nil, ErrCannotExtractDefault, path, field, err.Error())
		}
		if !found {
			continue
		}

		// we don't allow pointers to pointers
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Ptr {
			return NewError(nil, ErrNotSupported, path, field, "pointer to pointer is not allowed")
		}

		kind := field.Type.Kind()
		if kind == reflect.Ptr {
			kind = field.Type.Elem().Kind()
		}
		//nolint:nestif // can't extract to a function/method because of the error handling
		if defaulters, ok := r.defaulters[kind]; ok {
			if !fieldValue.CanSet() {
				if r.ignoreCannotSet {
					continue
				}
				return NewError(nil, ErrCannotSetField, path, field, "cannot set field")
			}

			var result *multierror.Error
			var somethingSet bool
			for _, defaulterWithPriority := range defaulters {
				defaulter := defaulterWithPriority.Defaulter

				var callNext bool
				var set bool
				callNext, set, err = defaulter.HandleField(defaultStr, path, field, fieldValue)
				// err is always nil for the existing defaulters. May not be nil for custom defaulters.
				if err != nil {
					result = multierror.Append(result, err)
					// we continue to the next defaulter
				}
				if set {
					somethingSet = true
				}
				if !callNext {
					break
				}
			}
			// if there's nothing set and there are errors, return an error
			if !somethingSet && result != nil {
				if result.Len() == 1 {
					return fmt.Errorf("failed to apply default value : %w", result.Errors[0])
				}
				return fmt.Errorf("failed to apply default value: %w", result)
			}
		} else {
			return NewError(nil, ErrNotSupported, path, field, fmt.Sprintf("no defaulters found for kind '%s'", kind))
		}
	}

	return nil
}

func addFieldToPath(path string, field reflect.StructField) string {
	return path + "." + field.Name
}
