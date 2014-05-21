package sessionauth

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestUser struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Age           int    `json:"age"`
	authenticated bool   `json:"-"`
}

func (u *TestUser) IsAuthenticated() bool {
	return u.authenticated
}

func (u *TestUser) Login() {
	u.authenticated = true
}

func (u *TestUser) Logout() {
	u.authenticated = false
}

func (u *TestUser) UniqueId() interface{} {
	return u.Id
}

func (u *TestUser) GetById(id interface{}) error {
	u.Id = id.(int)
	u.Name = "My Test User"
	u.Age = 42

	return nil
}

func NewUser() User {
	return &TestUser{}
}

func TestAuthenticateSession(t *testing.T) {
	store := sessions.NewCookieStore([]byte("secret123"))
	m := martini.Classic()

	m.Use(render.Renderer())
	m.Use(sessions.Sessions("my_session", store))
	m.Use(SessionUser(NewUser))

	m.Get("/setauth", func(session sessions.Session, user User) string {
		err := AuthenticateSession(session, user)
		if err != nil {
			t.Error(err)
		}
		return "OK"
	})

	m.Get("/private", LoginRequired, func(session sessions.Session, user User) string {
		return "OK"
	})

	m.Get("/logout", LoginRequired, func(session sessions.Session, user User) string {
		Logout(session, user)
		return "OK"
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/private", nil)
	m.ServeHTTP(res, req)
	if res.Code != 302 {
		t.Errorf("Private response should be 302, was %d", res.Code)
	}

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/setauth", nil)
	req1.Header.Set("Cookie", res.Header().Get("Set-Cookie"))
	m.ServeHTTP(res1, req1)
	if res1.Code != 200 {
		t.Errorf("Setauth response should be 200, was %d", res.Code)
	}

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/private", nil)
	req2.Header.Set("Cookie", res1.Header().Get("Set-Cookie"))
	m.ServeHTTP(res2, req2)
	if res2.Code != 200 {
		t.Errorf("Authenticated private response should be 200, was %d", res.Code)
	}

	res3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/logout", nil)
	req3.Header.Set("Cookie", res2.Header().Get("Set-Cookie"))
	m.ServeHTTP(res3, req3)
	if res3.Code != 302 {
		t.Errorf("Logout response should be 302, was %d", res.Code)
	}

	res4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("GET", "/private", nil)
	req4.Header.Set("Cookie", res3.Header().Get("Set-Cookie"))
	m.ServeHTTP(res4, req4)
	if res4.Code != 302 {
		t.Errorf("Unauthenticated private response should be 302, was %d", res.Code)
	}

}
