package cluster

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	v2 "github.com/selectel/mks-go/pkg/v2"
	"github.com/selectel/mks-go/pkg/v2/mksclient"
	mksmock "github.com/selectel/mks-go/pkg/v2/mksclient/mocks"
	"github.com/selectel/mks-go/pkg/v2/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.GetClusterV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.GetClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.ClusterResp{
					Cluster: &mksclient.ClusterDetailed{
						Id:                            clusterID,
						Name:                          "test-cluster",
						KubeVersion:                   "1.28.0",
						Basic:                         false,
						EnableAutorepair:              true,
						EnablePatchVersionAutoUpgrade: true,
						KubeApiIp:                     "10.0.0.1",
						CreatedAt:                     time.Now(),
						CniType:                       mksclient.ClusterDetailedCniType("cilium"),
						NetworkType:                   mksclient.ClusterDetailedNetworkType("default"),
						Status:                        mksclient.ClusterDetailedStatus("active"),
						AdditionalSoftware:            map[string]any{"nginx-ingress": "enabled"},
					},
				},
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.GetClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "internal server error",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.GetClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericNotFoundError{
					Error: struct {
						Id      string `json:"id"` //nolint:revive // it's generated struct
						Message string `json:"message"`
					}{
						Id:      clusterID,
						Message: "cluster not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "cluster not found",
			},
		},
		{
			name: "unknown status",
			clientResponse: &mksclient.GetClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Status:     http.StatusText(http.StatusServiceUnavailable),
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    http.StatusText(http.StatusServiceUnavailable),
			},
		},
		{
			name:        "http error",
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().GetClusterV2WithResponse(mock.Anything, clusterID).Return(test.clientResponse, test.clientError)

			cluster, err := Get(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID)

			if test.errExpected != nil {
				assert.Nil(t, cluster)
				require.Error(t, err)

				var mksErrExp *mksclient.MKSError
				if !errors.As(test.errExpected, &mksErrExp) {
					assert.ErrorIs(t, err, test.errExpected)

					return
				}

				var mksErr *mksclient.MKSError
				require.ErrorAs(t, err, &mksErr)
				assert.Equal(t, mksErrExp, mksErr)

				return
			}

			require.NoError(t, err)

			assert.Equal(t, test.clientResponse.JSON200.Cluster, cluster)
		})
	}
}

func TestCreate(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.CreateClusterV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.CreateClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusCreated,
					Status:     http.StatusText(http.StatusCreated),
				},
				JSON201: &mksclient.ClusterResp{
					Cluster: &mksclient.ClusterDetailed{
						Id:                            clusterID,
						Name:                          "test-cluster",
						KubeVersion:                   "1.28.0",
						Basic:                         false,
						EnableAutorepair:              true,
						EnablePatchVersionAutoUpgrade: true,
						KubeApiIp:                     "10.0.0.1",
						CreatedAt:                     time.Now(),
						CniType:                       mksclient.ClusterDetailedCniType("cilium"),
						NetworkType:                   mksclient.ClusterDetailedNetworkType("default"),
						Status:                        mksclient.ClusterDetailedStatus("active"),
						AdditionalSoftware:            map[string]any{"nginx-ingress": "enabled"},
					},
				},
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.CreateClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "internal server error",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "unknown status",
			clientResponse: &mksclient.CreateClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Status:     http.StatusText(http.StatusServiceUnavailable),
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    http.StatusText(http.StatusServiceUnavailable),
			},
		},
		{
			name:        "http error",
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().CreateClusterV2WithResponse(mock.Anything, mock.Anything).Return(test.clientResponse, test.clientError)

			cluster, err := Create(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, nil)

			if test.errExpected != nil {
				assert.Nil(t, cluster)
				require.Error(t, err)

				var mksErrExp *mksclient.MKSError
				if !errors.As(test.errExpected, &mksErrExp) {
					assert.ErrorIs(t, err, test.errExpected)

					return
				}

				var mksErr *mksclient.MKSError
				require.ErrorAs(t, err, &mksErr)
				assert.Equal(t, mksErrExp, mksErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.clientResponse.JSON201.Cluster, cluster)
		})
	}
}

