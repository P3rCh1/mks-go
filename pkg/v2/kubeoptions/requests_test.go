package kubeoptions

import (
	"context"
	"errors"
	"net/http"
	"testing"

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

func TestListFeatureGates(t *testing.T) {
	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.ListFeatureGatesV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.ListFeatureGatesV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.FeatureGatesList{
					FeatureGates: &[]mksclient.AvailableFeatureGates{
						{
							KubeVersionMinor: ptr("1.28"),
							Names:            ptr(mksclient.OptionNamesFG{"FeatureGate1", "FeatureGate2"}),
						},
					},
				},
			},
		},
		{
			name: "empty feature gates",
			clientResponse: &mksclient.ListFeatureGatesV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.FeatureGatesList{},
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.ListFeatureGatesV2Response{
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
			clientResponse: &mksclient.ListFeatureGatesV2Response{
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
			mksClient.EXPECT().ListFeatureGatesV2WithResponse(mock.Anything).Return(test.clientResponse, test.clientError)

			featureGates, err := ListFeatureGates(context.Background(), &v2.ServiceClient{MKSClient: mksClient})

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

			if test.clientResponse.JSON200.FeatureGates != nil {
				assert.Equal(t, *test.clientResponse.JSON200.FeatureGates, featureGates)
			} else {
				assert.Nil(t, featureGates)
			}
		})
	}
}

func TestListAdmissionControllers(t *testing.T) {
	httpError := errors.New("error")

	tests := []struct {
		name           string
		clientResponse *mksclient.ListAdmissionControllersV2Response
		clientError    error
		errExpected    error
	}{
		{
			name: "success",
			clientResponse: &mksclient.ListAdmissionControllersV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.AdmissionControllersList{
					AdmissionControllers: &[]mksclient.AvailableAdmissionControllers{
						{
							KubeVersionMinor: ptr("1.28"),
							Names:            ptr(mksclient.OptionNamesAC{"AC1", "AC2"}),
						},
					},
				},
			},
		},
		{
			name: "empty admission controllers",
			clientResponse: &mksclient.ListAdmissionControllersV2Response{
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
				},
				JSON200: &mksclient.AdmissionControllersList{},
			},
		},
		{
			name: "internal server error",
			clientResponse: &mksclient.ListAdmissionControllersV2Response{
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
			clientResponse: &mksclient.ListAdmissionControllersV2Response{
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
			mksClient.EXPECT().ListAdmissionControllersV2WithResponse(mock.Anything).Return(test.clientResponse, test.clientError)

			admissionControllers, err := ListAdmissionControllers(context.Background(), &v2.ServiceClient{MKSClient: mksClient})

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

			if test.clientResponse.JSON200.AdmissionControllers != nil {
				assert.Equal(t, *test.clientResponse.JSON200.AdmissionControllers, admissionControllers)
			} else {
				assert.Nil(t, admissionControllers)
			}
		})
	}
}
