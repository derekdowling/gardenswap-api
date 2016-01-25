package gardenswap

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/derekdowling/gardenswap-api/api/db"
	"github.com/derekdowling/gardenswap-api/api/password"

	"github.com/dgrijalva/jwt-go"
	"github.com/nu7hatch/gouuid"
)

const (
	tokenExpirationHours = 72
)

const privKeyPath string = "../keys/app.rsa"
const pubKeyPath string = "../keys/app.rsa.pub"

var signKey *rsa.PrivateKey
var verifyKey *rsa.PublicKey

func init() {

	var signBytes []byte
	var verifyBytes []byte
	var err error

	_, filename, _, _ := runtime.Caller(0)
	signBytes, err = ioutil.ReadFile(path.Join(path.Dir(filename), privKeyPath))
	if err != nil {
		log.Print(err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Print(err)
	}
	verifyBytes, err = ioutil.ReadFile(path.Join(path.Dir(filename), pubKeyPath))
	if err != nil {
		log.Print(err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Print(err)
	}
}

// User JSONAPISpec representation
type User struct {
	ID           string `json:"id" valid:"uuidv4,required" db:"id"`
	Name         string `json:"name" valid:"alpha,required" db:"name"`
	Email        string `json:"email" valid:"email" db:"email"`
	JWT          string `json:"jwt" valid:"alphanum,required" db:"jwt"`
	PasswordHash string `json:"-" db:"password_hash"`
	Password     string `json:"password,omitempty"`
}

// NewUser creates a new object with a UUID and type
func NewUser() (*User, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &User{
		ID: uuid.String(),
	}, nil
}

// RotateJWT creates a new JWT and updates the user model with it
func (u *User) RotateJWT() error {

	// generate a new token
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims["expires"] = time.Now().Add(time.Hour * tokenExpirationHours).Unix()

	tokenStr, err := token.SignedString(signKey)
	if err != nil {
		return err
	}

	// Save new JWT
	u.JWT = tokenStr
	return nil
}

// SetPassword prepares a new encrypted password for being saved in the
// database
func (u *User) SetPassword(pass string) error {
	password, err := password.Encrypt(pass)
	if err != nil {
		return err
	}

	u.PasswordHash = password
	u.Password = ""
	return nil
}

// Save persists the most recent user model to the database
func (u *User) Save(isNew bool) error {

	// var query string

	// if isNew {
	// } else {
	// db.Save(u)
	// }

	// result, err := GetDB().Exec(query)
	// if err != nil {
	// return fmt.Errorf("Error saving user: %s", err.Error)
	// }

	return nil
}

// FetchUser attempts to find a single user based on ID
func FetchUser(id string) (*User, error) {
	user := &User{}
	err := db.Get().Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf(
			"Error attempting to fetch user by id: %s", err.Error(),
		)
	}

	return user, nil
}

// ListUsers returns a list of all users
func ListUsers() ([]*User, error) {

	users := []*User{}
	err := db.Get().Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("Error retrieving users: %s", err.Error())
	}

	return users, nil
}