func TestUpdate(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.UpdateClusterV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.UpdateClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.ClusterResp{
					Cluster: &mksclient.ClusterDetailed{
						Id:                            clusterID,
						Name:                          "test-cluster",
						KubeVersion:                   "1.28.0",
						Basic:                         false,
						EnableAutorepair:              true,
						EnablePatchVersionAutoUpgrade: true,
						KubeApiIp:                     "10.0.0.1",
						CreatedAt:                     time.Now(),
						CniType:                       mksclient.ClusterDetailedCniType("cilium"),
						NetworkType:                   mksclient.ClusterDetailedNetworkType("default"),
						Status:                        mksclient.ClusterDetailedStatus("active"),
						AdditionalSoftware:            map[string]any{"nginx-ingress": "enabled"},
					},
				},
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.UpdateClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "internal server error",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "unknown status",
			clientResponse: &mksclient.UpdateClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Status:     http.StatusText(http.StatusServiceUnavailable),
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    http.StatusText(http.StatusServiceUnavailable),
			},
		},
		{
			name:        "http error",
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().UpdateClusterV2WithResponse(mock.Anything, clusterID, mock.Anything).Return(test.clientResponse, test.clientError)

			cluster, err := Update(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nil)

			if test.errExpected != nil {
				assert.Nil(t, cluster)
				require.Error(t, err)

				var mksErrExp *mksclient.MKSError
				if !errors.As(test.errExpected, &mksErrExp) {
					assert.ErrorIs(t, err, test.errExpected)

					return
				}

				var mksErr *mksclient.MKSError
				require.ErrorAs(t, err, &mksErr)
				assert.Equal(t, mksErrExp, mksErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.clientResponse.JSON200.Cluster, cluster)
		})
	}
}

