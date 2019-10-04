package rest_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/client"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/client/metadata"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/request"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/signer/v4"
	"github.com/shztki/nifcloud-sdk-go/awstesting/unit"
	"github.com/shztki/nifcloud-sdk-go/private/protocol/rest"
)

func TestUnsetHeaders(t *testing.T) {
	cfg := &nifcloud.Config{Region: nifcloud.String("us-west-2")}
	c := unit.Session.ClientConfig("testService", cfg)
	svc := client.New(
		*cfg,
		metadata.ClientInfo{
			ServiceName:   "testService",
			SigningName:   c.SigningName,
			SigningRegion: c.SigningRegion,
			Endpoint:      c.Endpoint,
			APIVersion:    "",
		},
		c.Handlers,
	)

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(rest.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(rest.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(rest.UnmarshalMetaHandler)
	op := &request.Operation{
		Name:     "test-operation",
		HTTPPath: "/",
	}

	input := &struct {
		Foo nifcloud.JSONValue `location:"header" locationName:"x-amz-foo" type:"jsonvalue"`
		Bar nifcloud.JSONValue `location:"header" locationName:"x-amz-bar" type:"jsonvalue"`
	}{}

	output := &struct {
		Foo nifcloud.JSONValue `location:"header" locationName:"x-amz-foo" type:"jsonvalue"`
		Bar nifcloud.JSONValue `location:"header" locationName:"x-amz-bar" type:"jsonvalue"`
	}{}

	req := svc.NewRequest(op, input, output)
	req.HTTPResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(nil)), Header: http.Header{}}
	req.HTTPResponse.Header.Set("X-Amz-Foo", "e30=")

	// unmarshal response
	rest.UnmarshalMeta(req)
	rest.Unmarshal(req)
	if req.Error != nil {
		t.Fatal(req.Error)
	}
}
