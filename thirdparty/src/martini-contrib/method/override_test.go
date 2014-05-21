package method

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var tests = []struct {
	Method         string
	OverrideMethod string
	ExpectedMethod string
}{
	{"POST", "PUT", "PUT"},
	{"POST", "PATCH", "PATCH"},
	{"POST", "DELETE", "DELETE"},
	{"GET", "GET", "GET"},
	{"HEAD", "HEAD", "HEAD"},
	{"GET", "PUT", "GET"},
	{"HEAD", "DELETE", "HEAD"},
}

func TestOverride(t *testing.T) {
	for _, test := range tests {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(test.Method, "/", nil)
		if err != nil {
			t.Error(err)
		}
		OverrideRequestMethod(r, test.OverrideMethod)
		Override().ServeHTTP(w, r)
		if r.Method != test.ExpectedMethod {
			t.Errorf("Expected %s, got %s", test.ExpectedMethod, r.Method)
		}
	}
}

func selectRoute(r martini.Router, method string, h martini.Handler) {
	switch method {
	case "GET":
		r.Get("/", h)
	case "PATCH":
		r.Patch("/", h)
	case "POST":
		r.Post("/", h)
	case "PUT":
		r.Put("/", h)
	case "DELETE":
		r.Delete("/", h)
	case "OPTIONS":
		r.Options("/", h)
	case "HEAD":
		r.Head("/", h)
	default:
		panic("bad method")
	}
}

func TestMartiniSelectiveRouter(t *testing.T) {
	for _, test := range tests {
		w := httptest.NewRecorder()
		r := martini.NewRouter()

		done := make(chan bool)
		selectRoute(r, test.ExpectedMethod, func(rq *http.Request) {
			done <- true
		})

		req, err := http.NewRequest(test.Method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		OverrideRequestMethod(req, test.OverrideMethod)

		m := martini.New()
		m.Use(Override())
		m.Action(r.Handle)
		go m.ServeHTTP(w, req)
		select {
		case <-done:
		case <-time.After(30 * time.Millisecond):
			t.Errorf("Expected router to route to %s, got something else (%v).", test.ExpectedMethod, test)
		}
	}
}

func TestInMartini(t *testing.T) {
	for _, test := range tests {
		w := httptest.NewRecorder()
		m := martini.New()
		m.Use(Override())
		m.Use(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != test.ExpectedMethod {
				t.Errorf("Expected %s, got %s", test.ExpectedMethod, r.Method)
			}
		})

		r, err := http.NewRequest(test.Method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		OverrideRequestMethod(r, test.OverrideMethod)

		m.ServeHTTP(w, r)
	}

}

func TestParamenterOverrideInMartini(t *testing.T) {
	for _, test := range tests {
		w := httptest.NewRecorder()
		m := martini.New()
		m.Use(Override())
		m.Use(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != test.ExpectedMethod {
				t.Errorf("Expected %s, got %s", test.ExpectedMethod, r.Method)
			}
		})

		query := "_method=" + test.OverrideMethod
		r, err := http.NewRequest(test.Method, "/?"+query, nil)
		if err != nil {
			t.Fatal(err)
		}

		m.ServeHTTP(w, r)
	}

}
