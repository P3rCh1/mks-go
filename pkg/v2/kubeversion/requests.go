package kubeversion

import (
	"context"

	v2 "github.com/selectel/mks-go/pkg/v2"
	"github.com/selectel/mks-go/pkg/v2/mksclient"
)

// List returns all supported Kubernetes versions.
func List(ctx context.Context, client *v2.ServiceClient) ([]mksclient.KubeVersionInfo, error) {
	responseResult, err := client.MKSClient.ListKubeVersionsV2WithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if responseResult.JSON200 != nil {
		if responseResult.JSON200.KubeVersions == nil {
			return nil, nil
		}

		return *responseResult.JSON200.KubeVersions, nil
	}

	return nil, mksclient.HandleAPIErrors(
		responseResult.StatusCode(), responseResult.Status(),
		responseResult.JSON500,
	)
}
