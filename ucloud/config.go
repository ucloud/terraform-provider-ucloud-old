package ucloud

import (
	"os"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/ucloud/ucloud-sdk-go/ucloud/log"

	"github.com/ucloud/ucloud-sdk-go/services/uaccount"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
	"github.com/ucloud/ucloud-sdk-go/services/udisk"
	"github.com/ucloud/ucloud-sdk-go/services/udpn"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/services/ulb"
	"github.com/ucloud/ucloud-sdk-go/services/umem"
	"github.com/ucloud/ucloud-sdk-go/services/unet"
	"github.com/ucloud/ucloud-sdk-go/services/vpc"

	pumem "github.com/ucloud/ucloud-sdk-go/private/services/umem"
)

type Config struct {
	PublicKey  string
	PrivateKey string
	Region     string
	ProjectId  string

	MaxRetries int

	Insecure bool
}

type UCloudClient struct {
	region    string
	projectId string

	// public services
	uhostconn    *uhost.UHostClient
	unetconn     *unet.UNetClient
	ulbconn      *ulb.ULBClient
	vpcconn      *vpc.VPCClient
	uaccountconn *uaccount.UAccountClient
	udiskconn    *udisk.UDiskClient
	udbconn      *udb.UDBClient
	umemconn     *umem.UMemClient
	udpnconn     *udpn.UDPNClient

	// private services
	pumemconn *pumem.UMemClient
}

// Client will returns a client with connections for all product
func (c *Config) Client() (*UCloudClient, error) {
	var client UCloudClient
	client.region = c.Region
	client.projectId = c.ProjectId

	// set common attributes (region, project id, etc ...)
	config := ucloud.NewConfig()
	config.Region = c.Region
	config.ProjectId = c.ProjectId

	// enable auto retry with http/connection error
	config.MaxRetries = c.MaxRetries
	if os.Getenv("UCLOUD_DEBUG") != "false" {
		config.LogLevel = log.DebugLevel
	} else {
		config.LogLevel = log.ErrorLevel
	}

	// set endpoint with insecure https connection
	if c.Insecure {
		config.BaseUrl = GetEndpointURL(c.Region)
	} else {
		config.BaseUrl = GetInsecureEndpointURL(c.Region)
	}

	// credential with publicKey/privateKey
	credential := auth.NewCredential()
	credential.PublicKey = c.PublicKey
	credential.PrivateKey = c.PrivateKey

	// initialize client connections
	client.uhostconn = uhost.NewClient(&config, &credential)
	client.unetconn = unet.NewClient(&config, &credential)
	client.ulbconn = ulb.NewClient(&config, &credential)
	client.vpcconn = vpc.NewClient(&config, &credential)
	client.uaccountconn = uaccount.NewClient(&config, &credential)
	client.udiskconn = udisk.NewClient(&config, &credential)
	client.udbconn = udb.NewClient(&config, &credential)
	client.umemconn = umem.NewClient(&config, &credential)
	client.udpnconn = udpn.NewClient(&config, &credential)

	// initialize client connections for private usage
	client.pumemconn = pumem.NewClient(&config, &credential)

	return &client, nil
}
