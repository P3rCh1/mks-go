package v2

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/selectel/mks-go/pkg/v2/mksclient"
)

const (
	// ResourceURLCluster is the API resource path for clusters.
	ResourceURLCluster = "clusters"
	// ResourceURLKubeversion is the API resource path for Kubernetes versions.
	ResourceURLKubeversion = "kubeversions"
	// ResourceURLKubeconfig is the API resource path for kubeconfig.
	ResourceURLKubeconfig = "kubeconfig"
	// ResourceURLRotateCerts is the API resource path for certificate rotation.
	ResourceURLRotateCerts = "rotate-certs"
	// ResourceURLUpgradePatchVersion is the API resource path for patch version upgrade.
	ResourceURLUpgradePatchVersion = "upgrade-patch-version"
	// ResourceURLUpgradeMinorVersion is the API resource path for minor version upgrade.
	ResourceURLUpgradeMinorVersion = "upgrade-minor-version"
	// ResourceURLTask is the API resource path for tasks.
	ResourceURLTask = "tasks"
	// ResourceURLNodegroup is the API resource path for nodegroups.
	ResourceURLNodegroup = "nodegroups"
	// ResourceURLResize is the API resource path for resize operations.
	ResourceURLResize = "resize"
	// ResourceURLReinstall is the API resource path for reinstall operations.
	ResourceURLReinstall = "reinstall"
	// ResourceURLFeatureGates is the API resource path for feature gates.
	ResourceURLFeatureGates = "feature-gates"
	// ResourceURLAdmissionControllers is the API resource path for admission controllers.
	ResourceURLAdmissionControllers = "admission-controllers"
)

const (
	// appName represents an application name.
	appName = "mks-go/v2"

	// appVersion is a version of the application.
	appVersion = "0.1.0"

	// userAgent contains a basic user agent that will be used in queries.
	userAgent = appName + "/" + appVersion

	// defaultHTTPTimeout represents the default timeout (in seconds) for HTTP requests.
	defaultHTTPTimeout = 120

	// defaultDialTimeout represents the default timeout (in seconds) for HTTP connection establishments.
	defaultDialTimeout = 60

	// defaultKeepaliveTimeout represents the default keep-alive period for an active network connection.
	defaultKeepaliveTimeout = 60

	// defaultMaxIdleConns represents the maximum number of idle (keep-alive) connections.
	defaultMaxIdleConns = 100

	// defaultIdleConnTimeout represents the maximum amount of time an idle (keep-alive) connection will remain
	// idle before closing itself.
	defaultIdleConnTimeout = 100

	// defaultTLSHandshakeTimeout represents the default timeout (in seconds) for TLS handshake.
	defaultTLSHandshakeTimeout = 60

	// defaultExpectContinueTimeout represents the default amount of time to wait for a server's first
	// response headers.
	defaultExpectContinueTimeout = 1
)

// ServiceClient stores details that are needed to work with Selectel Managed Kubernetes Service API.
type ServiceClient struct {
	MKSClient mksclient.ClientWithResponsesInterface

	// TokenID is a client authentication token.
	TokenID string

	// Endpoint represents an endpoint that will be used in all requests.
	Endpoint string

	// UserAgent contains user agent that will be used in all requests.
	UserAgent string
}

// NewMKSClientV2 initializes a new MKS client for the V2 API.
func NewMKSClientV2(tokenID, endpoint string) *ServiceClient {
	s := &ServiceClient{
		TokenID:   tokenID,
		Endpoint:  endpoint,
		UserAgent: userAgent,
	}

	s.MKSClient = s.newMKSClient(newHTTPClient(), endpoint)

	return s
}

// NewMKSClientV2WithCustomHTTP initializes a new MKS client for the V2 API using custom HTTP client.
// If custom HTTP client is nil - default HTTP client will be used.
func NewMKSClientV2WithCustomHTTP(customHTTPClient *http.Client, tokenID, endpoint string) *ServiceClient {
	if customHTTPClient == nil {
		customHTTPClient = newHTTPClient()
	}

	s := &ServiceClient{
		TokenID:   tokenID,
		Endpoint:  endpoint,
		UserAgent: userAgent,
	}

	s.MKSClient = s.newMKSClient(customHTTPClient, endpoint)

	return s
}

// newHTTPClient returns a reference to an initialized and configured HTTP client.
func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   defaultHTTPTimeout * time.Second,
		Transport: newHTTPTransport(),
	}
}

// newMKSClient returns a reference to an initialized and configured MKS client.
func (s *ServiceClient) newMKSClient(httpClient *http.Client, endpoint string) *mksclient.ClientWithResponses {
	client, err := mksclient.NewClientWithResponses(
		endpoint,
		mksclient.WithHTTPClient(httpClient),
		mksclient.WithRequestEditorFn(s.withUserAgent),
		mksclient.WithRequestEditorFn(s.withAuthToken),
	)
	if err != nil {
		panic(err)
	}

	return client
}

// withUserAgent adds User-Agent header to request.
func (s *ServiceClient) withUserAgent(_ context.Context, req *http.Request) error {
	req.Header.Set("User-Agent", s.UserAgent)

	return nil
}

// withAuthToken adds X-Auth-Token header to request.
func (s *ServiceClient) withAuthToken(_ context.Context, req *http.Request) error {
	req.Header.Set("X-Auth-Token", s.TokenID)

	return nil
}

// newHTTPTransport returns a reference to an initialized and configured HTTP transport.
func newHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultDialTimeout * time.Second,
			KeepAlive: defaultKeepaliveTimeout * time.Second,
		}).DialContext,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       defaultIdleConnTimeout * time.Second,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout * time.Second,
		ExpectContinueTimeout: defaultExpectContinueTimeout * time.Second,
	}
}
