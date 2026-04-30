package mksclient

import "fmt"

// errGotHTTPStatusCodeFmt is an error message format for unexpected HTTP status codes.
const errGotHTTPStatusCodeFmt = "mks-go: got the %d status code with message: %s"

type MKSError struct {
	StatusCode int
	Message    string
}

func (e *MKSError) Status() int {
	return e.StatusCode
}

func (e *MKSError) Error() string {
	return e.Message
}

func (e *GenericError) ToMKSError(statusCode int) *MKSError {
	return &MKSError{
		StatusCode: statusCode,
		Message:    e.Error.Message,
	}
}

func (e *GenericNotFoundError) ToMKSError(statusCode int) *MKSError {
	return &MKSError{
		StatusCode: statusCode,
		Message:    e.Error.Message,
	}
}

type APIError interface {
	ToMKSError(statusCode int) *MKSError
}

func HandleAPIErrors(statusCode int, statusMsg string, errors ...APIError) error {
	for _, err := range errors {
		if err != nil {
			return err.ToMKSError(statusCode)
		}
	}

	return fmt.Errorf(errGotHTTPStatusCodeFmt, statusCode, statusMsg)
}
