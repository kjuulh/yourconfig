# 🌱 yourconfig

A dead-simple Go library for loading config structs from environment variables — with optional tagging and required enforcement.

## 📦 Install

```bash
go get github.com/kjuulh/yourconfig
```

## 🧪 Basic Usage

```go
package main

import (
	"fmt"
	"github.com/kjuulh/yourconfig"
)

type Config struct {
	// Uses SNAKE_CASE: SOME_ITEM
	SomeItem string `conf:"required:true"`

	// Custom name
	OtherItem string `conf:"MY_ITEM"`

	// Loads YET_ANOTHER optionally from env
  YetAnother string `conf:""`

	// ignored field
  Ignored string
}

func main() {
	cfg := yourconfig.MustLoad[Config]()
	fmt.Println(cfg.SomeItem, cfg.OtherItem, cfg.YetAnother)
}
```

## 🏷️ Tag Syntax

You configure fields using the `conf` struct tag:

* **Custom env name:** First item (e.g. `conf:"MY_ENV"`)
* **Options:** `key:value` or flags (e.g. `required:true`, or just `required`)

Examples:

```go
// infers env var SOME_ITEM
SomeItem string `conf:""`

// uses DIFFERENT_NAME
SomeItem string `conf:"DIFFERENT_NAME"`

// requires the env var to be present
SomeItem string `conf:"required"`

// equivalent to above
SomeItem string `conf:"required:true"`

// combine name and required
SomeItem string `conf:"DIFFERENT_NAME,required"`

// combine name and required (option)
SomeItem string `conf:"DIFFERENT_NAME,required:true"`

```

## ❗ Required Fields

If a field is marked `required` and the environment variable is **not set**, loading will return an error. Use `MustLoad[T]()` to panic on error, or `Load[T]() (T, error)` to handle it gracefully.

Private (unexported) fields are ignored unless tagged — in which case they produce a `not settable` error, avoid setting `conf` on private items.

## License

MIT

## Contributing

```go
go test ./...  
```

## Release

Releases are handled through cuddle-please a pr based release tool, simply merge the pr and a release will be cut
