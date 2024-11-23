package grpcbase

import jsoniter "github.com/json-iterator/go"

// Response
type Response struct {
	Total   int         `json:"total,omitempty"`   //总数
	Message string      `json:"message,omitempty"` //错误消息
	Data    interface{} `json:"data,omitempty"`    //存储任意结构
}

func ErrorResponse(err error) Response {
	return Response{Message: err.Error()}
}

func (r Response) String() string {
	b, _ := jsoniter.Marshal(r)
	return string(b)
}
