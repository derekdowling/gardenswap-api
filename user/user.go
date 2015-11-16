package user

import (
	"crypto/rsa"
	"io/ioutil"
	"path"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/derekdowling/gardenswap-api/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/nu7hatch/gouuid"
)

const (
	tokenExpirationHours = 72
)

const privKeyPath string = "./keys/app.rsa"
const pubKeyPath string = "./keys/app.rsa.pub"

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

	db, err := db.Get()
	if err != nil {
		log.Fatal(err.Error())
	}

	db.CreateTable(&User{})
}

// User JSONAPISpec representation
type User struct {
	ID           string `json:"id" valid:"uuidv4,required"`
	Name         string `json:"name" valid:"alpha,required"`
	Email        string `json:"email" valid:"email"`
	JWT          string `json:"jwt" valid:"alphanum,required"`
	PasswordHash string `json:"password_hash"`
	Password     string `json:",omitempty"`
}

// New creates a new object with a UUID and type
func New() (*User, error) {
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

// Save persists the most recent user model to the database
func (u *User) Save() error {

	db, err := GetDB()
	if err != nil {
		return err
	}

	if db.NewRecord(u) {
		db.Create(u)
	} else {
		db.Save(u)
	}

	return nil
}
