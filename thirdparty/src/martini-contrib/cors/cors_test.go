package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/codegangsta/martini"
)

func Test_AllowAll(t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.New()
	m.Use(Allow(&Options{
		AllowAllOrigins: true,
	}))

	r, _ := http.NewRequest("PUT", "foo", nil)
	m.ServeHTTP(recorder, r)

	if recorder.HeaderMap.Get(headerAllowOrigin) != "*" {
		t.Errorf("Allow-Origin header should be *")
	}
}

func Test_AllowRegexMatch(t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.New()
	m.Use(Allow(&Options{
		AllowOrigins: []string{"https://aaa.com", "https://foo\\.*"},
	}))

	origin := "https://foo.com"
	r, _ := http.NewRequest("PUT", "foo", nil)
	r.Header.Add("Origin", origin)
	m.ServeHTTP(recorder, r)

	headerValue := recorder.HeaderMap.Get(headerAllowOrigin)
	if headerValue != origin {
		t.Errorf("Allow-Origin header should be %v, found %v", origin, headerValue)
	}
}

func Test_AllowRegexNoMatch(t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.New()
	m.Use(Allow(&Options{
		AllowOrigins: []string{"https://foo\\.*"},
	}))

	origin := "https://bar.com"
	r, _ := http.NewRequest("PUT", "foo", nil)
	r.Header.Add("Origin", origin)
	m.ServeHTTP(recorder, r)

	headerValue := recorder.HeaderMap.Get(headerAllowOrigin)
	if headerValue != "" {
		t.Errorf("Allow-Origin header should not exist, found %v", headerValue)
	}
}

func Test_OtherHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.New()
	m.Use(Allow(&Options{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"PATCH", "GET"},
		AllowHeaders:     []string{"Origin", "X-whatever"},
		ExposeHeaders:    []string{"Content-Length", "Hello"},
		MaxAge:           5 * time.Minute,
	}))

	r, _ := http.NewRequest("PUT", "foo", nil)
	m.ServeHTTP(recorder, r)

	credentialsVal := recorder.HeaderMap.Get(headerAllowCredentials)
	methodsVal := recorder.HeaderMap.Get(headerAllowMethods)
	headersVal := recorder.HeaderMap.Get(headerAllowHeaders)
	exposedHeadersVal := recorder.HeaderMap.Get(headerExposeHeaders)
	maxAgeVal := recorder.HeaderMap.Get(headerMaxAge)

	if credentialsVal != "true" {
		t.Errorf("Allow-Credentials is expected to be true, found %v", credentialsVal)
	}

	if methodsVal != "PATCH,GET" {
		t.Errorf("Allow-Methods is expected to be PATCH,GET; found %v", methodsVal)
	}

	if headersVal != "Origin,X-whatever" {
		t.Errorf("Allow-Headers is expected to be Origin,X-whatever; found %v", headersVal)
	}

	if exposedHeadersVal != "Content-Length,Hello" {
		t.Errorf("Expose-Headers are expected to be Content-Length,Hello. Found %v", exposedHeadersVal)
	}

	if maxAgeVal != "300" {
		t.Errorf("Max-Age is expected to be 300, found %v", maxAgeVal)
	}
}

func Test_Preflight(t *testing.T) {
	recorder := httptest.NewRecorder()
	m := martini.New()
	m.Use(Allow(&Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH"},
		AllowHeaders:    []string{"Origin", "X-whatever"},
	}))

	r, _ := http.NewRequest("OPTIONS", "foo", nil)
	r.Header.Add(headerRequestMethod, "PUT")
	r.Header.Add(headerRequestHeaders, "X-whatever")
	m.ServeHTTP(recorder, r)

	methodsVal := recorder.HeaderMap.Get(headerAllowMethods)
	headersVal := recorder.HeaderMap.Get(headerAllowHeaders)

	if methodsVal != "PUT,PATCH" {
		t.Errorf("Allow-Methods is expected to be PUT,PATCH, found %v", methodsVal)
	}

	if headersVal != "X-whatever" {
		t.Errorf("Allow-Headers is expected to be X-whatever, found %v", headersVal)
	}
}

func Benchmark_WithoutCORS(b *testing.B) {
	recorder := httptest.NewRecorder()
	m := martini.New()

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		r, _ := http.NewRequest("PUT", "foo", nil)
		m.ServeHTTP(recorder, r)
	}
}

func Benchmark_WithCORS(b *testing.B) {
	recorder := httptest.NewRecorder()
	m := martini.New()
	m.Use(Allow(&Options{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"PATCH", "GET"},
		AllowHeaders:     []string{"Origin", "X-whatever"},
		MaxAge:           5 * time.Minute,
	}))

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		r, _ := http.NewRequest("PUT", "foo", nil)
		m.ServeHTTP(recorder, r)
	}
}
