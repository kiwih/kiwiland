package kiwiserver

import (
	"net/http"

	"github.com/gocraft/web"
	"github.com/gorilla/sessions"
)

//UserStorer defines the interface something that can store and retrieve users must follow
type UserStorer interface {
	LoadUser(username string) (User, error)
	LoadUsernameFromSessionID(sessionID string) (string, error)
	Logout(username string) error
	AttemptLogin(propUsername string, propPassword string, rememeber bool) (string, error)
}

//Context is used in all requests as the root object of all web requests
type Context struct {
	ErrorMessages        []string
	NotificationMessages []string
	Username             string
	Data                 interface{}
	Store                *sessions.CookieStore
	Storage              UserStorer
}

//HELPER FUNCTIONS

//SetErrorMessage allows for a handler to set an error message as a "Flash" message which can be shown to the user in a later request
//(via a different handler) - it stores them in a session variable
func (c *Context) SetErrorMessage(rw web.ResponseWriter, req *web.Request, err string) {
	session, _ := c.Store.Get(req.Request, "error-messages")
	session.AddFlash(err)
	session.Save(req.Request, rw)
}

//SetNotificationMessage allows for a handler to set a notification message as a "Flash" message which can be shown to the user in a later request
//(via a different handler) - it stores them in a session variable
func (c *Context) SetNotificationMessage(rw web.ResponseWriter, req *web.Request, notification string) {
	session, _ := c.Store.Get(req.Request, "notification-messages")
	session.AddFlash(notification)
	session.Save(req.Request, rw)
}

//SetFailedRequestObject allows us to store a bad request from a form (eg not meeting the regex for the NHI parameter of system.Patient)
//so it can be recalled later for them to amend it
func (c *Context) SetFailedRequestObject(rw web.ResponseWriter, req *web.Request, requestedObject interface{}) {
	session, _ := c.Store.Get(req.Request, "error-form-requests")
	session.AddFlash(requestedObject)
	session.Save(req.Request, rw)
}

//CheckFailedRequestObject returns just one "flash" failed request object for a session. Any other request objects that were stored will be
//removed without retrieval. It allows for users to amend bad forms without needing to retype all the data
func (c *Context) CheckFailedRequestObject(rw web.ResponseWriter, req *web.Request) interface{} {
	session, _ := c.Store.Get(req.Request, "error-form-requests")
	flashes := session.Flashes()
	session.Save(req.Request, rw)

	if len(flashes) > 0 {
		//again, note that only the first one is returned. All other forms will be discarded.
		return flashes[0]
	}
	return nil
}

//MIDDLEWARE

//AssignStorageMiddleware provides a pointer to the storage object
func (c *Context) AssignStorageMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.Storage = &users
	next(rw, req)
}

//LoadUserMiddleware will load a user if possible from their session-security sessionID cookie
func (c *Context) LoadUserMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	session, _ := c.Store.Get(req.Request, "session-security")

	if session.Values["sessionID"] != nil {
		c.Username, _ = c.Storage.LoadUsernameFromSessionID(session.Values["sessionID"].(string))
	}
	next(rw, req)
}

//AssignTemplatesAndSessionsMiddleware will provide a pointer to the templates and cookie store
func (c *Context) AssignTemplatesAndSessionsMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.Store = store
	next(rw, req)
}

//GetErrorMessagesMiddleware returns any flash error messages that have been saved. Upon retrieving them, they will be deleted from the session
//(as they are "flash" session variables)
func (c *Context) GetErrorMessagesMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	session, _ := c.Store.Get(req.Request, "error-messages")
	flashes := session.Flashes()
	session.Save(req.Request, rw)

	if len(flashes) > 0 {
		//it is not possible in go to cast from []interface to []string
		strings := make([]string, len(flashes))
		for i := range flashes {
			strings[i] = flashes[i].(string)
		}
		c.ErrorMessages = strings
	}
	next(rw, req)
}

//GetNotificationMessagesMiddleware returns any flash notification messages that have been saved. Upon retrieving them, they will be deleted from the session
//(as they are "flash" session variables)
func (c *Context) GetNotificationMessagesMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	session, _ := c.Store.Get(req.Request, "notification-messages")
	flashes := session.Flashes()
	session.Save(req.Request, rw)

	if len(flashes) > 0 {
		//it is not possible in go to cast from []interface to []string
		strings := make([]string, len(flashes))
		for i := range flashes {
			strings[i] = flashes[i].(string)
		}
		c.NotificationMessages = strings
	}
	next(rw, req)
}

//RequireAccountMiddleware requires a user to be signed in
func (c *Context) RequireAccountMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if c.Username == "" {
		c.SetErrorMessage(rw, req, "You need to sign in to view this page!")
		http.Redirect(rw, req.Request, "/", http.StatusSeeOther)
	} else {
		next(rw, req)
	}
}

//HANDLERS

//GetHomeHandler returns the homepage
func (c *Context) GetHomeHandler(rw web.ResponseWriter, req *web.Request) {
	err := templates.ExecuteTemplate(rw, "indexPage", c)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

//LoginRequestForm is used when users are trying to log in
type LoginRequestForm struct {
	Username string
	Password string
	Remember bool
}

//PostSignInRequestHandler performs a sign in request (POST)
func (c *Context) PostSignInRequestHandler(rw web.ResponseWriter, req *web.Request) {
	req.ParseForm()

	var prop LoginRequestForm
	if err := decoder.Decode(&prop, req.PostForm); err != nil {
		c.SetErrorMessage(rw, req, "Decoding error: "+err.Error())
		http.Redirect(rw, req.Request, HomeURL.Make(), http.StatusSeeOther)
		return
	}

	sessionID, err := c.Storage.AttemptLogin(prop.Username, prop.Password, prop.Remember)

	if sessionID != "" && err == nil {
		//they have passed the login check. Save them to the session and redirect to management portal
		session, _ := c.Store.Get(req.Request, "session-security")
		session.Values["sessionID"] = sessionID
		//c.SetNotificationMessage(rw, req, "Hi, "+prop.Username+".") //uncomment if you want a welcome notification
		session.Save(req.Request, rw)
		http.Redirect(rw, req.Request, HomeURL.Make(), http.StatusFound)
		return
	}
	if err != nil {
		c.SetErrorMessage(rw, req, err.Error())
	} else {
		c.SetErrorMessage(rw, req, "Logging in failed (unspecified error).")
	}
	http.Redirect(rw, req.Request, HomeURL.Make(), http.StatusSeeOther)
}
