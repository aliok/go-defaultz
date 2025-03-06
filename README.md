# go-defaultz

[![Build Status](https://img.shields.io/github/actions/workflow/status/aliok/go-defaultz/unit-tests.yaml?branch=main&style=flat-square)](https://github.com/aliok/go-defaultz/actions/workflows/unit-tests.yaml)
[![codecov](https://codecov.io/github/aliok/go-defaultz/graph/badge.svg?token=VJ8020TZH7)](https://codecov.io/github/aliok/go-defaultz)
[![Go Report Card](https://goreportcard.com/badge/github.com/aliok/go-defaultz)](https://goreportcard.com/report/github.com/aliok/go-defaultz)
[![Go Reference](https://pkg.go.dev/badge/github.com/aliok/go-defaultz.svg)](https://pkg.go.dev/github.com/aliok/go-defaultz)
![License](https://img.shields.io/dub/l/vibe-d.svg)



[//]: # (TODO)
[//]: # ([![Release]&#40;https://img.shields.io/github/v/release/aliok/go-defaultz&#41;]&#40;https://github.com/aliok/go-defaultz/releases/latest&#41;)

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
- Works with pointers!
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

## Supported field types

- Primitive types: `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `string`, `bool`
- Slices of primitive types: `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`, `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`, `[]float32`, `[]float64`, `[]string`, `[]bool`
- Maps with keys or values of primitive types
- `time.Duration`, `[]time.Duration`, `map[...]time.Duration`, `map[time.Duration]...`
- All types via custom defaulters

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

[//]: # (### Defaulting custom types)

[//]: # (TODO)

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
	Field3       []string          `default:"a b"`
	Field4       map[float64]uint8 `default:"1.2:3 4.5:6"`
	Field5       time.Duration     `default:"1s"`
	Child        Child
	ChildPointer *Child
	Embedded
	Anonymous struct {
		AnonymousField string `default:"anonymousField"`
	}
}

type Child struct {
	ChildField string `default:"childField"`
}

type Embedded struct {
	EmbeddedField string `default:"embeddedField"`
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
	//      childfield: childField
	//   childpointer:
	//      childfield: childField
	//   embedded:
	//      embeddedfield: embeddedField
	//   anonymous:
	//      anonymousfield: anonymousField
}
```

[//]: # (TODO)

[//]: # (### Reading configuration)

[//]: # (TODO)

[//]: # (### More complex examples)

[//]: # (TODO)

[//]: # (## Comparing with other libraries)

[//]: # (TODO a table)

[//]: # (https://github.com/creasty/defaults)
[//]: # (https://github.com/mcuadros/go-defaults)

[//]: # (## FAQ)

[//]: # (TODO)
[//]: # (- Why no arrays?)
[//]: # (- )

[//]: # (## Development)

[//]: # (```shell)

[//]: # (#go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6)

[//]: # (#go get -u github.com/golangci/golangci-lint/cmd/golangci-lint)

[//]: # (```)

[//]: # ()
[//]: # (```shell)

[//]: # (golangci-lint run --fix)

[//]: # (goimports -w .)

[//]: # (go test -race -count=1 ./...)

[//]: # (```)
