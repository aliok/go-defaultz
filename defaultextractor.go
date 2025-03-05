package defaultz

import (
	"reflect"
	"strings"
)

// DefaultExtractor is an interface that defines a method to extract the default value specified in the tag of a struct
// field.
type DefaultExtractor interface {

	// ExtractDefault extracts the default value as string from the tag of a struct field.
	ExtractDefault(field reflect.StructField) (defaultStr string, found bool, err error)
}

var _ DefaultExtractor = &DefaultzExtractor{}

// DefaultzExtractor is a DefaultExtractor implementation that extracts the default value from the tag of a struct
// field.
//
// The tag name, prefix and separator can be configured.
//
// For example, if the tag name is "default", the prefix is "value=" and the separator is ",":
//
// - `default:"value=foo, name=bar, title=baz"` will yield "foo"
//
// - `default:"name=bar, value=foo, title=baz"` will also yield "bar"
//
// This implementation allows using existing struct tags for defaulting purposes.
// For example, if you are using `jsonschema` tag to generate JSON schema with default values, you can reuse the same
// tag for defaulting.
//
// `jsonschema:"title=myfield,default=foo"` will yield "foo", if the tag name is "jsonschema", the prefix is "default="
// and the separator is ",".
//
//nolint:revive		// Extractor would be a too generic name. Here, we're extracting default values.
type DefaultzExtractor struct {

	// TagName is the tag name to be used for extracting default values.
	// For example, for this struct:
	//
	//	  type MyStruct struct {
	//		   Field string `mytag:"hello"`
	//		 }
	//
	// the tagName should be set to "mytag" to be able to extract the default value "hello".
	TagName string

	// Prefix sets the prefix for the default value in the tag.
	//
	// For example, for this struct:
	//
	//	  type MyStruct struct {
	//		   Field string `default:"default=hello,name=myfield"`
	//		 }
	//
	// the tagStringPrefix should be set to "default=" to be able to extract the default value "hello".
	//
	// The tag will be split by the separator (comma in this case) and the prefix will be used to extract the
	// default value.
	//
	// In the following example, the tagStringPrefix should be set to "":
	//
	//	  type MyStruct struct {
	//		   Field string `default:"hello,name=myfield"`
	//		 }
	//
	// The default value to set will be "hello".
	//
	// This option especially useful if you already have some default values defined in your existing tags.
	Prefix string

	// Separator is the separator to be used for splitting the tag value.
	//
	// For example, for this struct:
	//
	//	  type MyStruct struct {
	//		   Field string `default:"hello,name=myfield"`
	//		 }
	//
	// the separator should be set to "," to be able to extract the default value "hello".
	Separator string
}

func NewDefaultzExtractor(tagName, prefix, separator string) DefaultExtractor {
	return &DefaultzExtractor{
		TagName:   tagName,
		Prefix:    prefix,
		Separator: separator,
	}
}

func (d DefaultzExtractor) ExtractDefault(field reflect.StructField) (string, bool, error) {
	tag, ok := field.Tag.Lookup(d.TagName)
	if !ok {
		return "", false, nil
	}

	if tag == "" {
		return "", false, nil
	}

	// split the tag value by separator
	tagParts := strings.Split(tag, d.Separator)
	for _, tagPart := range tagParts {
		tagPart = strings.TrimSpace(tagPart)
		if strings.HasPrefix(tagPart, d.Prefix) {
			return strings.TrimPrefix(tagPart, d.Prefix), true, nil
		}
	}

	return "", false, nil
}
