package api

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/derekdowling/gardenswap-api/api/password"
	"github.com/derekdowling/gardenswap-api/db"
	"github.com/derekdowling/gardenswap-api/user"
	"github.com/derekdowling/go-json-spec-handler"
)

// Register handles persisting a new user and their relevant identifiers
func (a *API) Register(w http.ResponseWriter, r *http.Request) {

	userObj, err := jsh.ParseObject(r)

	// TODO: Check email isn't already taken
	user, err := parseJSONUser(r.Body)
	if err != nil {
		msg := fmt.Sprintf("Unable to parse user data: %s.", err.Error())
		a.Logger.Error(msg)
		a.HandleError(w, &Error{Detail: msg, Status: 422})
		return
	}

	// Create a new user
	err = createUser(user)
	if err != nil {
		a.Logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Unable to create new user.")
		a.HandleError(w, &Error{Detail: "Unable to create account.", Status: 500})
		return
	}

	a.SendResponse(w, 201, userToData(user))
}

// GetUsers returns a list of all users
func (a *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []*user.User{}

	db, err := db.GetDB()
	if err != nil {
		a.Logger.Error(err.Error())
		a.HandleError(w, ISE(err))
		return
	}

	db.Find(&users)
	a.SendResponse(w, 200, usersToData(users))
}

// CreateUser converts an incoming JSON body into a User struct
func createUser(user *user.User) error {

	if user.Name == "" ||
		user.Password == "" ||
		user.Email == "" {
		return fmt.Errorf("Missing one of 'name', 'email', 'password' which are mandatory fields.")
	}

	password, err := password.Encrypt(user.Password)
	if err != nil {
		return err
	}

	user.PasswordHash = password
	user.Password = ""

	err = user.RotateJWT()
	if err != nil {
		return err
	}

	return user.Save()
}

func userToData(user *user.User) *Data {
	return &Data{
		Type:       "user",
		ID:         user.ID,
		Attributes: user,
	}
}

func usersToData(users []*user.User) []*Data {
	data := []*Data{}

	for _, user := range users {
		data = append(data, userToData(user))
	}

	return data
}

// userFromJSON returns a decoded User from JSON
func parseJSONUser(r io.ReadCloser) (*user.User, error) {
	jObject, err := parseJSONObject(r)
	if err != nil {
		return nil, fmt.Errorf("Error parsing user: %s", err.Error())
	}

	if jObject.Data == nil {
		return nil, fmt.Errorf("No data for request: %+v", jObject)
	}

	return jObject.Data.Attributes, nil
}
