package defaultz_test

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aliok/go-defaultz"
	"github.com/aliok/go-defaultz/testtypes"
)

func TestApplyDefaultsWithBasicDefaulters(t *testing.T) {
	tests := []struct {
		name       string
		options    []defaultz.DefaulterRegistryOption
		obj        interface{}
		expectJSON string
	}{
		{
			name: "Primitives",
			obj: &struct {
				BoolField1   bool    `default:"true"`
				BoolField2   bool    `default:"false"`
				IntField1    int     `default:"0"`
				IntField2    int8    `default:"123"`
				IntField3    int16   `default:"123"`
				IntField4    int32   `default:"-123"`
				IntField5    int64   `default:"1000000000000000000"` // 17 zeroes
				UintField1   uint    `default:"0"`
				UintField2   uint8   `default:"123"`
				UintField3   uint16  `default:"123"`
				UintField4   uint32  `default:"123"`
				UintField5   uint64  `default:"1000000000000000000"` // 17 zeroes
				FloatField1  float32 `default:"0.0"`
				FloatField2  float32 `default:"9"`
				FloatField3  float64 `default:"-0"`
				FloatField4  float64 `default:"123456789.123456789"`
				StringField1 string  `default:"defaultValue"`
				StringField2 string  `default:"default Value"`
				StringField3 string  `default:""`
			}{},
			expectJSON: `{
				"BoolField1":true,
				"BoolField2":false,
				"IntField1":0,
				"IntField2":123,
				"IntField3":123,
				"IntField4":-123,
				"IntField5":1e+18,
				"UintField1":0,
				"UintField2":123,
				"UintField3":123,
				"UintField4":123,
				"UintField5":1e+18,
				"FloatField1":0, 
				"FloatField2":9, 
				"FloatField3":-0, 
				"FloatField4":1.2345678912345679e+08, 
				"StringField1":"defaultValue",
				"StringField2":"default Value",
				"StringField3":""
			}`,
		},
		{
			name: "Empty primitives",
			obj: &struct {
				BoolField1   bool    `default:""`
				IntField1    int     `default:""`
				UintField1   uint    `default:""`
				FloatField1  float32 `default:""`
				StringField1 string  `default:""`
			}{},
			expectJSON: `{
				"BoolField1":false,
				"IntField1":0,
				"UintField1":0,
				"FloatField1":0, 
				"StringField1":""
			}`,
		},
		{
			name: "Additional types",
			obj: &struct {
				Duration1 time.Duration `default:"1"`
				Duration2 time.Duration `default:"1s"`
				Duration3 time.Duration `default:"1h30m20s999ms"`
			}{},
			expectJSON: `{
				"Duration1":1,
				"Duration2":1e+09,
				"Duration3":5.420999e+12
			}`,
		},
		{
			name: "Slices",
			obj: &struct {
				BoolSlice1 []bool  `default:"true"`
				BoolSlice2 []bool  `default:"true false"`
				IntSlice1  []int   `default:"0"`
				IntSlice2  []int   `default:"0 1"`
				IntSlice3  []int8  `default:"0"`
				IntSlice4  []int8  `default:"0 1"`
				IntSlice5  []int16 `default:"0"`
				IntSlice6  []int16 `default:"0 1"`
				IntSlice7  []int32 `default:"0"`
				IntSlice8  []int32 `default:"0 1"`
				IntSlice9  []int64 `default:"0"`
				IntSlice10 []int64 `default:"0 1"`
				UintSlice1 []uint  `default:"0"`
				UintSlice2 []uint  `default:"0 1"`
				// Go's encoding/json package automatically encodes []uint8 (which is an alias for []byte)
				// as a Base64-encoded string when serializing to JSON.
				// Thus, we can't really check the default value by converting the object to JSON.
				// But it will work!
				// UintSlice3   []uint8   `default:"0"`
				// UintSlice4   []uint8   `default:"0 1"`
				UintSlice5   []uint16  `default:"0"`
				UintSlice6   []uint16  `default:"0 1"`
				UintSlice7   []uint32  `default:"0"`
				UintSlice8   []uint32  `default:"0 1"`
				UintSlice9   []uint64  `default:"0"`
				UintSlice10  []uint64  `default:"0 1"`
				FloatSlice1  []float32 `default:"0.0"`
				FloatSlice2  []float32 `default:"0.0 1.2345"`
				FloatSlice3  []float64 `default:"0.0"`
				FloatSlice4  []float64 `default:"0.0 1.2345"`
				StringSlice1 []string  `default:"defaultValue"`
				StringSlice2 []string  `default:"defaultValue1 defaultValue2"`
			}{},
			expectJSON: `{
				"BoolSlice1":[true],
				"BoolSlice2":[true,false],
				"IntSlice1":[0],
				"IntSlice2":[0,1],
				"IntSlice3":[0],
				"IntSlice4":[0,1],
				"IntSlice5":[0],
				"IntSlice6":[0,1],
				"IntSlice7":[0],
				"IntSlice8":[0,1],
				"IntSlice9":[0],
				"IntSlice10":[0,1],
				"UintSlice1":[0],
				"UintSlice2":[0,1],
				"UintSlice5":[0],
				"UintSlice6":[0,1],
				"UintSlice7":[0],
				"UintSlice8":[0,1],
				"UintSlice9":[0],
				"UintSlice10":[0,1],
				"FloatSlice1":[0],
				"FloatSlice2":[0,1.2345],
				"FloatSlice3":[0],
				"FloatSlice4":[0,1.2345],
				"StringSlice1":["defaultValue"],
				"StringSlice2":["defaultValue1","defaultValue2"]
			}`,
		},
		{
			name: "Empty Slices",
			obj: &struct {
				BoolSlice1 []bool  `default:""`
				IntSlice1  []int   `default:""`
				IntSlice2  []int8  `default:""`
				IntSlice3  []int16 `default:""`
				IntSlice4  []int32 `default:""`
				IntSlice5  []int64 `default:""`
				UintSlice1 []uint  `default:""`
				// Go's encoding/json package automatically encodes []uint8 (which is an alias for []byte)
				// as a Base64-encoded string when serializing to JSON.
				// Thus, we can't really check the default value by converting the object to JSON.
				// But it will work!
				// UintSlice2   []uint8   `default:""`
				UintSlice3   []uint16  `default:""`
				UintSlice4   []uint32  `default:""`
				UintSlice5   []uint64  `default:""`
				FloatSlice1  []float32 `default:""`
				FloatSlice2  []float64 `default:""`
				StringSlice1 []string  `default:""`
			}{},
			expectJSON: `{
				"BoolSlice1":null,
				"IntSlice1":null,
				"IntSlice2":null,
				"IntSlice3":null,
				"IntSlice4":null,
				"IntSlice5":null,
				"UintSlice1":null,
				"UintSlice3":null,
				"UintSlice4":null,
				"UintSlice5":null,
				"FloatSlice1":null,
				"FloatSlice2":null,
				"StringSlice1":null
			}`,
		},
		{
			name: "Maps",
			obj: &struct {
				BoolMap1 map[string]bool `default:"a:true b:false"`
				BoolMap2 map[int]bool    `default:"3:true 4:false"`
				// Go's encoding/json package cannot encode map keys of type float64.
				// The default value will be set correctly, but we can't check it by converting the object to JSON.
				// BoolMap3 map[float64]bool `default:"123456.789:true -987654:false"`
				//
				IntMap1 map[string]int `default:"a:0 b:1"`
				IntMap2 map[int]int    `default:"3:0 4:1"`
				//
				UintMap1 map[string]uint `default:"a:0 b:1"`
				UintMap2 map[int]uint    `default:"3:0 4:1"`
				//
				FloatMap1 map[string]float32 `default:"a:0.0 b:1.2345"`
				FloatMap2 map[int]float32    `default:"3:0.0 4:1.2345"`
				//
				StringMap1 map[string]string `default:"a:defaultValue1 b:defaultValue2"`
				StringMap2 map[int]string    `default:"3:defaultValue1 4:defaultValue2"`
			}{},
			expectJSON: `{
				"BoolMap1":{"a":true,"b":false},
				"BoolMap2":{"3":true,"4":false},
				"IntMap1":{"a":0,"b":1},
				"IntMap2":{"3":0,"4":1},
				"UintMap1":{"a":0,"b":1},
				"UintMap2":{"3":0,"4":1},
				"FloatMap1":{"a":0,"b":1.2345},
				"FloatMap2":{"3":0,"4":1.2345},
				"StringMap1":{"a":"defaultValue1","b":"defaultValue2"},
				"StringMap2":{"3":"defaultValue1","4":"defaultValue2"}
			}`,
		},
		{
			name: "Empty Maps",
			obj: &struct {
				BoolMap1   map[string]bool    `default:""`
				IntMap1    map[string]int     `default:""`
				UintMap1   map[string]uint    `default:""`
				FloatMap1  map[string]float32 `default:""`
				StringMap1 map[string]string  `default:""`
			}{},
			expectJSON: `{
				"BoolMap1":null,
				"IntMap1":null,
				"UintMap1":null,
				"FloatMap1":null,
				"StringMap1":null
			}`,
		},
		{
			name: "Primitive pointers",
			obj: &struct {
				BoolField1   *bool    `default:"true"`
				BoolField2   *bool    `default:"false"`
				IntField1    *int     `default:"0"`
				IntField2    *int8    `default:"123"`
				IntField3    *int16   `default:"123"`
				IntField4    *int32   `default:"-123"`
				IntField5    *int64   `default:"1000000000000000000"` // 17 zeroes
				UintField1   *uint    `default:"0"`
				UintField2   *uint8   `default:"123"`
				UintField3   *uint16  `default:"123"`
				UintField4   *uint32  `default:"123"`
				UintField5   *uint64  `default:"1000000000000000000"` // 17 zeroes
				FloatField1  *float32 `default:"0.0"`
				FloatField2  *float32 `default:"9"`
				FloatField3  *float64 `default:"-0"`
				FloatField4  *float64 `default:"123456789.123456789"`
				StringField1 *string  `default:"defaultValue"`
				StringField2 *string  `default:"default Value"`
			}{},
			expectJSON: `{
				"BoolField1":true,
				"BoolField2":false,
				"IntField1":0,
				"IntField2":123,
				"IntField3":123,
				"IntField4":-123,
				"IntField5":1e+18,
				"UintField1":0,
				"UintField2":123,
				"UintField3":123,
				"UintField4":123,
				"UintField5":1e+18,
				"FloatField1":0,
				"FloatField2":9,
				"FloatField3":-0,
				"FloatField4":1.2345678912345679e+08,
				"StringField1":"defaultValue",
				"StringField2":"default Value"
			}`,
		},
		{
			name: "Additional type pointers",
			obj: &struct {
				Duration1 *time.Duration `default:"1"`
				Duration2 *time.Duration `default:"1s"`
				Duration3 *time.Duration `default:"1h30m20s999ms"`
			}{},
			expectJSON: `{
				"Duration1":1,
				"Duration2":1e+09,
				"Duration3":5.420999e+12
			}`,
		},
		{
			name: "No tag",
			obj: &struct {
				BoolField1  bool   `foo:"true"`
				StringField string `bar:"defaultValue"`
			}{},
			expectJSON: `{
				"BoolField1":false,
				"StringField":""
			}`,
		},
		{
			name: "Nested structs",
			obj: &struct {
				Field1 bool     `default:"true"`
				Field2 []string `default:"val1 val2"`
				Child  struct {
					Field3     int `default:"123"`
					GrandChild struct {
						Field4          map[string]float32 `default:"a:3.2 b:-214.11"`
						GrandGrandChild struct {
							Field5 []int `default:"1 2 3"`
						}
					}
				}
			}{},
			expectJSON: `{
			  "Field1" : true,
			  "Field2" : [ "val1", "val2" ],
			  "Child" : {
				"Field3" : 123,
				"GrandChild" : {
				  "Field4" : {
					"a" : 3.2,
					"b" : -214.11
				  },
				  "GrandGrandChild" : {
					"Field5" : [ 1, 2, 3 ]
				  }
				}
			  }
			}`,
		},
		{
			name: "Slice of structs",
			obj: &struct {
				Field []struct {
					Foo string `default:"bar"`
				}
			}{},
			expectJSON: `{
				"Field": null
			}`,
		},
		{
			name: "Array of structs",
			obj: &struct {
				Field [2]struct {
					Foo string `default:"bar"`
				}
			}{},
			expectJSON: `{
				"Field": [
					{"Foo": ""},
					{"Foo": ""}
				]
			}`,
		},
		{
			name: "Map of struct values",
			obj: &struct {
				Field map[string]struct {
					Foo string `default:"bar"`
				}
			}{},
			expectJSON: `{
				"Field": null
			}`,
		},
		{
			name: "Should not overwrite existing values",
			obj: &struct {
				// define values for these
				Field1 string             `default:"foo"`
				Field2 []int              `default:"1 2 3"`
				Field3 map[string]float32 `default:"a:3.2 b:-214.11"`

				// define no values for these
				Field4 string        `default:"foo"`
				Field5 time.Duration `default:"1s"`
			}{
				Field1: "bar",
				Field2: []int{4, 5, 6},
				Field3: map[string]float32{"c": 1.1, "d": 2.2},
			},
			expectJSON: `{
				"Field1":"bar",
				"Field2":[4,5,6],
				"Field3":{"c":1.1,"d":2.2},
				"Field4":"foo",
				"Field5":1000000000
			}`,
		},
		{
			name: "Should not overwrite existing values - nested struct",
			obj: &struct {
				Field1 string `default:"foo"`
				Nested struct {
					Field2 string `default:"bar"`
					Field3 string `default:"tmp"`
				}
			}{
				Field1: "baz",
				Nested: struct {
					Field2 string `default:"bar"`
					Field3 string `default:"tmp"`
				}{
					Field2: "qux",
					// Field3 is going to be defaulted
				},
			},
			expectJSON: `{
				"Field1":"baz",
				"Nested":{
					"Field2":"qux",
					"Field3":"tmp"
				}
			}`,
		},
		{
			name: "Use different tag",
			obj: &struct {
				BoolField1  bool   `foo:"true"`
				StringField string `bar:"defaultValue"`
			}{},
			options: []defaultz.DefaulterRegistryOption{
				defaultz.WithDefaultExtractor(
					defaultz.NewDefaultzExtractor("foo", "", ","),
				),
			},
			expectJSON: `{
				"BoolField1":true,
				"StringField":""
			}`,
		},
		{
			name: "Use different tag and prefix",
			obj: &struct {
				BoolField1  bool   `existingTag:"default=true"`
				StringField string `existingTag:"default=defaultValue"`
			}{},
			options: []defaultz.DefaulterRegistryOption{
				defaultz.WithDefaultExtractor(
					defaultz.NewDefaultzExtractor("existingTag", "default=", ","),
				),
			},
			expectJSON: `{
				"BoolField1":true,
				"StringField":"defaultValue"
			}`,
		},
		{
			name: "No prefix in the tag",
			obj: &struct {
				BoolField1  bool   `existingTag:""`
				IntField1   int    `existingTag:"name=IntField1"`
				UintField1  uint   `existingTag:"foo"`
				StringField string `existingTag:"what=hasDefault"`
			}{},
			options: []defaultz.DefaulterRegistryOption{
				defaultz.WithDefaultExtractor(
					defaultz.NewDefaultzExtractor("existingTag", "default=", ","),
				),
			},
			expectJSON: `{
				"BoolField1":false,
				"IntField1":0,
				"UintField1":0,
				"StringField":""
			}`,
		},
		{
			name: "Additional info in the tag - no prefix",
			obj: &struct {
				BoolField1   bool   `default:"true,foo=bar"`
				IntField1    int    `default:"123,foo=bar,baz=qux"`
				UintField1   uint   `default:"456,aaa,foo=bar"`
				StringField1 string `default:",foo=bar"`
				StringField2 string `default:"   hasDefault,foo=bar   "`
			}{},
			options: []defaultz.DefaulterRegistryOption{
				defaultz.WithDefaultExtractor(
					defaultz.NewDefaultzExtractor("default", "", ","),
				),
			},
			expectJSON: `{
				"BoolField1":true,
				"IntField1":123,
				"UintField1":456,
				"StringField1":"",
				"StringField2":"hasDefault"
			}`,
		},
		{
			name: "Additional info in the tag - with prefix",
			obj: &struct {
				BoolField1  bool   `default:"default=true,foo=bar"`
				IntField1   int    `default:"foo=bar,default=123,baz=qux"`
				UintField1  uint   `default:"foo=bar,baz=qux,default=123"`
				StringField string `default:"   default=hasDefault,foo=bar   "`
			}{},
			options: []defaultz.DefaulterRegistryOption{
				defaultz.WithDefaultExtractor(
					defaultz.NewDefaultzExtractor("default", "default=", ","),
				),
			},
			expectJSON: `{
				"BoolField1":true,
				"IntField1":123,
				"UintField1":123,
				"StringField":"hasDefault"
			}`,
		},
		{
			name: "Custom separator - no prefix",
			obj: &struct {
				BoolField1   bool   `default:"true#foo=bar"`
				IntField1    int    `default:"123#foo=bar#baz=qux"`
				UintField1   uint   `default:"456#aaa#foo=bar"`
				StringField1 string `default:"#foo=bar"`
				StringField2 string `default:"   hasDefault#foo=bar   "`
			}{},
			options: []defaultz.DefaulterRegistryOption{
				defaultz.WithDefaultExtractor(
					defaultz.NewDefaultzExtractor("default", "", "#"),
				),
			},
			expectJSON: `{
				"BoolField1":true,
				"IntField1":123,
				"UintField1":456,
				"StringField1":"",
				"StringField2":"hasDefault"
			}`,
		},
		{
			name: "Custom separator - with prefix",
			obj: &struct {
				BoolField1  bool   `default:"default=true#foo=bar"`
				IntField1   int    `default:"foo=bar#default=123#baz=qux"`
				UintField1  uint   `default:"foo=bar#baz=qux#default=123"`
				StringField string `default:"   default=hasDefault#foo=bar   "`
			}{},
			options: []defaultz.DefaulterRegistryOption{
				defaultz.WithDefaultExtractor(
					defaultz.NewDefaultzExtractor("default", "default=", "#"),
				),
			},
			expectJSON: `{
				"BoolField1":true,
				"IntField1":123,
				"UintField1":123,
				"StringField":"hasDefault"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := []defaultz.DefaulterRegistryOption{
				defaultz.WithBasicDefaulters(),
				defaultz.WithDefaultExtractor(
					defaultz.NewDefaultzExtractor("default", "", ","),
				),
			}
			if len(tt.options) > 0 {
				options = append(options, tt.options...)
			}
			registry := defaultz.NewDefaulterRegistry(options...)

			derr := registry.ApplyDefaults(tt.obj)
			if derr != nil {
				t.Fatalf("unexpected error: %v", derr)
			}
			require.NoError(t, derr, tt.name)

			jsonBytes, err := json.Marshal(tt.obj)
			require.NoError(t, err, tt.name)
			assert.JSONEq(t, tt.expectJSON, string(jsonBytes), tt.name)
		})
	}
}

func TestApplyDefaultsWithBasicDefaulters_InvalidCases(t *testing.T) {
	tests := []struct {
		name      string
		obj       interface{}
		expectErr string
	}{
		{
			name:      "No struct",
			obj:       "foo",
			expectErr: "object must be a pointer to a struct",
		},
		{
			name: "Struct, but not a pointer",
			obj: struct {
				Foo string `default:"bar"`
			}{},
			expectErr: "object must be a pointer to a struct",
		},
		{
			name: "No pointer to pointer",
			obj: &struct {
				Field **int `default:"3"`
			}{},
			expectErr: "not supported - " +
				"pointer to pointer is not allowed, " +
				"path:'<root>.Field`, " +
				"field:'Field **int `default:\"3\"`'",
		},
		{
			name: "Not convertible to the type",
			obj: &struct {
				Field int `default:"abc"`
			}{},
			expectErr: "failed to apply default value : (defaultz.IntDefaulter): invalid default value - " +
				"strconv.ParseInt: parsing \"abc\": invalid syntax, " +
				"path:'<root>.Field`, " +
				"field:'Field int `default:\"abc\"`'",
		},
		{
			name: "Overflow",
			obj: &struct {
				Field int `default:"99999999999999999999999999999999999999999999999999999999999999999"`
			}{},
			expectErr: "failed to apply default value : (defaultz.IntDefaulter): invalid default value - " +
				"strconv.ParseInt: " +
				"parsing \"99999999999999999999999999999999999999999999999999999999999999999\": value out of range, " +
				"path:'<root>.Field`, " +
				"field:'Field int `default:\"99999999999999999999999999999999999999999999999999999999999999999\"`'",
		},
		{
			name: "Slices of non-primitive types",
			obj: &struct {
				Field []rand.Rand `default:"foo"`
			}{},
			expectErr: "failed to apply default value : (defaultz.SliceDefaulter): invalid default value item - " +
				"unsupported type: rand.Rand, " +
				"path:'<root>.Field`, " +
				"field:'Field []rand.Rand `default:\"foo\"`'",
		},
		{
			// even though we support defaulting fields of type time.Duration, we don't support
			// defaulting slices of time.Duration
			name: "Slices of time.Duration",
			obj: &struct {
				Field []time.Duration `default:"1m 2m"`
			}{},
			expectErr: "failed to apply default value : (defaultz.SliceDefaulter): invalid default value item - " +
				"strconv.ParseInt: parsing \"1m\": invalid syntax, " +
				"path:'<root>.Field`, " +
				"field:'Field []time.Duration `default:\"1m 2m\"`'",
		},
		{
			name: "Maps with keys of non-primitive types",
			obj: &struct {
				Field map[rand.Rand]bool `default:"foo:true"`
			}{},
			expectErr: "failed to apply default value : (defaultz.MapDefaulter): invalid default value key - " +
				"unsupported type: rand.Rand, " +
				"path:'<root>.Field`, " +
				"field:'Field map[rand.Rand]bool `default:\"foo:true\"`'",
		},
		{
			name: "Maps with values of non-primitive types",
			obj: &struct {
				Field map[string]rand.Rand `default:"foo:bar"`
			}{},
			expectErr: "failed to apply default value : (defaultz.MapDefaulter): invalid default value item - " +
				"unsupported type: rand.Rand, " +
				"path:'<root>.Field`, " +
				"field:'Field map[string]rand.Rand `default:\"foo:bar\"`'",
		},
		{
			name: "Slice of structs",
			obj: &struct {
				Field []struct {
					Foo string `default:"bar"`
				} `default:"foo"`
			}{},
			expectErr: "failed to apply default value : (defaultz.SliceDefaulter): invalid default value item - " +
				"unsupported type: struct { Foo string \"default:\\\"bar\\\"\" }, " +
				"path:'<root>.Field`, " +
				"field:'Field []struct { Foo string \"default:\\\"bar\\\"\" } `default:\"foo\"`'",
		},
		{
			name: "Array of structs",
			obj: &struct {
				Field [5]struct {
					Foo string `default:"bar"`
				} `default:"foo"`
			}{},
			expectErr: "not supported - " +
				"no defaulters found for kind 'array', " +
				"path:'<root>.Field`, " +
				"field:'Field [5]struct { Foo string \"default:\\\"bar\\\"\" } `default:\"foo\"`'",
		},
		{
			name: "Map of struct keys",
			obj: &struct {
				Field map[struct {
					Foo string `default:"bar"`
				}]string `default:"foo:bar"`
			}{},
			expectErr: "failed to apply default value : (defaultz.MapDefaulter): invalid default value key - " +
				"unsupported type: struct { Foo string \"default:\\\"bar\\\"\" }, " +
				"path:'<root>.Field`, " +
				"field:'Field map[struct { Foo string \"default:\\\"bar\\\"\" }]string `default:\"foo:bar\"`'",
		},
		{
			name: "Map of struct values",
			obj: &struct {
				Field map[string]struct {
					Foo string `default:"bar"`
				} `default:"foo:bar"`
			}{},
			expectErr: "failed to apply default value : (defaultz.MapDefaulter): invalid default value item - " +
				"unsupported type: struct { Foo string \"default:\\\"bar\\\"\" }, " +
				"path:'<root>.Field`, " +
				"field:'Field map[string]struct { Foo string \"default:\\\"bar\\\"\" } `default:\"foo:bar\"`'",
		},
		{
			name: "Pointer to pointer deep",
			obj: &struct {
				Field struct {
					Field1 **int `default:"123"`
				}
			}{},
			expectErr: "not supported - " +
				"pointer to pointer is not allowed, " +
				"path:'<root>.Field.Field1`, " +
				"field:'Field1 **int `default:\"123\"`'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			assert.NotEmpty(t, tt.expectErr, "Invalid test case: expectErr is empty")
			if err = defaultz.ApplyDefaults(tt.obj); err != nil {
				assert.EqualError(t, err, tt.expectErr, tt.name)
			} else {
				t.Logf("expected error, got nil")
				var jsonBytes []byte
				if jsonBytes, err = json.Marshal(tt.obj); err != nil {
					t.Fatalf("failed to marshal object to json: %v", err)
				}
				t.Logf("got: %s", string(jsonBytes))
				t.Fail()
			}
		})
	}
}

