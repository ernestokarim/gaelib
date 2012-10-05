package app

import (
	"fmt"
	"net/http"

	"appengine"
)

// All handlers in the app must implement this type
type Handler func(r *Request) error

// Serves a http request
func (fn Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)

	w.Header().Set("X-UA-Compatible", "chrome=1")

	r := &Request{Req: req, W: w, C: c}

	defer func() {
		if rec := recover(); rec != nil {
			err := Error(fmt.Errorf("panic recovered error: %s", rec))
			r.processError(err)
		}
	}()

	if err := fn(r); err != nil {
		r.processError(err)
	}
}
