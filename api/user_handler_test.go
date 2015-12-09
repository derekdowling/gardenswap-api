package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/derekdowling/gardenswap-api/api/user"

	"github.com/derekdowling/go-json-spec-handler"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserHandler(t *testing.T) {

	Convey("User Handler Tests", t, func() {

		stack := BuildServer()
		server := httptest.NewServer(stack)
		defer server.Close()

		baseURL := server.URL

		testUser, err := buildTestUser()
		So(err, ShouldBeNil)

		Convey("->parseUser()", func() {
			jsonUser, err := testUserJSON(testUser)
			So(err, ShouldBeNil)
			req := testRequest(jsonUser)

			user, parseErr := parseUser(req)
			So(parseErr, ShouldBeNil)
			So(user, ShouldNotBeNil)
			So(user.ID, ShouldEqual, testUser.ID)
		})

		Convey("->createUser()", func() {
			user, err := buildTestUser()
			So(err, ShouldBeNil)

			err = createUser(user)
			So(err, ShouldBeNil)
			So(user.ID, ShouldNotBeNil)
			So(user.JWT, ShouldNotBeNil)

			savedUser, err := GetUser(user.ID)
			So(err, ShouldBeNil)
			So(savedUser, ShouldResemble, user)
		})

		Convey("POST /user", func() {

			Convey("should return a formatted output page", func() {

				user, err := buildTestUser()
				So(err, ShouldBeNil)

				req, err := getUserRequest("GET", baseURL, user)
				So(err, ShouldBeNil)

				resp, err := testRequest(req)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 201)

				// Decode profile and deep compare structs
				respUser, parseErr := parseJSONUser(resp.Body)
				So(parseErr, ShouldBeNil)
				So(respUser.ID, ShouldEqual, testUser.ID)
				So(respUser.Password, ShouldBeEmpty)
				So(respUser.PasswordHash, ShouldBeEmpty)
				So(respUser.JWT, ShouldNotBeEmpty)
			})
		})

		Convey("GET /users", func() {

		})
	})
}

func testRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Origin", "http://localhost")
	client := &http.Client{}
	return client.Do(req)
}

func getUserRequest(method string, baseURL string, user *user.User) (*http.Request, error) {
	obj, err := jsh.NewObject(user.ID, "user", user)
	if err != nil {
		return nil, err
	}

	url := &url.URL{Host: baseURL}

	req, reqErr := jsh.NewObjectRequest(method, url, obj)
	if reqErr != nil {
		return nil, reqErr
	}

	return req, nil
}

func buildTestUser() (*user.User, error) {
	user, err := user.New()
	if err != nil {
		return nil, err
	}

	user.Email = "test123"
	user.Name = "Derek"
	user.Password = "test456"
	return user, nil
}
