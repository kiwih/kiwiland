package kiwiserver

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//A User is a user of this website
type User struct {
	Username       string
	Password       string
	SessionID      string
	SessionExpires time.Time
}

//A Userlist is a slice of users, useful for the website to keep around
type Userlist struct {
	Users []User
}

//LoadUser will load a user from the userlist
func (ul *Userlist) LoadUser(username string) (User, error) {
	for _, u := range ul.Users {
		if u.Username == username {
			return u, nil
		}
	}
	return User{}, errors.New("Username not found")
}

//LoadUsernameFromSessionID will, given a valid sessionID, attempt to load a user
func (ul *Userlist) LoadUsernameFromSessionID(sessionID string) (string, error) {
	for _, u := range ul.Users {
		if u.SessionID == sessionID {
			//make sure session is still valid
			if u.SessionExpires.After(time.Now()) {
				return u.Username, nil
			}
		}
	}
	return "", errors.New("Username not found")
}

//Logout will, given a username, expire that user's session
func (ul *Userlist) Logout(username string) error {
	for i := 0; i < len(ul.Users); i++ {
		if ul.Users[i].Username == username {
			ul.Users[i].SessionID = ""
			ul.Users[i].SessionExpires = time.Now()
			SaveUsers()
			return nil
		}
	}
	return errors.New("Username not found")
}

//SaveUser updates all fields of a user except username
func (ul *Userlist) SaveUser(u User) error {
	for i := 0; i < len(ul.Users); i++ {
		if ul.Users[i].Username == u.Username {
			ul.Users[i] = u
			SaveUsers()
			return nil
		}
	}
	return errors.New("Username not found")
}

//GenerateValidationKey will generate a random md5 hash
func GenerateValidationKey() (string, error) {
	//generate validation key
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(b)), nil
}

//AttemptLogin will return an error if it fails to log in, or a sessionId string if it succeeds
func (ul *Userlist) AttemptLogin(propUsername string, propPassword string, remember bool) (string, error) {
	propUser, err := ul.LoadUser(propUsername)

	if err == nil {
		if err = bcrypt.CompareHashAndPassword([]byte(propUser.Password), []byte(propPassword)); err == nil {
			//successful login.
			//generate session
			propUser.SessionID, err = GenerateValidationKey()
			if err != nil {
				log.Println("Bad validation key")
				return "", err
			}
			if !remember {
				propUser.SessionExpires = time.Now().Add(3600 * time.Second) //expires in 1 hour if they don't want to be remembered
			} else {
				propUser.SessionExpires = time.Now().AddDate(0, 1, 0) //expires in 1 month if they want to be remembered
			}
			log.Printf("Login successful for user '%s'", propUsername)
			return propUser.SessionID, ul.SaveUser(propUser)
		}
		log.Println("Bad password")
		return "", errors.New("Invalid Username or Password")
	}
	log.Println("Bad username")
	return "", errors.New("Invalid Username or Password")
}

const (
	userFile = "users.json"
)

//LoadUsers will load the users from a user-provided json file
//or make a new user if none exist
func LoadUsers() {
	userBytes, err := ioutil.ReadFile(userFile)
	if err != nil {
		log.Printf("No user file provided. What will your username be?")
		username := ""
		password := ""
		fmt.Scanln(&username)
		log.Printf("What will your password be?")
		fmt.Scanln(&password)
		if username == "" || password == "" {
			log.Fatalf("Bad input for new username/password and no user file existing.")
		}
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			panic("something's wrong with bcrypt")
		}
		users = Userlist{Users: []User{User{Username: username, Password: string(hashedPass)}}}
		userBytes, _ = json.MarshalIndent(users, "", "\t")
		ioutil.WriteFile(userFile, userBytes, 0644)
		return
	}
	err = json.Unmarshal(userBytes, &users)
	if err != nil {
		log.Fatalf("User file is broken. Delete it and restart program.")
	}
}

//SaveUsers will save the user file
func SaveUsers() {
	userBytes, _ := json.MarshalIndent(users, "", "\t")
	ioutil.WriteFile(userFile, userBytes, 0644)
	log.Println("Saved user file")
	return
}