func TestApplyDefaultsWithNamedStructs(t *testing.T) {
	type Embedded1 struct {
		Field4 map[string]float32 `default:"a:3.2 b:-214.11"`
	}

	type Embedded2 struct {
		Field5 []int `default:"1 2 3"`
	}

	type GrandChild struct {
		Field3 string `default:"abc"`
		Embedded1
		*Embedded2
	}

	type Child struct {
		Field2     int `default:"123"`
		GrandChild GrandChild
	}

	type Parent struct {
		Field1 bool `default:"true"`
		Child  Child
	}

	obj := &Parent{}

	derr := defaultz.ApplyDefaults(obj)
	if derr != nil {
		t.Fatalf("unexpected error: %v", derr)
	}
	require.NoError(t, derr)

	assert.True(t, obj.Field1)
	assert.Equal(t, 123, obj.Child.Field2)
	assert.Equal(t, "abc", obj.Child.GrandChild.Field3)
	assert.Equal(t, map[string]float32{"a": 3.2, "b": -214.11}, obj.Child.GrandChild.Field4)
	assert.Equal(t, []int{1, 2, 3}, obj.Child.GrandChild.Field5)
}

type cyclicParent1 struct {
	Field1 bool `default:"true"`
	Child  cyclicChild1
}

type cyclicChild1 struct {
	Field2 int `default:"123"`
	Parent *cyclicParent1
}

