// +build !go1.6

package request_test

import (
	"errors"

	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
)

var errTimeout = awserr.New("foo", "bar", errors.New("net/http: request canceled Timeout"))
