package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
	"os"
	"sync"
)

// A Method is an integer specifier
// for a http method.
type Method int

const (
	// GET is the enumerator for the http get method.
	GET Method = iota
	// POST is the enumerator for the http post method.
	POST
	// DELETE is the enumerator for the http delete method.
	DELETE
)

var (
	// The methods map holds all the string
	// associations for a Method.
	methods = map[Method]string{
		GET:    "GET",
		POST:   "POST",
		DELETE: "DELETE",
	}
)

// A RouteDefinition is a definition that may be used
// in shorthand when adding fiber handlers.
type RouteDefinition struct {
	// The fiber handles.
	handlers []fiber.Handler
	// The usable methods.
	methods []Method
	// The path of the route.
	path string
}

// An api is a simple structure
// for a fiber application that
// may serve function.
type api struct {
	// The synchronization mutex, mu.
	mu sync.Mutex
	// A pointer to the fiber application.
	fiber *fiber.App
	// The logger.
	logger zerolog.Logger
	// The application's port.
	port string
	// The routes that the api serves.
	routes []RouteDefinition
}

// DefineRoute is a simple function used for defining fiber.Handler s with respective
// guidelines as RouteDefinitions, to later be used when registering the route with fiber.
func DefineRoute(path string, methods []Method, handlers ...fiber.Handler) *RouteDefinition {
	return &RouteDefinition{
		handlers: handlers,
		methods:  methods,
		path:     path,
	}
}

// String gets a Method's string form from the method-string
// association map.
func (m Method) String() string {
	return methods[m]
}

// fiberInit is a function across all api implementations that
// will initialize fiber with a cors configuration as a goroutine.
func (a *api) fiberInit() {
	go func() {
		origin := os.Getenv("FRONTEND_URL")
		// todo: failure reading data but
		// 	the requests still work.
		a.fiber.Use(cors.New(cors.Config{
			AllowCredentials: true,
			AllowOrigins:     origin,
			AllowHeaders:     "Origin",
			ExposeHeaders:    "Origin",
		}))
		err := a.fiber.Listen(":" + a.port)
		a.handleError(err)
	}()
}

// Logger returns an api's logger.
func (a *api) Logger() *zerolog.Logger {
	return &a.logger
}

// Fiber returns an api's fiber instance.
func (a *api) Fiber() *fiber.App {
	return a.fiber
}

// Port returns an api's port.
func (a *api) Port() string {
	return a.port
}

// AddRoutes adds the specified routes to an api.
func (a *api) AddRoutes(routes ...*RouteDefinition) {
	for _, route := range routes {
		a.AddRoute(route)
	}
}

// AddRoute adds a singular route to an api.
func (a *api) AddRoute(route *RouteDefinition) {
	for _, method := range route.methods {
		a.fiber.Add(method.String(), route.path, route.handlers...)
	}
}

// Routes returns an api's routes.
func (a *api) Routes() []RouteDefinition {
	return a.routes
}

// handleError is an error handling function built-in
// to an api.
func (a *api) handleError(err error) {
	if err != nil {
		// If the error exists, log it.
		a.logger.Warn().Err(err)
	}
}