func TestApplyDefaultsCyclicReference(t *testing.T) {
	obj := &cyclicParent1{}

	derr := defaultz.ApplyDefaults(obj)
	require.Error(t, derr)
	assert.Equal(t, "type definition must not have cycles", derr.Error())
}

func TestApplyDefaultsUnexportedFields(t *testing.T) {
	obj := &testtypes.TestExportedWithUnexportedField{}

	derr := defaultz.ApplyDefaults(obj)
	require.Error(t, derr)
	assert.Equal(
		t,
		"cannot set field - cannot set field, path:'github.com/aliok/go-defaultz/testtypes."+
			"(TestExportedWithUnexportedField).unexportedField`, field:'unexportedField string `default:\"bar\"`'",
		derr.Error(),
	)
}

func TestApplyDefaultsIgnoreUnexportedFields(t *testing.T) {
	obj := &testtypes.TestExportedWithUnexportedField{}
	d := defaultz.NewDefaulterRegistry(
		defaultz.WithBasicDefaulters(),
		defaultz.WithDefaultExtractor(
			defaultz.NewDefaultzExtractor("default", "", ","),
		),
		defaultz.WithIgnoreCannotSet(true),
	)

	derr := d.ApplyDefaults(obj)
	require.NoError(t, derr)
	assert.Equal(t, "foo", obj.ExportedField)
}

type customDefaulter struct{}

var _ defaultz.Defaulter = customDefaulter{}

func (c customDefaulter) Name() string {
	return "test.customDefaulter"
}

func (c customDefaulter) HandledKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Bool}
}

//nolint:lll
func (c customDefaulter) HandleField(value string, _ string, _ reflect.StructField, fieldValue reflect.Value) (bool, bool, *defaultz.Error) {
	if value == "yay" {
		fieldValue.SetBool(true)
		return false, true, nil
	}
	return true, false, nil
}

func TestApplyDefaultsCustomDefaulter(t *testing.T) {
	// the custom defaulter will be parsing `yay` as true
	obj := &struct {
		Field bool `default:"yay"`
	}{}

	d := defaultz.NewDefaulterRegistry(
		defaultz.WithBasicDefaulters(),
		defaultz.WithDefaultExtractor(
			defaultz.NewDefaultzExtractor("default", "", ","),
		),
	)
	d.Register(2000, customDefaulter{}) // should run after the core BoolDefaulter

	err := d.ApplyDefaults(obj)
	require.NoError(t, err)
	assert.True(t, obj.Field)
}
