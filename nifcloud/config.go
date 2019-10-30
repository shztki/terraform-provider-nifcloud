package nifcloud

import (
	"fmt"
	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/credentials"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/session"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
	"github.com/shztki/nifcloud-sdk-go/service/rdb"
)

// Config is struct
type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Endpoint  string
}

// NifcloudClient is struct
type NifcloudClient struct {
	computingconn *computing.Computing
	rdbconn       *rdb.Rdb
}

// Client is function
func (c *Config) Client() (interface{}, error) {
	if c.Region == "" {
		return nil, fmt.Errorf("[Err] No Region Name for Nifcloud")
	}

	var credential *credentials.Credentials
	if c.AccessKey != "" && c.SecretKey != "" {
		credential = credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	} else {
		credential = credentials.NewEnvCredentials()
	}

	config := nifcloud.Config{
		Region:      nifcloud.String(c.Region),
		Credentials: credential,
	}

	if c.Endpoint != "" {
		config.WithEndpoint(*nifcloud.String(c.Endpoint))
	}

	sess := session.Must(session.NewSession(&config))

	var client NifcloudClient

	client.computingconn = computing.New(sess)
	client.rdbconn = rdb.New(sess)

	return &client, nil
}
