package user

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUser(t *testing.T) {

	Convey("User Tests", t, func() {
		user, err := New()
		So(err, ShouldBeNil)
		So(user, ShouldNotBeNil)

		Convey("->SetPassword()", func() {
			err := user.SetPassword("foo")
			So(err, ShouldBeNil)
			So(user.Password, ShouldBeEmpty)
			So(user.PasswordHash, ShouldNotBeEmpty)
		})
	})
}
