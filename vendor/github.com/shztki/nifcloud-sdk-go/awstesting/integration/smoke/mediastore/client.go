// +build integration

//Package mediastore provides gucumber integration tests support.
package mediastore

import (
	"github.com/shztki/nifcloud-sdk-go/awstesting/integration/smoke"
	"github.com/shztki/nifcloud-sdk-go/service/mediastore"
	"github.com/gucumber/gucumber"
)

func init() {
	gucumber.Before("@mediastore", func() {
		gucumber.World["client"] = mediastore.New(smoke.Session)
	})
}
