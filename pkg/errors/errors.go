package apierrors

type DataNotFoundErrorWrapper struct {
	ReturnedStatusCode int
	OriginalError      error
}

// Error the original error message remains as it is for logging reasons etc.
// and the wrapper error message is empty because we don't want the client to see anything
func (err DataNotFoundErrorWrapper) Error() string {
	return ""
}

func (err DataNotFoundErrorWrapper) Unwrap() error {
	return err.OriginalError
}

type MalformedEventDataError struct {
	ReturnedStatusCode   int
	OriginalError        error
	MalformedParameterId string
}

// Error the original error message remains as it is for logging reasons etc.
// and the wrapper error message is empty because we don't want the client to see anything
func (err MalformedEventDataError) Error() string {
	return "malformed solar panel data, check parameter " + err.MalformedParameterId
}

func (err MalformedEventDataError) Unwrap() error {
	return err.OriginalError
}

type EmptySolarDataError struct {
	ReturnedStatusCode int
}

// Error the original error message remains as it is for logging reasons etc.
// and the wrapper error message is empty because we don't want the client to see anything
func (err EmptySolarDataError) Error() string {
	return "solar data is empty on request"
}
