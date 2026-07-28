package node

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
		nodeID      = "test-node-id"
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.GetNodeV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.GetNodeV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.NodeResp{
					Node: mksclient.Node{
						CreatedAt:        time.Now(),
						Hostname:         "test-hostname",
						Id:               nodeID,
						Ip:               "10.0.0.1",
						NodegroupId:      nodegroupID,
						ProviderServerId: "provider-server-id",
						UpdatedAt:        ptr(time.Now()),
					},
				},
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.GetNodeV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "node not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "node not found",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.GetNodeV2Response{
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
			clientResponse: &mksclient.GetNodeV2Response{
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
			mksClient.EXPECT().GetNodeV2WithResponse(mock.Anything, clusterID, nodegroupID, nodeID).Return(test.clientResponse, test.clientError)

			node, err := Get(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nodegroupID, nodeID)

			if test.errExpected != nil {
				assert.Nil(t, node)
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

			assert.Equal(t, test.clientResponse.JSON200.Node.Id, node.Id)
			assert.Equal(t, test.clientResponse.JSON200.Node.Hostname, node.Hostname)
			assert.Equal(t, test.clientResponse.JSON200.Node.NodegroupId, node.NodegroupId)
		})
	}
}

func TestReinstall(t *testing.T) {
	const (
		clusterID   = "test-cluster-id"
		nodegroupID = "test-nodegroup-id"
		nodeID      = "test-node-id"
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.ReinstallNodeV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.ReinstallNodeV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNoContent,
					Status:     http.StatusText(http.StatusNoContent),
				},
				JSON204: ptr(any(struct{}{})),
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.ReinstallNodeV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "node not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "node not found",
			},
		},
		{
			name: "bad request",
			clientResponse: &mksclient.ReinstallNodeV2Response{
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
			clientResponse: &mksclient.ReinstallNodeV2Response{
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
			clientResponse: &mksclient.ReinstallNodeV2Response{
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
			clientResponse: &mksclient.ReinstallNodeV2Response{
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
			mksClient.EXPECT().ReinstallNodeV2WithResponse(mock.Anything, clusterID, nodegroupID, nodeID).Return(test.clientResponse, test.clientError)

			err := Reinstall(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nodegroupID, nodeID)

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
		nodeID      = "test-node-id"
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.DeleteNodeV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.DeleteNodeV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNoContent,
					Status:     http.StatusText(http.StatusNoContent),
				},
				JSON204: ptr(any(struct{}{})),
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.DeleteNodeV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: "node not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "node not found",
			},
		},
		{
			name: "bad request",
			clientResponse: &mksclient.DeleteNodeV2Response{
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
			clientResponse: &mksclient.DeleteNodeV2Response{
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
			clientResponse: &mksclient.DeleteNodeV2Response{
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
			clientResponse: &mksclient.DeleteNodeV2Response{
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
			mksClient.EXPECT().DeleteNodeV2WithResponse(mock.Anything, clusterID, nodegroupID, nodeID).Return(test.clientResponse, test.clientError)

			err := Delete(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, nodegroupID, nodeID)

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
