// Package unit performs initialization and validation for unit tests
package unit

import (
	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/credentials"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/session"
)

// Session is a shared session for unit tests to use.
var Session = session.Must(session.NewSession(nifcloud.NewConfig().
	WithCredentials(credentials.NewStaticCredentials("AKID", "SECRET", "SESSION")).
	WithRegion("mock-region")))
