package app

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"appengine"
)

// All handlers in the app must implement this type
type Handler func(r *Request) error

// Serves a http request
func (fn Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Creates the AppEngine context from the request.
	c := appengine.NewContext(req)

	// Ask for chrome frame if we're in IE.
	w.Header().Set("X-UA-Compatible", "chrome=1")

	// Create the request.
	r := &Request{Req: req, W: w, C: c}

	// Defers the panic recovering.
	defer func() {
		if rec := recover(); rec != nil {
			err := fmt.Errorf("panic recovered error: %+v\n%s", rec, debug.Stack())
			r.LogError(err)
			r.internalServerError(err.Error())
		}
	}()

	// Call the handler.
	if err := fn(r); err != nil {
		r.LogError(err)
		r.internalServerError(err.Error())
	}
}
