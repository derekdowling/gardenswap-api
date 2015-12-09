package api

import (
	"fmt"
	"net/http"

	"github.com/derekdowling/gardenswap-api/api/db"
	"github.com/derekdowling/gardenswap-api/api/user"

	sq "github.com/Masterminds/squirrel"
	"github.com/derekdowling/go-json-spec-handler"
)

// Register handles persisting a new user and their relevant identifiers
func (a *API) Register(w http.ResponseWriter, r *http.Request) {

	user, err := parseUser(r)
	if err != nil {
		a.Logger.Errorf("Error parsing user request data: %s", err.Error())
		jsh.Send(w, r, err)
		return
	}
	// Create a new user
	registrationErr := registerUser(user)
	if registrationErr != nil {
		a.Logger.Errorf("Unable to create new user: %s", registrationErr.Error())
		jsh.Send(w, r, jsh.ISE("Unable to create account."))
		return
	}

	obj, err := jsh.NewObject(user.ID, "user", user)
	if err != nil {
		a.Logger.Errorf("Error creating user response: %s", err.Error())
		jsh.Send(w, r, err)
		return
	}

	jsh.Send(w, r, obj)
}

// ListUsers returns a list of all users
func (a *API) ListUsers(w http.ResponseWriter, r *http.Request) {

	users, err := user.All()
	if err != nil {
		a.Logger.Errorf("Error getting all users: %s", err.Error())
		jsh.Send(w, r, jsh.ISE("Unable to get users."))
		return
	}

	list := &jsh.List{}

	for _, user := range users {
		obj, err := jsh.NewObject(user.ID, "user", user)
		if err != nil {
			a.Logger.Errorf("Error converting user to response object: %s", err.Error())
			jsh.Send(w, r, err)
		}
		list.Add(obj)
	}

	jsh.Send(w, r, list)
}

func parseUser(r *http.Request) (*user.User, jsh.SendableError) {

	userObj, err := jsh.ParseObject(r)
	if err != nil {
		return nil, err
	}

	user := &user.User{ID: userObj.ID}
	err = userObj.Unmarshal("user", user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser converts an incoming JSON body into a User struct
func registerUser(usr *user.User) error {

	err := usr.SetPassword(usr.Password)
	if err != nil {
		return err
	}

	err = usr.RotateJWT()
	if err != nil {
		return err
	}

	query, _, err := sq.Insert("users").
		Columns("id", "name", "email", "jwt", "password").
		Values(usr.ID, usr.Name, usr.Email, usr.JWT, usr.PasswordHash).
		ToSql()

	if err != nil {
		return fmt.Errorf("Error creating SQL: %s", err.Error())
	}

	_, err = db.Get().Exec(query)
	if err != nil {
		return fmt.Errorf("Error inserting new user: %s", err.Error())
	}

	return nil
}
