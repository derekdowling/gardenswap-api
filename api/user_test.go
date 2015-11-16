package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/derekdowling/gardenswap-api/db"
	"github.com/derekdowling/gardenswap-api/user"
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

		Convey("->parseJSONUser()", func() {
			jsonUser, err := testUserJSON(testUser)
			So(err, ShouldBeNil)
			req := testRequest(jsonUser)

			user, parseErr := parseJSONUser(req.Body)
			So(parseErr, ShouldBeNil)
			So(user, ShouldNotBeNil)
		})

		Convey("->createUser()", func() {
			user, err := buildTestUser()
			So(err, ShouldBeNil)

			err = createUser(user)
			So(err, ShouldBeNil)
			So(user.ID, ShouldNotBeNil)
			So(user.JWT, ShouldNotBeNil)

			db, err := db.GetDB()
			So(db.NewRecord(user), ShouldBeFalse)
		})

		Convey("POST /user", func() {

			Convey("should return a formatted output page", func() {
				usersURL := combineURL(baseURL, "/users")
				So(err, ShouldBeNil)

				jsonStr, err := testUserJSON(testUser)
				So(err, ShouldBeNil)

				resp, err := testResponse(usersURL, "POST", jsonStr)
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
	})
}

func testResponse(url string, method string, json []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", jSONType)
	req.Header.Set("Origin", "http://localhost")

	client := &http.Client{}
	return client.Do(req)
}

func testUserJSON(user *user.User) ([]byte, error) {
	return json.Marshal(&JSONObject{
		Data: userToData(user),
	})
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
