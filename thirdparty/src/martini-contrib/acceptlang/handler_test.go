package acceptlang

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type acceptLanguageTest struct {
	path     string
	header   string
	expected AcceptLanguages
}

var acceptLanguageTests = []acceptLanguageTest{
	// Test an empty header
	{"/none", "", make(AcceptLanguages, 0)},

	// Test a single unqualified header value
	{"/single", "en-gb", AcceptLanguages{AcceptLanguage{"en-gb", 1}}},

	// Test a single qualified header value
	{"/single_qualified", "en-gb;q=0.8", AcceptLanguages{AcceptLanguage{"en-gb", 0.8}}},

	// Test multiple unqualified header values
	{"/multiple", "en-gb, nl,en-us", AcceptLanguages{
		AcceptLanguage{"en-gb", 1}, AcceptLanguage{"nl", 1}, AcceptLanguage{"en-us", 1},
	}},

	// Test multiple qualified header values
	{"/multiple_qualified", "en-gb;q=0.2, nl;q=1,en-us;q=0.5", AcceptLanguages{
		AcceptLanguage{"nl", 1}, AcceptLanguage{"en-us", 0.5}, AcceptLanguage{"en-gb", 0.2},
	}},
}

func TestAcceptLanguageTests(t *testing.T) {
	for _, test := range acceptLanguageTests {
		m := martini.Classic()
		m.Get(test.path, Languages(), func(result AcceptLanguages) {
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Unexpected test result:\nExpected: %#v\nResult: %#v", test.expected, result)
			}
		})

		recorder := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", test.path, nil)
		if test.header != "" {
			r.Header.Add(acceptLanguageHeader, test.header)
		}
		m.ServeHTTP(recorder, r)
	}
}

func BenchmarkLanguages1(b *testing.B) {
	m := newBenchmarkMartini()

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/benchmark", nil)
	r.Header.Add(acceptLanguageHeader, "en-us;q=0.7")

	for n := 0; n < b.N; n++ {
		m.ServeHTTP(recorder, r)
	}
}

func BenchmarkLanguages6(b *testing.B) {
	m := newBenchmarkMartini()

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/benchmark", nil)
	r.Header.Add(acceptLanguageHeader, "en-us;q=0.7, en-GB;q=0.8, de;q=1, nl;q=0.1, fr-FR;q=0.3, es")

	for n := 0; n < b.N; n++ {
		m.ServeHTTP(recorder, r)
	}
}

func newBenchmarkMartini() *martini.ClassicMartini {
	router := martini.NewRouter()
	base := martini.New()
	base.Action(router.Handle)

	m := &martini.ClassicMartini{base, router}
	m.Get("/benchmark", Languages(), func(result AcceptLanguages) {
		//b.Logf("Parsed languages: %s", result)
	})

	return m
}
