package kiwiserver

import (
	"html/template"
)

var funcMap = template.FuncMap{
	"GetSignInURL":         SignInURL.Make,
	"GetSignOutURL":        SignOutURL.Make,
	"GetHomeURL":           HomeURL.Make,
	"GetTVCommandURL":      GetTVCommandURL,
	"GetToshibaCommandURL": GetToshibaCommandURL,
} //this provides templates with the ability to run useful functions

//GetTVCommandURL makes a tv command URL
func GetTVCommandURL(command string) string {
	return TVCommandURL.Make("command", command)
}

//GetToshibaCommandURL makes a tv command URL
func GetToshibaCommandURL(command string) string {
	return ToshibaCommandURL.Make("command", command)
}
