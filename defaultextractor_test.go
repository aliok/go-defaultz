package defaultz_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aliok/go-defaultz"
)

// Test struct with various tag formats
type testStruct struct {
	EmptyTag       string `customTag:""`
	WhitespaceTag  string `customTag:" "`
	SimpleValue    string `customTag:"foo"`
	CommaSeparated string `customTag:"a,name=x"`
	SpaceSeparated string `customTag:" b ,name=x"`

	NoDefault         string `customTag:""`
	DefaultAlone      string `customTag:"default="`
	DefaultWithValue  string `customTag:"default=b"`
	DefaultWithSpaces string `customTag:" default=c"`
	MultipleDefaults  string `customTag:"title=foo,default=m n, name=x"`

	PipeSeparated             string `altTag:"key1=value1|default=val|key2=value2"`
	SemicolonDefault          string `altTag:"title=hello;default=xyz;foo=bar"`
	SlashSeparated            string `altTag:"a/b/default=c/d"`
	EmptyDefault              string `altTag:"default=;name=x"`
	MultipleDefaultsSemicolon string `altTag:"default=first;default=second"`

	// Edge Cases
	OnlyDefault      string `customTag:"default=onlyone"`
	DefaultAtStart   string `customTag:"default=start,foo=bar"`
	DefaultAtEnd     string `customTag:"foo=bar,default=end"`
	MissingEqualSign string `customTag:"default"`
	MalformedTag     string `customTag:"default=abc,default"` // Malformed (missing value after second default)
	DoubleSeparators string `altTag:"title=hello;;default=test;;foo=bar"`
}

func TestDefaultzExtractor_ExtractDefault(t *testing.T) {
	tests := []struct {
		fieldName string
		tagName   string
		separator string
		prefix    string
		expected  string
		ok        bool
	}{
		// Basic Cases
		{"EmptyTag", "customTag", ",", "", "", false},
		{"WhitespaceTag", "customTag", ",", "", "", true},
		{"SimpleValue", "customTag", ",", "", "foo", true},
		{"CommaSeparated", "customTag", ",", "", "a", true},
		{"SpaceSeparated", "customTag", ",", "", "b", true},

		// Default Extraction
		{"NoDefault", "customTag", ",", "default=", "", false},
		{"DefaultAlone", "customTag", ",", "default=", "", true},
		{"DefaultWithValue", "customTag", ",", "default=", "b", true},
		{"DefaultWithSpaces", "customTag", ",", "default=", "c", true},
		{"MultipleDefaults", "customTag", ",", "default=", "m n", true},

		// Alternative Separators
		{"PipeSeparated", "altTag", "|", "default=", "val", true},
		{"SemicolonDefault", "altTag", ";", "default=", "xyz", true},
		{"SlashSeparated", "altTag", "/", "default=", "c", true},
		{"EmptyDefault", "altTag", ";", "default=", "", true},
		{"MultipleDefaultsSemicolon", "altTag", ";", "default=", "first", true},

		// Edge Cases
		{"OnlyDefault", "customTag", ",", "default=", "onlyone", true},
		{"DefaultAtStart", "customTag", ",", "default=", "start", true},
		{"DefaultAtEnd", "customTag", ",", "default=", "end", true},
		{"MissingEqualSign", "customTag", ",", "default=", "", false}, // No '=' means invalid default
		{"MalformedTag", "customTag", ",", "default=", "abc", true},   // Should still extract 'abc'
		{"DoubleSeparators", "altTag", ";", "default=", "test", true}, // Handles double separators gracefully
	}

	testType := reflect.TypeOf(testStruct{})
	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			extractor := defaultz.NewDefaultzExtractor(tt.tagName, tt.prefix, tt.separator)

			field, _ := testType.FieldByName(tt.fieldName) // Get struct field

			result, ok, err := extractor.ExtractDefault(field)
			assert.NoError(t, err)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.expected, result)
		})
	}
}