func TestDelete(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.DeleteClusterV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.DeleteClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNoContent,
					Status:     http.StatusText(http.StatusNoContent),
				},
				JSON204: testutils.Ptr(any(struct{}{})),
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.DeleteClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericNotFoundError{
					Error: struct {
						Id      string `json:"id"` //nolint:revive // it's generated struct
						Message string `json:"message"`
					}{
						Id:      clusterID,
						Message: "cluster not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "cluster not found",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.DeleteClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "internal server error",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "unknown status",
			clientResponse: &mksclient.DeleteClusterV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Status:     http.StatusText(http.StatusServiceUnavailable),
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    http.StatusText(http.StatusServiceUnavailable),
			},
		},
		{
			name:        "http error",
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().DeleteClusterV2WithResponse(mock.Anything, clusterID).Return(test.clientResponse, test.clientError)

			err := Delete(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID)

			if test.errExpected != nil {
				require.Error(t, err)

				var mksErrExp *mksclient.MKSError
				if !errors.As(test.errExpected, &mksErrExp) {
					assert.ErrorIs(t, err, test.errExpected)

					return
				}

				var mksErr *mksclient.MKSError
				require.ErrorAs(t, err, &mksErr)
				assert.Equal(t, mksErrExp, mksErr)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetKubeconfig(t *testing.T) {
	const clusterID = "test-cluster-id"

	testKubeconfig := []byte("kubeconfig")

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.GetClusterKubeconfigV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.GetClusterKubeconfigV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				Body: testKubeconfig,
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.GetClusterKubeconfigV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericNotFoundError{
					Error: struct {
						Id      string `json:"id"` //nolint:revive // it's generated struct
						Message string `json:"message"`
					}{
						Id:      clusterID,
						Message: "cluster not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "cluster not found",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.GetClusterKubeconfigV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "internal server error",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "unknown status",
			clientResponse: &mksclient.GetClusterKubeconfigV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Status:     http.StatusText(http.StatusServiceUnavailable),
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    http.StatusText(http.StatusServiceUnavailable),
			},
		},
		{
			name:        "http error",
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().GetClusterKubeconfigV2WithResponse(mock.Anything, clusterID).Return(test.clientResponse, test.clientError)

			kubeconfig, err := GetKubeconfig(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID)

			if test.errExpected != nil {
				require.Error(t, err)

				var mksErrExp *mksclient.MKSError
				if !errors.As(test.errExpected, &mksErrExp) {
					assert.ErrorIs(t, err, test.errExpected)

					return
				}

				var mksErr *mksclient.MKSError
				require.ErrorAs(t, err, &mksErr)
				assert.Equal(t, mksErrExp, mksErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.clientResponse.Body, kubeconfig)
		})
	}
}

func TestRotateCerts(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.RotateClusterCertsV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.RotateClusterCertsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNoContent,
					Status:     http.StatusText(http.StatusNoContent),
				},
				JSON204: testutils.Ptr(any(struct{}{})),
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.RotateClusterCertsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericNotFoundError{
					Error: struct {
						Id      string `json:"id"` //nolint:revive // it's generated struct
						Message string `json:"message"`
					}{
						Id:      clusterID,
						Message: "cluster not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "cluster not found",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.RotateClusterCertsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "internal server error",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "unknown status",
			clientResponse: &mksclient.RotateClusterCertsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Status:     http.StatusText(http.StatusServiceUnavailable),
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    http.StatusText(http.StatusServiceUnavailable),
			},
		},
		{
			name:        "http error",
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().RotateClusterCertsV2WithResponse(mock.Anything, clusterID).Return(test.clientResponse, test.clientError)

			err := RotateCerts(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID)

			if test.errExpected != nil {
				require.Error(t, err)

				var mksErrExp *mksclient.MKSError
				if !errors.As(test.errExpected, &mksErrExp) {
					assert.ErrorIs(t, err, test.errExpected)

					return
				}

				var mksErr *mksclient.MKSError
				require.ErrorAs(t, err, &mksErr)
				assert.Equal(t, mksErrExp, mksErr)

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUpgradePatchVersion(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.UpgradePatchVersionV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.UpgradePatchVersionV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.ClusterResp{
					Cluster: &mksclient.ClusterDetailed{
						Id:                            clusterID,
						Name:                          "test-cluster",
						KubeVersion:                   "1.28.1",
						Basic:                         false,
						EnableAutorepair:              true,
						EnablePatchVersionAutoUpgrade: true,
						KubeApiIp:                     "10.0.0.1",
						CreatedAt:                     time.Now(),
						CniType:                       mksclient.ClusterDetailedCniType("cilium"),
						NetworkType:                   mksclient.ClusterDetailedNetworkType("default"),
						Status:                        mksclient.ClusterDetailedStatus("active"),
						AdditionalSoftware:            map[string]any{"nginx-ingress": "enabled"},
					},
				},
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.UpgradePatchVersionV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericNotFoundError{
					Error: struct {
						Id      string `json:"id"` //nolint:revive // it's generated struct
						Message string `json:"message"`
					}{
						Id:      clusterID,
						Message: "cluster not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "cluster not found",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.UpgradePatchVersionV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "internal server error",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "unknown status",
			clientResponse: &mksclient.UpgradePatchVersionV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Status:     http.StatusText(http.StatusServiceUnavailable),
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    http.StatusText(http.StatusServiceUnavailable),
			},
		},
		{
			name:        "http error",
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().UpgradePatchVersionV2WithResponse(mock.Anything, clusterID).Return(test.clientResponse, test.clientError)

			cluster, err := UpgradePatchVersion(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID)

			if test.errExpected != nil {
				assert.Nil(t, cluster)
				require.Error(t, err)

				var mksErrExp *mksclient.MKSError
				if !errors.As(test.errExpected, &mksErrExp) {
					assert.ErrorIs(t, err, test.errExpected)

					return
				}

				var mksErr *mksclient.MKSError
				require.ErrorAs(t, err, &mksErr)
				assert.Equal(t, mksErrExp, mksErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.clientResponse.JSON200.Cluster, cluster)
		})
	}
}

func TestUpgradeMinorVersion(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.UpgradeMinorVersionV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.UpgradeMinorVersionV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.ClusterResp{
					Cluster: &mksclient.ClusterDetailed{
						Id:                            clusterID,
						Name:                          "test-cluster",
						KubeVersion:                   "1.29.0",
						Basic:                         false,
						EnableAutorepair:              true,
						EnablePatchVersionAutoUpgrade: true,
						KubeApiIp:                     "10.0.0.1",
						CreatedAt:                     time.Now(),
						CniType:                       mksclient.ClusterDetailedCniType("cilium"),
						NetworkType:                   mksclient.ClusterDetailedNetworkType("default"),
						Status:                        mksclient.ClusterDetailedStatus("active"),
						AdditionalSoftware:            map[string]any{"nginx-ingress": "enabled"},
					},
				},
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.UpgradeMinorVersionV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericNotFoundError{
					Error: struct {
						Id      string `json:"id"` //nolint:revive // it's generated struct
						Message string `json:"message"`
					}{
						Id:      clusterID,
						Message: "cluster not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "cluster not found",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.UpgradeMinorVersionV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "internal server error",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "unknown status",
			clientResponse: &mksclient.UpgradeMinorVersionV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Status:     http.StatusText(http.StatusServiceUnavailable),
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusServiceUnavailable,
				Message:    http.StatusText(http.StatusServiceUnavailable),
			},
		},
		{
			name:        "http error",
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().UpgradeMinorVersionV2WithResponse(mock.Anything, clusterID).Return(test.clientResponse, test.clientError)

			cluster, err := UpgradeMinorVersion(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID)

			if test.errExpected != nil {
				assert.Nil(t, cluster)
				require.Error(t, err)

				var mksErrExp *mksclient.MKSError
				if !errors.As(test.errExpected, &mksErrExp) {
					assert.ErrorIs(t, err, test.errExpected)

					return
				}

				var mksErr *mksclient.MKSError
				require.ErrorAs(t, err, &mksErr)
				assert.Equal(t, mksErrExp, mksErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.clientResponse.JSON200.Cluster, cluster)
		})
	}
}
