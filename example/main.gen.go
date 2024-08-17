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
