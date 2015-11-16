package password

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
)

func TestSpec(t *testing.T) {

	Convey("Encrypt Testing", t, func() {

		Convey("generateSalt()", func() {
			salt, err := generateSalt()
			So(err, ShouldBeNil)
			So(salt, ShouldNotBeBlank)
			So(len(salt), ShouldEqual, SaltLength)
		})

		Convey("hashPassword()", func() {
			salt, err := generateSalt()
			So(err, ShouldBeNil)

			hash, err := hashPassword(salt, "hershmahgersh")
			So(err, ShouldBeNil)
			So(hash, ShouldNotBeNil)

			cost, err := bcrypt.Cost([]byte(hash))
			So(err, ShouldBeNil)
			So(cost, ShouldEqual, EncryptCost)

			pw := combine(salt, string(hash))
			log.Printf("pw = %+v\n", pw)
			parsedSalt, parsedHash := getPWPieces(pw)
			So(salt, ShouldEqual, parsedSalt)
			So(string(hash), ShouldEqual, parsedHash)
		})

		Convey("->Encrypt()", func() {
			passString := "mmmPassword1"
			password, err := Encrypt(passString)
			So(err, ShouldBeNil)
			So(password, ShouldNotBeEmpty)

			salt, _ := getPWPieces(password)
			So(len(salt), ShouldEqual, SaltLength)
		})

		Convey("->IsMatch()", func() {
			password := "megaman49"
			hash, err := Encrypt(password)
			So(err, ShouldBeNil)

			match, err := IsMatch(password, hash)
			So(match, ShouldBeTrue)
			So(err, ShouldBeNil)

			match, err = IsMatch("lolfail", hash)
			So(match, ShouldBeFalse)
			So(err, ShouldNotBeNil)

			match, err = IsMatch("Megaman49", hash)
			So(match, ShouldBeFalse)
			So(err, ShouldNotBeNil)
		})
	})
}
