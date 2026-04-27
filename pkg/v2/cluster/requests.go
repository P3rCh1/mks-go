package cluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	v2 "github.com/selectel/mks-go/v2"
	"github.com/selectel/mks-go/v2/mksclient"
	"gopkg.in/yaml.v3"
)

// Get returns a single cluster by its id.
func Get(ctx context.Context, client *v2.ServiceClient, clusterID string) (*mksclient.GetClusterV2Response, error) {
	responseResult, err := client.MKSClient.GetClusterV2WithResponse(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.JSON200 != nil:
		return responseResult, nil

	case responseResult.JSON404 != nil:
		return responseResult, errors.New(responseResult.JSON404.Error.Message)

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}

// List gets a list of all clusters.
func List(ctx context.Context, client *v2.ServiceClient) (*mksclient.ListClustersV2Response, error) {
	responseResult, err := client.MKSClient.ListClustersV2WithResponse(ctx)
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.JSON200 != nil:
		return responseResult, nil

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}

// Create requests a creation of a new cluster.
func Create(ctx context.Context, client *v2.ServiceClient, opts *mksclient.ClusterCreateStruct) (*mksclient.CreateClusterV2Response, error) {
	responseResult, err := client.MKSClient.CreateClusterV2WithResponse(ctx, mksclient.ClusterCreateBody{Cluster: opts})
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.JSON201 != nil:
		return responseResult, nil

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}

// Update requests an update of an existing cluster.
func Update(ctx context.Context, client *v2.ServiceClient, clusterID string, opts *mksclient.ClusterUpdateStruct) (*mksclient.UpdateClusterV2Response, error) {
	responseResult, err := client.MKSClient.UpdateClusterV2WithResponse(ctx, clusterID, mksclient.UpdateClusterV2JSONRequestBody{Cluster: opts})
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.JSON200 != nil:
		return responseResult, nil

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}

// Delete deletes a single cluster by its id.
func Delete(ctx context.Context, client *v2.ServiceClient, clusterID string) (*mksclient.DeleteClusterV2Response, error) {
	responseResult, err := client.MKSClient.DeleteClusterV2WithResponse(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.JSON204 != nil:
		return responseResult, nil

	case responseResult.JSON404 != nil:
		return responseResult, errors.New(responseResult.JSON404.Error.Message)

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}

// GetKubeconfig returns a kubeconfig by cluster id.
func GetKubeconfig(ctx context.Context, client *v2.ServiceClient, clusterID string) (*mksclient.GetClusterKubeconfigV2Response, error) {
	responseResult, err := client.MKSClient.GetClusterKubeconfigV2WithResponse(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.StatusCode() == http.StatusOK:
		return responseResult, nil

	case responseResult.JSON404 != nil:
		return responseResult, errors.New(responseResult.JSON404.Error.Message)

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}

// GetParsedKubeconfig returns a parsed kubeconfig by cluster id.
func GetParsedKubeconfig(ctx context.Context, client *v2.ServiceClient, clusterID string) (*KubeconfigInfo, error) {
	resp, err := GetKubeconfig(ctx, client, clusterID)
	if err != nil {
		return nil, err
	}

	var kubecfg Kubeconfig
	if err := yaml.Unmarshal(resp.Body, &kubecfg); err != nil {
		return nil, fmt.Errorf("invalid kubeconfig: %w", err)
	}

	info := &KubeconfigInfo{KubeconfigRaw: string(resp.Body)}

	if len(kubecfg.Clusters) == 0 {
		return nil, errors.New("invalid kubeconfig: no clusters found")
	}

	info.Server = kubecfg.Clusters[0].Cluster.Server
	info.ClusterCA = kubecfg.Clusters[0].Cluster.CertificateAuthorityData

	if len(kubecfg.Users) == 0 {
		return nil, errors.New("invalid kubeconfig: no users found")
	}

	info.ClientCert = kubecfg.Users[0].User.ClientCertificateData
	info.ClientKey = kubecfg.Users[0].User.ClientKeyData

	return info, nil
}

// RotateCerts requests a rotation of cluster certificates by cluster id.
func RotateCerts(ctx context.Context, client *v2.ServiceClient, clusterID string) (*mksclient.RotateClusterCertsV2Response, error) {
	responseResult, err := client.MKSClient.RotateClusterCertsV2WithResponse(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.JSON204 != nil:
		return responseResult, nil

	case responseResult.JSON404 != nil:
		return responseResult, errors.New(responseResult.JSON404.Error.Message)

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}

// UpgradePatchVersion requests a Kubernetes patch version upgrade by cluster id.
func UpgradePatchVersion(ctx context.Context, client *v2.ServiceClient, clusterID string) (*mksclient.UpgradePatchVersionV2Response, error) {
	responseResult, err := client.MKSClient.UpgradePatchVersionV2WithResponse(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.JSON200 != nil:
		return responseResult, nil

	case responseResult.JSON404 != nil:
		return responseResult, errors.New(responseResult.JSON404.Error.Message)

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}

// UpgradeMinorVersion requests a Kubernetes minor version upgrade by cluster id.
func UpgradeMinorVersion(ctx context.Context, client *v2.ServiceClient, clusterID string) (*mksclient.UpgradeMinorVersionV2Response, error) {
	responseResult, err := client.MKSClient.UpgradeMinorVersionV2WithResponse(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	switch {
	case responseResult.JSON200 != nil:
		return responseResult, nil

	case responseResult.JSON404 != nil:
		return responseResult, errors.New(responseResult.JSON404.Error.Message)

	case responseResult.JSON500 != nil:
		return responseResult, errors.New(responseResult.JSON500.Error.Message)
	}

	return responseResult, fmt.Errorf(v2.ErrGotHTTPStatusCodeFmt, responseResult.StatusCode())
}
