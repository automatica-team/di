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
