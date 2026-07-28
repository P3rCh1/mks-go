package task

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

func TestGet(t *testing.T) {
	const (
		clusterID = "test-cluster-id"
		taskID    = "test-task-id"
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.GetTaskV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.GetTaskV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.TaskResp{
					Task: mksclient.Task{
						ClusterId:   clusterID,
						Id:          taskID,
						NodegroupId: ptr("test-nodegroup-id"),
						StartedAt:   time.Now(),
						Status:      mksclient.INPROGRESS,
						Type:        "test-type",
						UpdatedAt:   ptr(time.Now()),
						ErrorDetails: &mksclient.TaskErrorDetails{
							Code:    1,
							Details: ptr("test-error"),
							Name:    "TEST_ERROR",
						},
					},
				},
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.GetTaskV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				JSON404: &mksclient.GenericNotFoundError{
					Error: struct {
						Id      string `json:"id"` //nolint:revive // it's generated struct
						Message string `json:"message"`
					}{
						Id:      taskID,
						Message: "task not found",
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusNotFound,
				Message:    "task not found",
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.GetTaskV2Response{
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
			clientResponse: &mksclient.GetTaskV2Response{
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
			mksClient.EXPECT().GetTaskV2WithResponse(mock.Anything, clusterID, taskID, mock.Anything).Return(test.clientResponse, test.clientError)

			task, err := Get(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, taskID)

			if test.errExpected != nil {
				assert.Nil(t, task)
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

			assert.Equal(t, test.clientResponse.JSON200.Task.Id, task.Id)
			assert.Equal(t, test.clientResponse.JSON200.Task.ClusterId, task.ClusterId)
		})
	}
}

func TestList(t *testing.T) {
	const (
		clusterID = "test-cluster-id"
		limit     = 10
		offset    = 1
	)

	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.ListTasksV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.ListTasksV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.TaskList{
					Count: 2,
					Tasks: []mksclient.Task{
						{
							ClusterId:   clusterID,
							Id:          "test-task-id-1",
							NodegroupId: ptr("test-nodegroup-id"),
							StartedAt:   time.Now(),
							Status:      mksclient.INPROGRESS,
							Type:        "test-type",
							UpdatedAt:   ptr(time.Now()),
							ErrorDetails: &mksclient.TaskErrorDetails{
								Code:    1,
								Details: ptr("test-error"),
								Name:    "TEST_ERROR",
							},
						},
						{
							ClusterId:   clusterID,
							Id:          "test-task-id-2",
							NodegroupId: ptr("test-nodegroup-id"),
							StartedAt:   time.Now(),
							Status:      mksclient.DONE,
							Type:        "test-type",
							UpdatedAt:   ptr(time.Now()),
							ErrorDetails: &mksclient.TaskErrorDetails{
								Code:    2,
								Details: ptr("test-error"),
								Name:    "TEST_ERROR",
							},
						},
					},
				},
			},
		},
		{
			name: "not found",
			clientResponse: &mksclient.ListTasksV2Response{
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
			clientResponse: &mksclient.ListTasksV2Response{
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
			clientResponse: &mksclient.ListTasksV2Response{
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
			mksClient.EXPECT().ListTasksV2WithResponse(mock.Anything, clusterID, mock.Anything).Return(test.clientResponse, test.clientError)

			tasks, err := List(context.Background(), &v2.ServiceClient{MKSClient: mksClient}, clusterID, uint64(limit), uint64(offset))

			if test.errExpected != nil {
				assert.Nil(t, tasks)
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
			assert.Equal(t, test.clientResponse.JSON200.Tasks, tasks)
		})
	}
}
