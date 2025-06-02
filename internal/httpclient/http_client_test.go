package httpclient_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/billz-2/ofd_connector/internal/httpclient"
)

var (
	ctx context.Context
)

func TestMain(m *testing.M) {
	mySetupFunction()
	retCode := m.Run()
	myTeardownFunction()
	os.Exit(retCode)
}

func mySetupFunction() {
	println("start httpclient package testing")

	ctx = context.Background()

	err := os.Setenv("ENV_FILE_PATH", "../../../configs/.env.testing")
	if err != nil {
		panic(err)
	}

}

func myTeardownFunction() {
	println("success end httpclient package testing")
}

func TestHTTPClient(t *testing.T) {
	Convey("HTTPClient", t, func() {
		client := httpclient.NewHTTPClient(20)
		So(client, ShouldImplement, (*httpclient.HTTPClient)(nil))

		req, err := httpclient.NewHTTPRequest("https://jsonplaceholder.typicode.com/todos/1",
			http.MethodGet, "application/json", nil, nil)
		So(err, ShouldBeNil)

		resp := client.Request(ctx, req)
		So(resp.Error, ShouldBeNil)
		So(resp.Body, ShouldNotBeEmpty)
		So(resp.StatusCode, ShouldEqual, http.StatusOK)
	})
}
