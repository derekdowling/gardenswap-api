package api

import (
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {

	Convey("Server Tests", t, func() {
		Convey("->BuildRouter()", func() {
			api := &API{}
			router := buildRouter(api)
			So(router, ShouldNotBeNil)

			server := httptest.NewServer(router)
			defer server.Close()
		})

		Convey("->BuildServer()", func() {
			server := BuildServer()
			So(server, ShouldNotBeNil)
		})
	})
}
