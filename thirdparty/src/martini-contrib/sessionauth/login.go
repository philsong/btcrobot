// Package login is a middleware for Martini that provides a simple way to track user sessions
// in on a website. Please see https://github.com/codegangsta/martini-contrib/blob/master/sessionauth/README.md
// for a more detailed description of the package.
package sessionauth

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini-contrib/sessions"
	"log"
	"net/http"
)

// These are the default configuration values for this package. They
// can be set at anytime, probably during the initial setup of Martini.
var (
	// RedirectUrl should be the relative URL for your login route
	RedirectUrl string = "/login"

	// RedirectParam is the query string parameter that will be set
	// with the page the user was trying to visit before they were
	// intercepted.
	RedirectParam string = "next"

	// SessionKey is the key containing the unique ID in your session
	SessionKey string = "AUTHUNIQUEID"
)

// User defines all the functions necessary to work with the user's authentication.
// The caller should implement these functions for whatever system of authentication
// they choose to use
type User interface {
	// Return whether this user is logged in or not
	IsAuthenticated() bool

	// Set any flags or extra data that should be available
	Login()

	// Clear any sensitive data out of the user
	Logout()

	// Return the unique identifier of this user object
	UniqueId() interface{}

	// Populate this user object with values
	GetById(id interface{}) error
}

// SessionUser will try to read a unique user ID out of the session. Then it tries
// to populate an anonymous user object from the database based on that ID. If this
// is successful, the valid user is mapped into the context. Otherwise the anonymous
// user is mapped into the contact.
// The newUser() function should provide a valid 0value structure for the caller's
// user type.
func SessionUser(newUser func() User) martini.Handler {
	return func(s sessions.Session, c martini.Context, l *log.Logger) {
		userId := s.Get(SessionKey)
		user := newUser()

		if userId != nil {
			err := user.GetById(userId)
			if err != nil {
				l.Printf("Login Error: %v\n", err)
			} else {
				user.Login()
			}
		}

		c.MapTo(user, (*User)(nil))
	}
}

// AuthenticateSession will mark the session and user object as authenticated. Then
// the Login() user function will be called. This function should be called after
// you have validated a user.
func AuthenticateSession(s sessions.Session, user User) error {
	user.Login()
	return UpdateUser(s, user)
}

// Logout will clear out the session and call the Logout() user function.
func Logout(s sessions.Session, user User) {
	user.Logout()
	s.Delete(SessionKey)
}

// LoginRequired verifies that the current user is authenticated. Any routes that
// require a login should have this handler placed in the flow. If the user is not
// authenticated, they will be redirected to /login with the "next" get parameter
// set to the attempted URL.
func LoginRequired(r render.Render, user User, req *http.Request) {
	if user.IsAuthenticated() == false {
		path := fmt.Sprintf("%s?%s=%s", RedirectUrl, RedirectParam, req.URL.Path)
		r.Redirect(path, 302)
	}
}

// UpdateUser updates the User object stored in the session. This is useful incase a change
// is made to the user model that needs to persist across requests.
func UpdateUser(s sessions.Session, user User) error {
	s.Set(SessionKey, user.UniqueId())
	return nil
}
