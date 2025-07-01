package httpclient_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
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

func TestHttpReqUrlEncoded(t *testing.T) {
	Convey("TestHttpReqUrlEncoded", t, func() {
		bodyValues := map[string]any{
			"a": 1,
			"b": 2,
		}
		bodyBytes, err := json.Marshal(bodyValues)
		So(err, ShouldBeNil)
		req, err := httpclient.NewHTTPRequest("https://jsonplaceholder.typicode.com/todos/1",
			http.MethodGet, "application/x-www-form-urlencoded", bodyBytes, nil)
		So(err, ShouldBeNil)
		So(req.ContentType, ShouldEqual, "application/x-www-form-urlencoded")
		So(req.Body, ShouldNotBeEmpty)
		Convey("TestHttpReqUrlEncoded", func() {
			urlValues, err := url.ParseQuery(string(req.Body))
			So(err, ShouldBeNil)
			So(urlValues.Get("a"), ShouldEqual, "1")
			So(urlValues.Get("b"), ShouldEqual, "2")
		})
	})

}
