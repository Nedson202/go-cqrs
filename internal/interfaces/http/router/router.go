package router

import (
	"net/http"

	"github.com/nedson202/go-cqrs/internal/interfaces/http/handlers"
)

// Middleware defines a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Handlers contains all HTTP handlers for the application
type Handlers struct {
    User *handlers.UserHandler
}

// Router handles HTTP routing with middleware support
type Router struct {
    handlers   *Handlers
    middleware []Middleware
    mux        *http.ServeMux
}

// New creates a new Router instance
func New(h *Handlers, mw ...Middleware) *Router {
    return &Router{
        handlers:   h,
        middleware: mw,
        mux:        http.NewServeMux(),
    }
}

// Setup configures all routes and middleware
func (r *Router) Setup() http.Handler {
    // Register routes
    r.mux.HandleFunc("POST /users", r.handlers.User.HandleCreateUser)
    r.mux.HandleFunc("PUT /users/{id}", r.handlers.User.HandleUpdateUser)
    r.mux.HandleFunc("GET /users/{id}", r.handlers.User.HandleGetUser)
    r.mux.HandleFunc("GET /users", r.handlers.User.HandleListUsers)

    // Apply middleware in reverse order
    var handler http.Handler = r.mux
    for i := len(r.middleware) - 1; i >= 0; i-- {
        handler = r.middleware[i](handler)
    }

    return handler
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    r.Setup().ServeHTTP(w, req)
} 
