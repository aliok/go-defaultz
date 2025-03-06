# go-defaultz

[![Build Status](https://img.shields.io/github/actions/workflow/status/aliok/go-defaultz/unit-tests.yaml?branch=main&style=flat-square)](https://github.com/aliok/go-defaultz/actions/workflows/unit-tests.yaml)
[![codecov](https://codecov.io/github/aliok/go-defaultz/graph/badge.svg?token=VJ8020TZH7)](https://codecov.io/github/aliok/go-defaultz)
[![Go Report Card](https://goreportcard.com/badge/github.com/aliok/go-defaultz)](https://goreportcard.com/report/github.com/aliok/go-defaultz)
[![Go Reference](https://pkg.go.dev/badge/github.com/aliok/go-defaultz.svg)](https://pkg.go.dev/github.com/aliok/go-defaultz)
[![Release](https://img.shields.io/github/v/release/aliok/go-defaultz)](https://github.com/aliok/go-defaultz/releases/latest)
![License](https://img.shields.io/dub/l/vibe-d.svg)




**Defaulting Go structs with field tags!**

## Install

```shell
go get github.com/aliok/go-defaultz
```

**Note:** go-defaultz uses [Go Modules](https://go.dev/wiki/Modules) to manage dependencies.

## What is go-defaultz?

go-defaultz is a library that provides a way to set default values to Go structs with field tags.

- No need to write boilerplate code to set default values.
- Works with nested structs.
- Supports slices and maps of primitive types.
- Works with pointers.
- Supports custom types.
- Supports custom field tag formats.
- Keeps the original value if it is not the zero value.

## Usage

```go
package main

import (
	"fmt"
	"github.com/aliok/go-defaultz"
)

type Config struct {
	Host string `default:"localhost"`
	Port int    `default:"8080"`
}

func main() {
	cfg := Config{}
	_ = defaultz.ApplyDefaults(&cfg)
	fmt.Printf("%+v\n", cfg)
	// Output: {Host:localhost Port:8080}
}
```

See [examples](#examples) for more complex examples.

## Supported field types

- Primitive types: `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `string`, `bool`

```go
  Field1 int    `default:"42"`
  Field2 bool   `default:"true"`
```

- Slices of primitive types: `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`, `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`, `[]float32`, `[]float64`, `[]string`, `[]bool`

```go
  // slice items are separated by space, instead of tediously escaping things for 
  // the https://pkg.go.dev/encoding#TextMarshaler format
  Field3       []string          `default:"a b"`
  Field4       []int64           `default:"1 2"`
````

- Maps with keys or values of primitive types
```go
  // map pairs are also separated by space. key and value are separated by colon.
  Field5       map[string]bool   `default:"a:true b:false"`
  Field6       map[float64]uint8 `default:"1.2:3 4.5:6"`
```

- `time.Duration`, `[]time.Duration`, `map[...]time.Duration`, `map[time.Duration]...`
```go
  Field7       time.Duration     `default:"1s"`
  Field8       []time.Duration   `default:"1s 2m"`
```

- All types via custom defaulters
```go
  MaxContentSize    FileSize  `default:"1MB"`
  Distance          Metric    `default:"1km"`
  Date              MyDate    `default:"2022-01-01"`
```

[//]: # (TODO: map[time.Duration]int and map[int]time.Duration ?)

## Configuring go-defaultz

### Customizing how default values are extracted

By default, go-defaultz uses `default` field tag to extract default values.

You can customize the field tag by setting up the extractor.

```go
type Config struct {
	Host string `othertag:"foo" mytag:"theDefault=localhost,otherInfo=foo"`
}

func main() {
	reg := defaultz.NewDefaulterRegistry(
		// Register the core defaulters.
		defaultz.WithBasicDefaulters(),
		
		// Register the custom default extractor.
		// This will extract the default value from the `mytag` tag: `theDefault=localhost,otherInfo=foo"`
		// It will then split the tag string by `,` and look for the key `theDefault=` prefix: `theDefault=localhost`.
		// Then it will trim the prefix to find out the default value: `localhost`.
		defaultz.WithDefaultExtractor(defaultz.NewDefaultzExtractor("mytag", "theDefault=", ",")),
	)

	cfg := Config{}
	_ = reg.ApplyDefaults(&cfg)
	fmt.Printf("%+v\n", cfg)
	// {Host:localhost}
}
```

Being able to customize the default value extraction is very important, if you want to be able to use existing field tags in your structs. Otherwise, you would end up with the same default value definition in multiple tags.

An example is using `jsonschema` tag for generating JSON schema with [invopop/jsonschema](https://github.com/invopop/jsonschema):

```go
type Config struct {
    Host string `jsonschema:"default=localhost,description=The host name"`
}
```

If your needs are not met by `defaultz.NewDefaultzExtractor()`, you can implement your own `defaultz.DefaultExtractor` interface.

Following example shows how to implement a custom extractor that extracts default values from a field tag in [piglatin](https://en.wikipedia.org/wiki/Pig_Latin) and converts it to English.

```go
type CustomExtractor struct {}

func (c CustomExtractor) ExtractDefault(field reflect.StructField) (defaultStr string, found bool, err error) {
    val, ok := field.Tag.Lookup("piglatin") // val is `igpay`
    if !ok {
        return "", false, nil
    }
    // ...
    // custom logic to extract the default value and convert piglatin to English
    // ...
    return "pig", true, nil
}

type Config struct {
    Host string `piglatin:"igpay"`
}

func main() {
    reg := defaultz.NewDefaulterRegistry(
	    // Register the core defaulters.
        defaultz.WithBasicDefaulters(),
		// Set the custom extractor
        defaultz.WithDefaultExtractor(&CustomExtractor{}),
    )

    cfg := Config{}
    _ = reg.ApplyDefaults(&cfg)
    fmt.Printf("%+v\n", cfg)
    // Output: {Host:pig}
}
```

### Type aliases

Type aliases work out of the box. 

```go
// Age is a type alias for int
type Age int

func (a Age) CanVote() bool {
	return a >= 18
}

type Person struct {
	Age Age `default:"20"`
}

func main() {
	obj := Person{}

	_ = defaultz.ApplyDefaults(&obj)

	fmt.Printf("Age: %v\n", obj.Age)
	fmt.Printf("Can vote: %v\n", obj.Age.CanVote())
	// Output:
	// Age: 20
	// Can vote: true
}
```

### Defaulting custom types

In the previous example, `Age` didn't need a custom parsing or defaulting logic.

However, if you have different parsing or defaulting rules for the type alias, you can implement a custom defaulter.

Consider this type:
```go
import "github.com/dustin/go-humanize"

type FileSize uint64

func (f FileSize) String() string {
	return humanize.Bytes(uint64(f))
}

func ParseFileSize(s string) (FileSize, error) {
	v, err := humanize.ParseBytes(s)
	if err != nil {
		return FileSize(0), err
	}
	return FileSize(v), nil
}
```

It uses `github.com/dustin/go-humanize` to parse and format the `FileSize`.

Before talking about defaulting, let's see how `FileSize` works:

```go
type Config struct {
	MaxContentSize FileSize		// no defaulting tags yet
}

func main() {
	size, _ := ParseFileSize("1 MB")
	obj := Config{
		MaxContentSize: size,
	}

	fmt.Printf("MaxContentSize: %v\n", obj.MaxContentSize)
	// Output:
	// MaxContentSize: 1.0 rMB
}
```

It outputs the size in human-readable format.

Since we would like to support the human-readable format in the defaulting logic, we need to implement a custom defaulter.

The defaulting instruction for the `FileSize` field will be in human-readable format, like `1 MB`.
While at it, let's add a few more fields.
```go
type Config struct {
	MaxContentSize    FileSize  `default:"1 MB"`
	MaxContentSizePtr *FileSize `default:"2 MB"`
	Other             uint64    `default:"3 MB"`        // this shouldn't be set
}
```

Now, let's implement a custom defaulter for `FileSize`.

The defaulter will handle `reflect.Uint64` kind as `FileSize` is an alias for `uint64`.
```go
type FileSizeDefaulter struct{}

func (c FileSizeDefaulter) Name() string {
	return "custom.FileSizeDefaulter"
}

func (c FileSizeDefaulter) HandledKinds() []reflect.Kind {
	// FileSize is a uint64 under the hood.
	return []reflect.Kind{reflect.Uint64}
}
```

The defaulter will parse the human-readable size and set it to the field.
```go
func (c FileSizeDefaulter) HandleField(value string, path string, field reflect.StructField,
	fieldValue reflect.Value) (callNext bool, set bool, err error) {

	zeroFileSize := FileSize(0)
	if field.Type == reflect.TypeOf(FileSize(0)) || field.Type == reflect.TypeOf(&zeroFileSize) {
		// we can handle the field
		// we won't support **FileSize fields, but we wouldn't want to default that with tags anyway.

		fileSize, err := ParseFileSize(value)

		if err != nil {
			// we cannot parse the specified default value, let's leave it to the next defaulter.
			// maybe they know how to handle it.

			// alternatively, we could return an error here, if we want to stop the defaulting process entirely as we know that this is a FileSize field but the default value is not parsable.
			// callNext=true, set=false, err=error
			return true, false, err
		}

		// Handle pointer cases
		if fieldValue.Kind() == reflect.Ptr {
			// if the field is a pointer, we need to allocate a new value and set it
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem())) // Allocate new FileSize pointer
			}
			fieldValue.Elem().Set(reflect.ValueOf(fileSize).Convert(fieldValue.Type().Elem())) // Set the actual FileSize value
		} else {
			// if the field is not a pointer, we can directly set the value
			fieldValue.Set(reflect.ValueOf(fileSize).Convert(field.Type)) // Direct FileSize assignment
		}

		// we stop on purpose, we don't want to call the next defaulter
		// callNext=false, set=true, err=nil
		return false, true, nil
	}

	// we cannot handle the field. it is a uint64 field, but not of FileSize type.
	// let's leave it to the next defaulter.
	// callNext=true, set=false, err=nil
	return true, false, nil
}
````

Let's register the custom defaulter and apply the defaults.

```go
func main() {
	cfg := Config{}

	reg := defaultz.NewDefaulterRegistry(
		defaultz.WithBasicDefaulters(),
		defaultz.WithDefaultExtractor(
			defaultz.NewDefaultzExtractor("default", "", ","),
		),
	)

	// should run after the core UintDefaulter which has a precedence of 1000
	reg.Register(2000, FileSizeDefaulter{})

	_ = reg.ApplyDefaults(&cfg)

	fmt.Printf("cfg.MaxContentSize: %s\n", cfg.MaxContentSize)
	fmt.Printf("cfg.MaxContentSizePtr: %s\n", cfg.MaxContentSizePtr)
	fmt.Printf("cfg.Other: %d\n", cfg.Other)
	// Output:
	// cfg.MaxContentSize: 1.0 MB
	// cfg.MaxContentSizePtr: 2.0 MB
	// cfg.Other: 0
}
```

You can see that:
- `MaxContentSize` and `MaxContentSizePtr` are set to the specified values.
- `Other` is not set as it is not a `FileSize` field.

See [defaultz.Defaulter](https://pkg.go.dev/github.com/aliok/go-defaultz#Defaulter) interface for more information.

## Examples

### Complex example

```go
import (
    "time"
    
    "gopkg.in/yaml.v3"
    "github.com/aliok/go-defaultz"
)

type Parent struct {
	Field1       string            `default:"field1"`
	Field2       *int              `default:"42"`
    // slice items are separated by space, instead of tediously escaping things for 
	// the https://pkg.go.dev/encoding#TextMarshaler format
	Field3       []string          `default:"a b"`
    // map pairs are also separated by space. key and value are separated by colon.
	Field4       map[float64]uint8 `default:"1.2:3 4.5:6"`
	Field5       time.Duration     `default:"1s"`
	Child        Child
	ChildPointer *Child
	Embedded
	Anonymous struct {
		AnonymousField string `default:"a"`
	}
}

type Child struct {
	ChildField string `default:"c"`
}

type Embedded struct {
	EmbeddedField string `default:"e"`
}

func main() {
	obj := Parent{}
	_ = defaultz.ApplyDefaults(&obj)

	// convert to YAML for pretty printing
	out, _ := yaml.Marshal(obj)
	println(string(out))
	// Output:
	//   field1: field1
	//   field2: 42
	//   field3:
	//      - a
	//      - b
	//   field4:
	//      1.2: 3
	//      4.5: 6
	//   field5: 1s
	//   child:
	//      childfield: c
	//   childpointer:
	//      childfield: c
	//   embedded:
	//      embeddedfield: e
	//   anonymous:
	//      anonymousfield: a
}
```

### Reading configuration files

See [aliok/best-go-config-setup](https://github.com/aliok/best-go-config-setup) for a complete example of integrating these features:

- Reading configuration from configuration files, with the ability to override through environment variables and flags. ([spf13/viper](https://github.com/spf13/viper))
- Filling in the default values (defined in tags) not provided in the user passed configuration. ([go-defaultz](https://github.com/aliok/go-defaultz))
- Validating the configuration with the rules defined in the tags. ([go-playground/validator](https://github.com/go-playground/validator/))
- Generating a JSON schema for the configuration file. ([invopop/jsonschema](https://github.com/invopop/jsonschema/))

So that:
- User is allowed to only define the values they want to change.
- Default values and validation rules are defined in a single place.
- User has IDE auto-completion and validation support via JSON schema.

## Comparing with other libraries

| Feature                       | [go-defaultz](https://github.com/aliok/go-defaultz) | [creasty/defaults](https://github.com/creasty/defaults) | [mcuadros/go-defaults](https://github.com/mcuadros/go-defaults) |
|-------------------------------|-----------------------------------------------------|---------------------------------------------------------|-----------------------------------------------------------------|
| Defaulting with field tags    | ✅                                                   | ✅                                                       | ✅                                                               |
| Nested structs                | ✅                                                   | ✅                                                       | ❌                                                               |
| Pointers                      | ✅                                                   | ✅                                                       | ❌                                                               |
| Type aliases                  | ✅                                                   | ✅                                                       | ❌                                                               |
| Preserve non-zero values      | ✅                                                   | ✅                                                       | ❌                                                               |
| Custom defaulting logic       | ✅                                                   | ❌                                                       | ❌                                                               |
| Customizable field tags       | ✅                                                   | ❌                                                       | ❌                                                               |
| Customizable field tag format | ✅                                                   | ❌                                                       | ❌                                                               |
| No-escape necessary in tags   | ✅                                                   | ❌                                                       | ❌                                                               |

## FAQ

#### Why are pointers to pointers not supported?

Because it doesn't make sense to have some default values for them.


## Development

Pre-requisites:
- [Go](https://golang.org/dl/)
- [golangci-lint](https://golangci-lint.run/usage/install/) min version `v1.64.6`
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)

Pre-commit checks:
```shell
go test -race -count=1 ./...

golangci-lint run --fix

goimports -w .
```

Make sure the code coverage is not decreased.
