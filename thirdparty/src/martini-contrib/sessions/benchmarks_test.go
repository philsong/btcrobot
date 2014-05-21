package sessions

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkNoSessionsMiddleware(b *testing.B) {
	m := testMartini()
	m.Get("/foo", func() string {
		return "Foo"
	})

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo", nil)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.ServeHTTP(recorder, r)
	}
}

func BenchmarkSessionsNoWrites(b *testing.B) {
	m := testMartini()
	store := NewCookieStore([]byte("secret123"))
	m.Use(Sessions("my_session", store))
	m.Get("/foo", func() string {
		return "Foo"
	})

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo", nil)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.ServeHTTP(recorder, r)
	}
}

func BenchmarkSessionsWithWrite(b *testing.B) {
	m := testMartini()
	store := NewCookieStore([]byte("secret123"))
	m.Use(Sessions("my_session", store))
	m.Get("/foo", func(s Session) string {
		s.Set("foo", "bar")
		return "Foo"
	})

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo", nil)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.ServeHTTP(recorder, r)
	}
}

func BenchmarkSessionsWithRead(b *testing.B) {
	m := testMartini()
	store := NewCookieStore([]byte("secret123"))
	m.Use(Sessions("my_session", store))
	m.Get("/foo", func(s Session) string {
		s.Get("foo")
		return "Foo"
	})

	recorder := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo", nil)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.ServeHTTP(recorder, r)
	}
}

func testMartini() *martini.ClassicMartini {
	m := martini.Classic()
	m.Handlers()
	return m
}
