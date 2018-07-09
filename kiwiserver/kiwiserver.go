package kiwiserver

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

var (
	users Userlist
	store *sessions.CookieStore

	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("./media/templates/*")) //this initializes the template engine
	decoder   = schema.NewDecoder()                                                             //this initializes the schema (HTML form decoding) engine
)

//StartServer will start a kiwiserver listening at the given address and with the provided cookie store salt
func StartServer(serverAddress string, cookieStoreSalt string) {
	// pass, err := bcrypt.GenerateFromPassword([]byte("testing1+"), 10)
	// if err != nil {
	// 	panic("something's wrong with bcrypt")
	// }
	//log.Printf("Password: '%s'", pass)

	LoadUsers()

	decoder.RegisterConverter(false, ConvertBool)

	store = sessions.NewCookieStore([]byte(cookieStoreSalt))

	router := initRouter()

	log.Println("Server running at " + serverAddress)
	if err := http.ListenAndServe(serverAddress, router); err != nil {
		log.Println("Error:", err.Error())
	}
}

//ConvertBool is used to convert checkboxes to golang boolean values
func ConvertBool(value string) reflect.Value {
	if value == "on" {
		return reflect.ValueOf(true)
	} else if v, err := strconv.ParseBool(value); err == nil {
		return reflect.ValueOf(v)
	}

	return reflect.ValueOf(false)
}
