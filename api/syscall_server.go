package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"os"
	"sync"
)

var instance *syscallApi

// A syscallApi is an api implementation
// for a note service with access to postgres.
type syscallApi struct {
	// The embedded api features.
	*api
	// If the api is currently alive.
	alive bool
}

// Get returns the instance of a syscallApi if it exists,
// otherwise, creates one.
func Get() *syscallApi {
	if instance == nil {
		instance = &syscallApi{
			api: &api{
				mu:     sync.Mutex{},
				fiber:  fiber.New(),
				logger: zerolog.New(os.Stderr),
				port:   os.Getenv("PORT"),
				routes: nil,
			},
			alive: false,
		}
	}
	return instance
}

// IsRunning returns whether a syscallApi is alive.
func (a *syscallApi) IsRunning() bool {
	return a.alive
}

// Start is the function that starts
// both fiber and the postgres driver.
// Under the hood, Start also calls
// the syscallApi's embedded fiberInit.
func (a *syscallApi) Start() {
	if !a.alive {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.fiberInit()
	}
}
