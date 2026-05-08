package nodegroup

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	v2 "github.com/selectel/mks-go/pkg/v2"
	"github.com/selectel/mks-go/pkg/v2/mksclient"
	mksmock "github.com/selectel/mks-go/pkg/v2/mksclient/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func ptr[T any](value T) *T {
	return &value
}

func TestGet(t *testing.T) {
	const (
		clusterID   = "test-cluster-id"
		nodegroupID = "test-nodegroup-id"
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.GetNodegroupV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.GetNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.NodegroupResp{
					Nodegroup: mksclient.NodegroupDetailed{
						AutoscaleMaxNodes: ptr(int64(10)),
						AutoscaleMinNodes: ptr(int64(1)),
						Cidr:              ptr("192.168.1.0/24"),
						CloudNodegroupConfig: &mksclient.CloudNodegroupConfig{
							AffinityPolicy: "anti-affinity",
							Cpus:           4,
							FlavorId:       "test-flavor",
							KeypairName:    "test-key",
							LocalVolume:    false,
							RamMb:          8192,
							VolumeGb:       100,
							VolumeType:     "standard",
						},
						ClusterId:                 clusterID,
						CreatedAt:                 time.Now(),
						DedicatedNodegroupConfig:  nil,
						EnableAutoscale:           true,
						Id:                        nodegroupID,
						InstallNvidiaDevicePlugin: true,
						Labels:                    map[string]string{"key": "value"},
						NodegroupType:             mksclient.NodegroupDetailedNodegroupType("standard"),
						Nodes: []mksclient.Node{{
							Id:          "node-1",
							Hostname:    "node-1",
							Ip:          "10.0.0.1",
							NodegroupId: nodegroupID,
							CreatedAt:   time.Now(),
						}},
						Preemptible: false,
						Segment:     "default-segment",
						Status:      mksclient.NodegroupDetailedStatus("active"),
						Taints: []mksclient.NodegroupTaint{{
							Key:    "key",
							Value:  "value",
							Effect: mksclient.NoSchedule,
						}},
						UpdatedAt: ptr(time.Now()),
						UserData:  "user data",
					},
				},
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.GetNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericNotFoundError{
					Error: struct {
						Id      string `json:"id"` //nolint:revive // it's generated struct
						Message string `json:"message"`
					}{
						Id:      nodegroupID,
						Message: "nodegroup not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "nodegroup not found",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.GetNodegroupV2Response{
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
			clientResponse: &mksclient.GetNodegroupV2Response{
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
			mksClient.EXPECT().GetNodegroupV2WithResponse(mock.Anything, clusterID, nodegroupID).Return(test.clientResponse, test.clientError)

			nodegroup, err := Get(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nodegroupID)

			if test.errExpected != nil {
				assert.Nil(t, nodegroup)
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

			assert.Equal(t, test.clientResponse.JSON200.Nodegroup.Id, nodegroup.Id)
			assert.Equal(t, test.clientResponse.JSON200.Nodegroup.ClusterId, nodegroup.ClusterId)
		})
	}
}

func TestList(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.ListNodegroupsV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.ListNodegroupsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.NodegroupList{
					Nodegroups: []mksclient.NodegroupListItem{
						{
							AutoscaleMaxNodes:       ptr(int64(10)),
							AutoscaleMinNodes:       ptr(int64(1)),
							AvailableAdditionalInfo: mksclient.NodegorupAdditionalInfo{UserData: true},
							Cidr:                    ptr("192.168.1.0/24"),
							CloudNodegroupConfig: &mksclient.CloudNodegroupConfig{
								AffinityPolicy: "anti-affinity",
								Cpus:           4,
								FlavorId:       "test-flavor",
								KeypairName:    "test-key",
								LocalVolume:    false,
								RamMb:          8192,
								VolumeGb:       100,
								VolumeType:     "standard",
							},
							ClusterId:                 clusterID,
							CreatedAt:                 time.Now(),
							DedicatedNodegroupConfig:  nil,
							EnableAutoscale:           true,
							Id:                        "test-nodegroup-id-1",
							InstallNvidiaDevicePlugin: true,
							Labels:                    map[string]string{"key": "value"},
							NodegroupType:             mksclient.STANDARD,
							Nodes: []mksclient.Node{{
								Id:          "node-1",
								Hostname:    "node-1",
								Ip:          "10.0.0.1",
								NodegroupId: "test-nodegroup-id-1",
								CreatedAt:   time.Now(),
							}},
							Preemptible: false,
							Segment:     "default-segment",
							Status:      mksclient.NodegroupListItemStatus("active"),
							Taints: []mksclient.NodegroupTaint{{
								Key:    "key",
								Value:  "value",
								Effect: mksclient.NoSchedule,
							}},
							UpdatedAt: ptr(time.Now()),
						},
					},
				},
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.ListNodegroupsV2Response{
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
			name: "bad request",
			clientResponse: &mksclient.ListNodegroupsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     http.StatusText(http.StatusBadRequest),
				},
				JSON400: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "bad request",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusBadRequest,
				Message:    "bad request",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.ListNodegroupsV2Response{
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
			clientResponse: &mksclient.ListNodegroupsV2Response{
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
			mksClient.EXPECT().ListNodegroupsV2WithResponse(mock.Anything, clusterID).Return(test.clientResponse, test.clientError)

			nodegroups, err := List(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID)

			if test.errExpected != nil {
				assert.Nil(t, nodegroups)
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
			assert.Equal(t, test.clientResponse.JSON200.Nodegroups, nodegroups)
		})
	}
}

func TestCreate(t *testing.T) {
	const clusterID = "test-cluster-id"

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.CreateNodegroupsV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: func() *mksclient.CreateNodegroupsV2Response {
				var v any = struct{}{}

				return &mksclient.CreateNodegroupsV2Response{
					HTTPResponse: &http.Response{
						StatusCode: http.StatusNoContent,
						Status:     http.StatusText(http.StatusNoContent),
					},
					JSON204: &v,
				}
			}(),
		},
		{
			name: "not found",
			clientResponse: &mksclient.CreateNodegroupsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
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
			name: "bad request",
			clientResponse: &mksclient.CreateNodegroupsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     http.StatusText(http.StatusBadRequest),
				},
				JSON400: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "bad request",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusBadRequest,
				Message:    "bad request",
			},
		},
		{
			name: "conflict",
			clientResponse: &mksclient.CreateNodegroupsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusConflict,
					Status:     http.StatusText(http.StatusConflict),
				},
				JSON409: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "conflict",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusConflict,
				Message:    "conflict",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.CreateNodegroupsV2Response{
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
			clientResponse: &mksclient.CreateNodegroupsV2Response{
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
			mksClient.EXPECT().CreateNodegroupsV2WithResponse(mock.Anything, clusterID, mock.Anything).Return(test.clientResponse, test.clientError)

			err := Create(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nil)

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

func TestDelete(t *testing.T) {
	const (
		clusterID   = "test-cluster-id"
		nodegroupID = "test-nodegroup-id"
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.DeleteNodegroupV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: func() *mksclient.DeleteNodegroupV2Response {
				var v any = struct{}{}

				return &mksclient.DeleteNodegroupV2Response{
					HTTPResponse: &http.Response{
						StatusCode: http.StatusNoContent,
						Status:     http.StatusText(http.StatusNoContent),
					},
					JSON204: &v,
				}
			}(),
		},
		{
			name: "not found",
			clientResponse: &mksclient.DeleteNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "nodegroup not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "nodegroup not found",
			},
		},
		{
			name: "bad request",
			clientResponse: &mksclient.DeleteNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     http.StatusText(http.StatusBadRequest),
				},
				JSON400: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "bad request",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusBadRequest,
				Message:    "bad request",
			},
		},
		{
			name: "conflict",
			clientResponse: &mksclient.DeleteNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusConflict,
					Status:     http.StatusText(http.StatusConflict),
				},
				JSON409: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "conflict",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusConflict,
				Message:    "conflict",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.DeleteNodegroupV2Response{
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
			clientResponse: &mksclient.DeleteNodegroupV2Response{
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
			mksClient.EXPECT().DeleteNodegroupV2WithResponse(mock.Anything, clusterID, nodegroupID).Return(test.clientResponse, test.clientError)

			err := Delete(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nodegroupID)

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

func TestResize(t *testing.T) {
	const (
		clusterID   = "test-cluster-id"
		nodegroupID = "test-nodegroup-id"
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.ResizeNodegroupV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: func() *mksclient.ResizeNodegroupV2Response {
				var v any = struct{}{}

				return &mksclient.ResizeNodegroupV2Response{
					HTTPResponse: &http.Response{
						StatusCode: http.StatusNoContent,
						Status:     http.StatusText(http.StatusNoContent),
					},
					JSON204: &v,
				}
			}(),
		},
		{
			name: "not found",
			clientResponse: &mksclient.ResizeNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "nodegroup not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "nodegroup not found",
			},
		},
		{
			name: "bad request",
			clientResponse: &mksclient.ResizeNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     http.StatusText(http.StatusBadRequest),
				},
				JSON400: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "bad request",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusBadRequest,
				Message:    "bad request",
			},
		},
		{
			name: "conflict",
			clientResponse: &mksclient.ResizeNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusConflict,
					Status:     http.StatusText(http.StatusConflict),
				},
				JSON409: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "conflict",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusConflict,
				Message:    "conflict",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.ResizeNodegroupV2Response{
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
			clientResponse: &mksclient.ResizeNodegroupV2Response{
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
			mksClient.EXPECT().ResizeNodegroupV2WithResponse(mock.Anything, clusterID, nodegroupID, mock.Anything).Return(test.clientResponse, test.clientError)

			err := Resize(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nodegroupID, 0)

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

func TestUpdate(t *testing.T) {
	const (
		clusterID   = "test-cluster-id"
		nodegroupID = "test-nodegroup-id"
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.UpdateNodegroupV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: func() *mksclient.UpdateNodegroupV2Response {
				var v any = struct{}{}

				return &mksclient.UpdateNodegroupV2Response{
					HTTPResponse: &http.Response{
						StatusCode: http.StatusNoContent,
						Status:     http.StatusText(http.StatusNoContent),
					},
					JSON204: &v,
				}
			}(),
		},
		{
			name: "not found",
			clientResponse: &mksclient.UpdateNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "nodegroup not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "nodegroup not found",
			},
		},
		{
			name: "bad request",
			clientResponse: &mksclient.UpdateNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     http.StatusText(http.StatusBadRequest),
				},
				JSON400: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "bad request",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusBadRequest,
				Message:    "bad request",
			},
		},
		{
			name: "conflict",
			clientResponse: &mksclient.UpdateNodegroupV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusConflict,
					Status:     http.StatusText(http.StatusConflict),
				},
				JSON409: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "conflict",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusConflict,
				Message:    "conflict",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.UpdateNodegroupV2Response{
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
			clientResponse: &mksclient.UpdateNodegroupV2Response{
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
			mksClient.EXPECT().UpdateNodegroupV2WithResponse(mock.Anything, clusterID, nodegroupID, mock.Anything).Return(test.clientResponse, test.clientError)

			err := Update(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nodegroupID, mksclient.NodegroupUpdateStruct{})

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
