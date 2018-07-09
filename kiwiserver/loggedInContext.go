package kiwiserver

import (
	"net/http"

	"github.com/gocraft/web"
)

//LoggedInContext is a helper struct for designating which handlers require a user to be logged in
type LoggedInContext struct {
	*Context
}

//SignOutRequestHandler performs the logout request
func (c *LoggedInContext) SignOutRequestHandler(rw web.ResponseWriter, req *web.Request) {
	session, _ := c.Store.Get(req.Request, "session-security")
	session.Values["sessionId"] = nil
	c.SetNotificationMessage(rw, req, "Goodbye!")

	c.Storage.Logout(c.Username)

	session.Save(req.Request, rw)
	http.Redirect(rw, req.Request, HomeURL.Make(), http.StatusFound)
}

//GetTVCommandHandler calls a command on the TV
func (c *LoggedInContext) GetTVCommandHandler(rw web.ResponseWriter, req *web.Request) {
	command, ok := req.PathParams["command"]
	if !ok {
		http.Error(rw, "400: Command not provided", http.StatusBadRequest)
		return
	}

	var cresp string
	var err error

	switch command {
	case "powerstatus":
		cresp, err = TVGetStatus()
	case "poweron":
		cresp, err = TVTurnOn()
	case "poweroff":
		cresp, err = TVTurnOff()
	case "hdmi4":
		cresp, err = TVSelectHDMI4()
	case "hdmi2":
		cresp, err = TVSelectHDMI2()
	case "hdmi1":
		cresp, err = TVSelectHDMI1()
	case "volumeup":
		cresp, err = TVVolumeUp()
	case "volumedown":
		cresp, err = TVVolumeDown()
	default:
		//unknown command
		http.Error(rw, "400: Bad tv command: "+command, http.StatusBadRequest)
		return
	}

	c.SetNotificationMessage(rw, req, cresp)
	if err != nil {
		c.SetErrorMessage(rw, req, err.Error())
	}

	http.Redirect(rw, req.Request, HomeURL.Make(), http.StatusFound)
}

//GetToshibaCommandHandler calls a command on the toshiba laptop
func (c *LoggedInContext) GetToshibaCommandHandler(rw web.ResponseWriter, req *web.Request) {
	command, ok := req.PathParams["command"]
	if !ok {
		http.Error(rw, "400: Command not provided", http.StatusBadRequest)
		return
	}

	var cresp string
	var err error

	switch command {
	case "wol":
		cresp, err = ToshibaWOL()
	default:
		//unknown command
		http.Error(rw, "400: Bad toshiba command: "+command, http.StatusBadRequest)
		return
	}

	c.SetNotificationMessage(rw, req, cresp)
	if err != nil {
		c.SetErrorMessage(rw, req, err.Error())
	}

	http.Redirect(rw, req.Request, HomeURL.Make(), http.StatusFound)
}

// func (c *Context) DoCreateFactHandler(rw web.ResponseWriter, req *web.Request) {

// 	req.ParseForm()

// 	var f fact.Fact

// 	if err := decoder.Decode(&f, req.PostForm); err != nil {
// 		c.SetErrorMessage(rw, req, "Decoding error: "+err.Error())
// 		http.Redirect(rw, req.Request, CreateFactUrl.Make(), http.StatusSeeOther)
// 		return
// 	}

// 	f.AccountId = c.Account.Id

// 	if err := fact.CreateFact(c.Storage, &f); err != nil {
// 		c.SetFailedRequestObject(rw, req, f)
// 		c.SetErrorMessage(rw, req, err.Error())
// 		http.Redirect(rw, req.Request, CreateFactUrl.Make(), http.StatusSeeOther)
// 		return
// 	}

// 	c.SetNotificationMessage(rw, req, "Fact submitted successfully!")
// 	http.Redirect(rw, req.Request, ViewFactUrl.Make("factId", strconv.FormatInt(f.Id, 10)), http.StatusFound)
// }
