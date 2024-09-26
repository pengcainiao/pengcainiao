package apierror

type ErrorWithData interface {
	Error
	Data() interface{}
}

func NewApiErrWithData(code int, msg string, data interface{}) ErrorWithData {
	return &ApiErrWithData{Error: NewApiError(code, msg), data: data}
}

type ApiErrWithData struct {
	Error
	data interface{}
}

func (e *ApiErrWithData) Data() interface{} { return e.data }
