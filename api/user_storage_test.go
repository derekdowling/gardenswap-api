package api

import (
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/gardenswap-api/gardenswap"
	"github.com/derekdowling/go-json-spec-handler/client"

	"github.com/derekdowling/go-json-spec-handler"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserHandler(t *testing.T) {

	Convey("User Handler Tests", t, func() {

		stack := BuildServer(true)
		server := httptest.NewServer(stack)
		defer server.Close()

		baseURL := server.URL

		user, err := testUser()
		So(err, ShouldBeNil)

		Convey("->parseUser()", func() {
			jsonUser, err := testJSONUser()
			So(err, ShouldBeNil)

			req, err := jsc.PostRequest(baseURL, jsonUser)
			So(err, ShouldBeNil)

			user, parseErr := parseUser(req)
			So(parseErr, ShouldBeNil)
			So(user, ShouldNotBeNil)
			So(user.ID, ShouldEqual, user.ID)
		})

		Convey("->registerUser()", func() {
			user, err = testUser()
			So(err, ShouldBeNil)

			err = registerUser(user)
			So(err, ShouldBeNil)
			So(user.ID, ShouldNotBeNil)
			So(user.JWT, ShouldNotBeNil)

			savedUser, err := gardenswap.FetchUser(user.ID)
			So(err, ShouldBeNil)
			So(savedUser, ShouldResemble, user)
		})

		Convey("POST /user", func() {

			Convey("should return a formatted output page", func() {

				jsonUser, err := testJSONUser()
				So(err, ShouldBeNil)

				doc, resp, err := jsc.Post(baseURL, jsonUser)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 201)
				So(doc, ShouldNotBeNil)
			})
		})

		Convey("GET /users", func() {

		})
	})
}

func testUser() (*gardenswap.User, error) {
	user, err := gardenswap.NewUser()
	if err != nil {
		return nil, err
	}

	user.Email = "test123"
	user.Name = "Derek"
	user.Password = "test456"
	return user, nil
}

func testJSONUser() (*jsh.Object, error) {
	testUser, err := testUser()
	if err != nil {
		return nil, err
	}

	return jsh.NewObject(testUser.ID, UserType, testUser)
}
