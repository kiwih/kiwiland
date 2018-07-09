package kiwiserver

import (
	"strings"

	"github.com/gocraft/web"
)

//URL is a helper type for URL string types
type URL string

//XxxxUrl are the URLs used in this web application
const (
	HomeURL           URL = "/"
	SignInURL         URL = "/signin"
	SignOutURL        URL = "/signout"
	TVCommandURL      URL = "/tv/:command"
	ToshibaCommandURL URL = "/toshiba/:command"
)

//String() converts a URL to a string
func (u URL) String() string {
	return string(u)
}

//Make is a helper function to helpfully set params up for urls that use regex
func (u URL) Make(param ...string) string {
	if len(param)%2 != 0 {
		panic("Make URL " + u.String() + " had non-even number of params")
	}

	retStr := u.String()

	for i := 0; i < len(param); i += 2 {
		retStr = strings.Replace(retStr, ":"+param[i], param[i+1], 1)
	}
	return retStr
}

func initRouter() *web.Router {

	rootRouter := web.New(Context{})
	rootRouter.Middleware(web.LoggerMiddleware)
	rootRouter.Middleware(web.ShowErrorsMiddleware)
	rootRouter.Middleware(web.StaticMiddleware("./media/public", web.StaticOption{Prefix: "/public"})) // "public" is a directory to serve files from.)
	rootRouter.Middleware((*Context).AssignStorageMiddleware)
	rootRouter.Middleware((*Context).AssignTemplatesAndSessionsMiddleware)
	rootRouter.Middleware((*Context).LoadUserMiddleware)
	rootRouter.Middleware((*Context).GetErrorMessagesMiddleware)
	rootRouter.Middleware((*Context).GetNotificationMessagesMiddleware)

	//rootRouter web paths
	rootRouter.Get(HomeURL.String(), (*Context).GetHomeHandler)

	//sign in
	rootRouter.Post(SignInURL.String(), (*Context).PostSignInRequestHandler)

	//must be logged in for some handlers...
	loggedInRouter := rootRouter.Subrouter(LoggedInContext{}, "/")
	loggedInRouter.Middleware((*LoggedInContext).RequireAccountMiddleware)

	//sign out handler
	loggedInRouter.Get(SignOutURL.String(), (*LoggedInContext).SignOutRequestHandler)

	//handlers
	loggedInRouter.Get(TVCommandURL.String(), (*LoggedInContext).GetTVCommandHandler)
	loggedInRouter.Get(ToshibaCommandURL.String(), (*LoggedInContext).GetToshibaCommandHandler)

	//create, delete fact handlers

	return rootRouter
}
