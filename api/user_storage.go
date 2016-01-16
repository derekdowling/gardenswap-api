package api

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/derekdowling/gardenswap-api/api/db"
	"github.com/derekdowling/gardenswap-api/api/user"
	"github.com/derekdowling/go-stdlogger"

	sq "github.com/Masterminds/squirrel"
	"github.com/derekdowling/go-json-spec-handler"
)

// UserAPI implements the jsh-api CRUD interface
type UserAPI struct {
	Logger std.Logger
}

// List returns a list of all users
func (u *UserAPI) List(ctx context.Context) (jsh.List, jsh.ErrorType) {
	users, err := user.All()
	if err != nil {
		return nil, jsh.ISE(err.Error())
	}

	list := jsh.List{}

	for _, user := range users {
		obj, err := jsh.NewObject(user.ID, "users", user)
		if err != nil {
			return nil, jsh.ISE(fmt.Sprintf("Error converting user to response object: %s", err.Error()))
		}
		list = append(list, obj)
	}

	return list, nil
}

// Save persistes a User object
func (u *UserAPI) Save(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	return object, nil
}

// Get retrieves a user by ID
func (u *UserAPI) Get(ctx context.Context, id string) (*jsh.Object, jsh.ErrorType) {
	user := &user.User{}

	object, err := jsh.NewObject(id, "user", user)
	if err != nil {
		return nil, err
	}

	return object, nil
}

// Update modifies an existing user
func (u *UserAPI) Update(ctx context.Context, object *jsh.Object) (*jsh.Object, jsh.ErrorType) {
	return object, nil
}

// Delete removes a user by ID
func (u *UserAPI) Delete(ctx context.Context, id string) jsh.ErrorType {
	return nil
}

// Register handles persisting a new user and their relevant identifiers
func (u *UserAPI) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	user, err := parseUser(r)
	if err != nil {
		u.Logger.Printf("Error parsing user request data: %s", err.Error())
		jsh.Send(w, r, err)
		return
	}

	// Create a new user
	registrationErr := registerUser(user)
	if registrationErr != nil {
		u.Logger.Printf("Unable to create new user: %s", registrationErr.Error())
		jsh.Send(w, r, jsh.ISE("Unable to create account."))
		return
	}

	obj, err := jsh.NewObject(user.ID, "user", user)
	if err != nil {
		u.Logger.Printf("Error creating user response: %s", err.Error())
		jsh.Send(w, r, err)
		return
	}

	jsh.Send(w, r, obj)
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

// parseUser marshals a user from a JSONAPI request into our internal User
// type
func parseUser(r *http.Request) (*user.User, jsh.ErrorType) {
	userObj, err := jsh.ParseObject(r)
	if err != nil {
		return nil, err
	}

	user := &user.User{ID: userObj.ID}
	unmarshalErr := userObj.Unmarshal("user", user)
	if err != nil {
		return nil, unmarshalErr
	}

	return user, nil
}
