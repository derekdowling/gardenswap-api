package api

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {

	Convey("Server Tests", t, func() {
		Convey("->BuildServer()", func() {
			server := BuildServer(true)
			So(server, ShouldNotBeNil)
		})
	})
}
