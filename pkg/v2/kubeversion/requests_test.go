package kubeversion

import (
	"context"
	"errors"
	"net/http"
	"testing"

	v2 "github.com/selectel/mks-go/pkg/v2"
	"github.com/selectel/mks-go/pkg/v2/mksclient"
	mksmock "github.com/selectel/mks-go/pkg/v2/mksclient/mocks"
	"github.com/selectel/mks-go/pkg/v2/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	httpError := errors.New("error")

	versions := []mksclient.KubeVersionInfo{
		{
			Version:   testutils.Ptr("v1.27.0"),
			IsDefault: testutils.Ptr(false),
		},
		{
			Version:   testutils.Ptr("v1.28.0"),
			IsDefault: testutils.Ptr(true),
		},
	}

	tests := []struct {
		name           string
		clientResponse *mksclient.ListKubeVersionsV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: testutils.NameSuccess,
			clientResponse: &mksclient.ListKubeVersionsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.KubeVersionsList{
					KubeVersions: &versions,
				},
			},
		},
		{
			name: testutils.NameEmptyKubeVersions,
			clientResponse: &mksclient.ListKubeVersionsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.KubeVersionsList{},
			},
		},
		{
			name: testutils.NameInternalError,
			clientResponse: &mksclient.ListKubeVersionsV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
				},
				JSON500: &mksclient.GenericError{
					Error: struct {
						Message string `json:"message"`
					}{
						Message: testutils.MsgInternalError,
					},
				},
			},
			errExpected: &mksclient.MKSError{
				StatusCode: http.StatusInternalServerError,
				Message:    testutils.MsgInternalError,
			},
		},
		{
			name: testutils.NameUnknownStatus,
			clientResponse: &mksclient.ListKubeVersionsV2Response{
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
			name:        testutils.NameHTTPError,
			clientError: httpError,
			errExpected: httpError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mksClient := mksmock.NewMockClientWithResponsesInterface(t)
			mksClient.EXPECT().ListKubeVersionsV2WithResponse(mock.Anything).Return(test.clientResponse, test.clientError)

			kubeVersions, err := List(context.Background(), &v2.ServiceClient{MKSClient: mksClient})

			if test.errExpected != nil {
				assert.Nil(t, kubeVersions)
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

			if test.clientResponse.JSON200 != nil && test.clientResponse.JSON200.KubeVersions != nil {
				assert.Equal(t, *test.clientResponse.JSON200.KubeVersions, kubeVersions)
			} else {
				assert.Empty(t, kubeVersions)
			}
		})
	}
}
