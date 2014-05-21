package secure

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_No_Config(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
	// nothing here to configure
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Body.String(), `bar`)
}

func Test_No_AllowHosts(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		AllowedHosts: []string{},
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Body.String(), `bar`)
}

func Test_Good_Single_AllowHosts(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		AllowedHosts: []string{"www.example.com"},
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Body.String(), `bar`)
}

func Test_Bad_Single_AllowHosts(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		AllowedHosts: []string{"sub.example.com"},
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusInternalServerError)
}

func Test_Good_Multiple_AllowHosts(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		AllowedHosts: []string{"www.example.com", "sub.example.com"},
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "sub.example.com"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Body.String(), `bar`)
}

func Test_Bad_Multiple_AllowHosts(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		AllowedHosts: []string{"www.example.com", "sub.example.com"},
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www3.example.com"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusInternalServerError)
}

func Test_SSL(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		SSLRedirect: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "https"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func Test_SSL_In_Dev_Mode(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Dev
	m.Use(Secure(Options{
		SSLRedirect: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func Test_SSL_In_Dev_Mode_But_Disable_Prod_Check(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Dev
	m.Use(Secure(Options{
		SSLRedirect:      true,
		DisableProdCheck: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://www.example.com/foo")
}

func Test_Basic_SSL(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		SSLRedirect: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://www.example.com/foo")
}

func Test_Basic_SSL_With_Host(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		SSLRedirect: true,
		SSLHost:     "secure.example.com",
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://secure.example.com/foo")
}

func Test_Bad_Proxy_SSL(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		SSLRedirect: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://www.example.com/foo")
}

func Test_Custom_Proxy_SSL(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func Test_Custom_Proxy_SSL_In_Dev_Mode(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Dev
	m.Use(Secure(Options{
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "http")

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func Test_Custom_Proxy_And_Host_SSL(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		SSLHost:         "secure.example.com",
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
}

func Test_Custom_Bad_Proxy_And_Host_SSL(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		SSLRedirect:     true,
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "superman"},
		SSLHost:         "secure.example.com",
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.Host = "www.example.com"
	req.URL.Scheme = "http"
	req.Header.Add("X-Forwarded-Proto", "https")

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusMovedPermanently)
	expect(t, res.Header().Get("Location"), "https://secure.example.com/foo")
}

func Test_STS_Header(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		STSSeconds: 315360000,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("Strict-Transport-Security"), "max-age=315360000")
}

func Test_STS_Header_In_Dev_Mode(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Dev
	m.Use(Secure(Options{
		STSSeconds: 315360000,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("Strict-Transport-Security"), "")
}

func Test_STS_Header_With_Subdomain(t *testing.T) {
	m := martini.Classic()
	martini.Env = martini.Prod
	m.Use(Secure(Options{
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("Strict-Transport-Security"), "max-age=315360000; includeSubdomains")
}

func Test_Frame_Deny(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		FrameDeny: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("X-Frame-Options"), "DENY")
}

func Test_Custom_Frame_Value(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		CustomFrameOptionsValue: "SAMEORIGIN",
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("X-Frame-Options"), "SAMEORIGIN")
}

func Test_Custom_Frame_Value_With_Deny(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("X-Frame-Options"), "SAMEORIGIN")
}

func Test_Content_Nosniff(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		ContentTypeNosniff: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("X-Content-Type-Options"), "nosniff")
}

func Test_XSS_Protection(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		BrowserXssFilter: true,
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("X-XSS-Protection"), "1; mode=block")
}

func Test_CSP(t *testing.T) {
	m := martini.Classic()
	m.Use(Secure(Options{
		ContentSecurityPolicy: "default-src 'self'",
	}))

	m.Get("/foo", func() string {
		return "bar"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, http.StatusOK)
	expect(t, res.Header().Get("Content-Security-Policy"), "default-src 'self'")
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
