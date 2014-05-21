// package method implements http method override
// using the X-HTTP-Method-Override http header.
package method

import (
	"errors"
	"net/http"
)

// HeaderHTTPMethodOverride is a commonly used
// Http header to override the method.
const HeaderHTTPMethodOverride = "X-HTTP-Method-Override"

// ParamHTTPMethodOverride is a commonly used
// HTML form parameter to override the method.
const ParamHTTPMethodOverride = "_method"

var httpMethods = []string{"PUT", "PATCH", "DELETE"}

// ErrInvalidOverrideMethod is returned when
// an invalid http method was given to OverrideRequestMethod.
var ErrInvalidOverrideMethod = errors.New("invalid override method")

func isValidOverrideMethod(method string) bool {
	for _, m := range httpMethods {
		if m == method {
			return true
		}
	}
	return false
}

// Override checks for the X-HTTP-Method-Override header
// or the HTML for parameter, `_method`
// and uses (if valid) the http method instead of
// Request.Method.
// This is especially useful for http clients
// that don't support many http verbs.
// It isn't secure to override e.g a GET to a POST,
// so only Request.Method which are POSTs are considered.
func Override() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			m := r.FormValue(ParamHTTPMethodOverride)
			if isValidOverrideMethod(m) {
				OverrideRequestMethod(r, m)
			}
			m = r.Header.Get(HeaderHTTPMethodOverride)
			if isValidOverrideMethod(m) {
				r.Method = m
			}
		}
	})
}

// OverrideRequestMethod overrides the http
// request's method with the specified method.
func OverrideRequestMethod(r *http.Request, method string) error {
	if !isValidOverrideMethod(method) {
		return ErrInvalidOverrideMethod
	}
	r.Header.Set(HeaderHTTPMethodOverride, method)
	return nil
}
