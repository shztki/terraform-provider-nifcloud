// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package computing

import (
	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/client"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/client/metadata"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/corehandlers"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/request"
	"github.com/shztki/nifcloud-sdk-go/private/protocol/computing"
	"github.com/shztki/nifcloud-sdk-go/private/signer/v2computing"
)

// Computing provides the API operation methods for making requests to
// NIFCLOUD Computing. See this package's package overview docs
// for details on the service.
//
// Computing methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type Computing struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Service information constants
const (
	ServiceName = "computing" // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the Computing client with a session.
// If additional configuration is needed for the client instance use the optional
// nifcloud.Config parameter to add your extra config.
//
// Example:
//     // Create a Computing client from just a session.
//     svc := computing.New(mySession)
//
//     // Create a Computing client with additional configuration
//     svc := computing.New(mySession, nifcloud.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*nifcloud.Config) *Computing {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg nifcloud.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) *Computing {
	svc := &Computing{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "3.0",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v2computing.SignRequestHandler)
	svc.Handlers.Sign.PushBackNamed(corehandlers.BuildContentLengthHandler)
	svc.Handlers.Build.PushBackNamed(computing.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(computing.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(computing.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(computing.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a Computing operation and runs any
// custom request initialization.
func (c *Computing) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
