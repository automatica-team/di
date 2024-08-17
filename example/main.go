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
