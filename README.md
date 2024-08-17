#  📦 [![GoDoc](https://pkg.go.dev/badge/automatica.team/di)](https://pkg.go.dev/automatica.team/di)

`automatica.team/di`  is a dead simple dependency injection tool for Go. Supports declarative configuration.

## Installation

> `$ go get -u automatica.team/di`

## Usage

Declarative configuration is optional, but is very handy and recommended to use:

```yaml
# example/di.yml

# A version of the current dependency state
version: 0

# Dependencies are described here
di:
  # Order does matter, keep individual dependencies on the top
  x/db:
    # Immediate dependency configuration
    prepare_stmt: true
    # A '$' sign points to look up the env
    path: $DB_PATH # "example.db" is also valid
  # `x/hit` depends on `x/db`, that's why it's defined as the second
  x/hit: {}

# A list of Go imports for external dependency integration
imports:
  # Generation tool will add this import path, thus it will
  # be able to inject the external dependency if it satisfies
  # the `di.Dependency` interface
  - github.com/external/dependency

```

Use a `cmd/di` generation tool to automatically inject the dependencies:

```go
// example/main.gen.go

package main

import (
	"automatica.team/di"
	"automatica.team/di/example/db"
	"automatica.team/di/example/hit"
)

// Injections generated automatically via `go:generate`.
func init() {
	di.Inject(db.New())
	di.Inject(hit.New())
}
```

The first defined dependency that opens the database connection:

```go
// example/db/db.go

package db

import (
	"fmt"

	"automatica.team/di"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ = (di.D)(&DB{})

// A shortcut for `di` generation to avoid searching for a proper type.
func New() *DB { return &DB{} }

// DB implements `di.D` that is `Dependency`.
type DB struct {
	*gorm.DB
}

// Name is a unique identifier of your dependency.
func (DB) Name() string {
	// Use the "x/" prefix for the local dependencies.
	return "x/db"
}

// New defines a constructor for dependency
// The `di.C` stands for `Config`.
func (DB) New(c di.C) (di.D, error) {
	// Using built-in configuration helper, `err` will be
	// properly formatted to display the missing required param.
	path, err := c.EnvString("path")
	if err != nil {
		return nil, err
	}

	conf := &gorm.Config{
		// Use `Must` to get zero value for optional params.
		PrepareStmt: di.Must(c.Bool("prepare_stmt")),
	}

	db, err := gorm.Open(sqlite.Open(path), conf)
	if err != nil {
		return nil, fmt.Errorf("x/db: %w", err)
	}

	return &DB{DB: db}, nil
}
```

Dependency can inject other dependencies declared above it. In this example, `x/hit` depends on the `x/db`.
`Hitter` provides some server logic to be used later in the final `Server` runnable:

```go
// example/hit/hit.go

package hit

import (
	"time"

	"automatica.team/di"
	"automatica.team/di/example/db"
)

var _ = (di.D)(&Hitter{})

func New() *Hitter {
	return &Hitter{}
}

type Hitter struct {
	db *db.DB `di:"x/db"`
}

func (Hitter) Name() string {
	return "x/hit"
}

func (h Hitter) New(c di.C) (di.D, error) {
	// Automatically migrate the `hits` table.
	// `h.db` is already accessible.
	return New(), h.db.AutoMigrate(&Hit{})
}

// Hit represents a database model for server hits.
type Hit struct {
	When time.Time
	IP   string
}

// Specify a table name for GORM.
func (Hit) TableName() string {
	return "hits"
}

// Hit adds a new IP hit to the database.
func (h Hitter) Hit(ip string) error {
	return h.db.Create(&Hit{
		When: time.Now(),
		IP:   ip,
	}).Error
}
```


You end up by calling a single final runnable in the `main` function:

```go
// example/main.go

package main

import (
	"net/http"

	"automatica.team/di"
	"automatica.team/di/example/hit"

	"github.com/labstack/echo/v4"
)

//go:generate $(go env GOBIN)/di

func main() {
	// Parse is optional, but it's handy for configuring
	// the dependencies in a single config file.
	if err := di.Parse(); err != nil {
		panic(err)
	}
	// Running a final handler that in some way uses
	// the injected dependencies.
	if err := di.Run[Server](); err != nil {
		panic(err)
	}
}

// Server is a final runnable entity.
type Server struct {
	// Inject any declared dependency by its name.
	hit *hit.Hitter `di:"x/hit"`
}

// Run implements `di.R` that is `Runnable`.
func (s Server) Run() error {
	e := echo.New()
	e.HideBanner = true
	e.GET("/", s.onHit)
	return e.Start(":8080")
}

func (s Server) onHit(c echo.Context) error {
	// Using injected `x/hit` dependency.
	if err := s.hit.Hit(c.RealIP()); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
```
