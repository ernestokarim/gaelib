package app

import (
	"encoding/json"
	"net/http"
	"strings"

	"appengine"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var (
	schemaDecoder = schema.NewDecoder()

	errHandler       Handler = nil
	forbiddenHandler Handler = nil
	notFoundHandler  Handler = nil
)

type Request struct {
	Req *http.Request
	W   http.ResponseWriter
	C   appengine.Context
}

// Load the request data using gorilla schema into a struct
func (r *Request) LoadData(data interface{}) error {
	if err := r.Req.ParseForm(); err != nil {
		return Error(err)
	}

	if err := schemaDecoder.Decode(data, r.Req.Form); err != nil {
		e, ok := err.(schema.MultiError)
		if ok {
			// Delete the invalid path errors
			for k, v := range e {
				if strings.Contains(v.Error(), "schema: invalid path") {
					delete(e, k)
				}
			}

			// Return directly if there are no other kind of errors
			if len(e) == 0 {
				return nil
			}
		}

		// Not a MultiError, log it
		if err != nil {
			return Error(err)
		}
	}

	return nil
}

func (r *Request) LoadJsonData(data interface{}) error {
	if err := json.NewDecoder(r.Req.Body).Decode(data); err != nil {
		return Error(err)
	}

	return nil
}

func (r *Request) EmitJson(data interface{}) error {
	if err := json.NewEncoder(r.W).Encode(data); err != nil {
		return Error(err)
	}

	return nil
}

func (r *Request) IsPOST() bool {
	return r.Req.Method == "POST"
}

func (r *Request) Path() string {
	u := r.Req.URL.Path
	query := r.Req.URL.RawQuery
	if len(query) > 0 {
		u += query
	}
	return u
}

// It returns a nil error always for easy of use inside the handlers.
// Example: return r.Redirect("/foo")
func (r *Request) Redirect(path string) error {
	http.Redirect(r.W, r.Req, path, http.StatusFound)
	return nil
}

// It returns a nil error always for easy of use inside the handlers.
// Example: return r.RedirectPermanently("/foo")
func (r *Request) RedirectPermanently(path string) error {
	http.Redirect(r.W, r.Req, path, http.StatusMovedPermanently)
	return nil
}

func (r *Request) ExecuteTemplate(names []string, data interface{}) error {
	return RawExecuteTemplate(r.W, names, data)
}

func (r *Request) JsonResponse(data interface{}) error {
	if err := json.NewEncoder(r.W).Encode(data); err != nil {
		return Error(err)
	}
	return nil
}

func (r *Request) processError(err error) {
	e, ok := (err).(*AppError)
	if !ok {
		e = Error(err).(*AppError)
	}

	e.Log(r.C)
	if e.Code == 500 && errHandler != nil {
		if err := errHandler(r); err == nil {
			return
		}
	} else if e.Code == 404 && notFoundHandler != nil {
		if err := notFoundHandler(r); err == nil {
			return
		}
	} else if e.Code == 403 && forbiddenHandler != nil {
		if err := forbiddenHandler(r); err == nil {
			return
		}
	}

	http.Error(r.W, "", e.Code)
}

func SetErrorHandler(f Handler) {
	errHandler = f
}

func SetForbiddenHandler(f Handler) {
	forbiddenHandler = f
}

func SetNotFoundHandler(r *mux.Router, f Handler) {
	notFoundHandler = f
	r.NotFoundHandler = f
}
