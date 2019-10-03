package awstesting

import (
	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/client"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/client/metadata"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/defaults"
)

// NewClient creates and initializes a generic service client for testing.
func NewClient(cfgs ...*nifcloud.Config) *client.Client {
	info := metadata.ClientInfo{
		Endpoint:    "http://endpoint",
		SigningName: "",
	}
	def := defaults.Get()
	def.Config.MergeIn(cfgs...)

	if v := nifcloud.StringValue(def.Config.Endpoint); len(v) > 0 {
		info.Endpoint = v
	}

	return client.New(*def.Config, info, def.Handlers)
}
