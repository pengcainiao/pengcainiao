package httprouter

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Response struct {
	Code             int         `json:"code"`                      //响应编码
	InternalError    string      `json:"dbg_error,omitempty"`       //内部错误信息，详情但不对外输出
	RequestParameter interface{} `json:"dbg_parameter,omitempty"`   //请求携带的参数，仅用于内部查看
	Message          string      `json:"message,omitempty"`         // 错误消息
	Data             interface{} `json:"data,omitempty"`            //有效荷载
	ScrollID         string      `json:"scroll_id,omitempty"`       // 滚动id
	CompleteTotal    int         `json:"complete_total,omitempty"`  // 完成事项总数
	NotTodayTotal    int         `json:"not_today_total,omitempty"` // 非今天明确事项总数
	Total            int         `json:"total,omitempty"`           //总数
	Cursor           int         `json:"cursor,omitempty"`          //游标
}

// SetRequestParameter 设置请求参数
func (r Response) SetRequestParameter(parameter interface{}) Response {
	r.RequestParameter = parameter
	return r
}

func (r Response) SetOuterErrorMessage(errMessage string) Response {
	if errMessage == "" {
		return r
	}
	r.Message = errMessage
	return r
}

func (r Response) String() string {
	v, _ := jsoniter.MarshalToString(r)
	return v
}

// Success 成功返回
func Success(maps ...map[string]interface{}) Response {
	res := Response{Code: 0}
	if len(maps) > 0 {
		if data, ok := maps[0]["data"]; ok {
			res.Data = data
		}

		if total, ok := maps[0]["total"]; ok {
			res.Total = int(total.(int64))
		}

		if completeTotal, ok := maps[0]["complete_total"]; ok {
			res.CompleteTotal = int(completeTotal.(int64))
		}

		if notTodayTotal, ok := maps[0]["not_today_total"]; ok {
			res.NotTodayTotal = int(notTodayTotal.(int64))
		}

		if scrollID, ok := maps[0]["scroll_id"]; ok {
			res.ScrollID = scrollID.(string)
		}

		if cursor, ok := maps[0]["cursor"]; ok {
			res.Cursor = cursor.(int)
		}
	}

	return res
}

// ErrorSame 内部错误与外部错误相同时使用
func ErrorSame(code int, err string) Response {
	return Error(code, err, err)
}

// Error error，外部错误简单不暴露实现细节，内部错误应为真实错误
func Error(code int, outError string, internalError interface{}) Response {
	switch e := internalError.(type) {
	case error:
		return Response{Code: code, Message: outError, InternalError: e.Error()}
	case string:
		return Response{Code: code, Message: outError, InternalError: e}
	}
	return Response{Code: code}
}

// HttpCode 标准化HTTP状态码
func HttpCode(customResponseCode int) int {
	if customResponseCode == -1 {
		return 500
	} else if customResponseCode == 0 {
		return 200
	} else {
		var tempCode = customResponseCode
		for tempCode = tempCode / 10; tempCode > 999; {
			tempCode = tempCode / 10
		}
		if v := http.StatusText(tempCode); v == "" {
			panic(fmt.Sprintf("invalid custom response code %d", customResponseCode))
		}
		return tempCode
	}
}

// ResponseJSONContent 返回json
func ResponseJSONContent(c *gin.Context, resp interface{}) {
	span := trace.SpanFromContext(c.Request.Context())
	switch v := resp.(type) {
	case string:
		span.AddEvent("response", trace.WithAttributes(
			attribute.String("response.data", v),
		))
		c.Header("Content-Type", "application/json")
		c.Writer.WriteHeader(200)
		_, _ = c.Writer.WriteString(v)
	case Response:
		if v.InternalError != "" {
			span.SetStatus(codes.Error, v.InternalError)
		}
		span.AddEvent("response", trace.WithAttributes(
			attribute.String("response.data", v.String()),
		))

		if v.Code != 0 || c.Writer.Status() != http.StatusOK {
			header, _ := c.Get("header")
			logEvent := logx.NewTraceLogger(c.Request.Context()).Err(errors.New(v.InternalError)).
				Str("url", c.Request.URL.String()).
				Str("method", c.Request.Method).
				Interface("header", header)

			if request, ok := c.Get("request"); ok {
				if req, ok := request.([]byte); ok && len(req) > 0 {
					logEvent.RawJSON("request", req)
				}
			}

			logEvent.Interface("response", v).Msg("HTTP响应请求返回失败")
		}

		if env.IsProduction() {
			v.InternalError = ""
			v.RequestParameter = nil
		}
		c.JSON(HttpCode(v.Code), v)
	}
}
