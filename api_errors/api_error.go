package api_error

var CodeName = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	500: "Internal Server Error",
	503: "Service Unavailable",
}

type Code int

func (code Code) Name() string {
	return CodeName[int(code)]
}

type APIError struct {
	Code
	Msg string
}

func NewAPIError(code Code, msg string) *APIError {
	return &APIError{
		Code: code,
		Msg:  msg,
	}
}

func (e *APIError) Error() string {
	return e.Msg
}
