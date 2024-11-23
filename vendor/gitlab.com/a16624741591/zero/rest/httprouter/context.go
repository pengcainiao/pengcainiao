package httprouter

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/a16624741591/zero/core/logx"
)

var (
	traceIdKey = http.CanonicalHeaderKey("x-request-id")
)

// HeaderData header数据
type HeaderData struct {
	Authorization string `header:"Authorization" json:"-" `
	UserID        string `header:"X-Auth-User" binding:"required"  json:"user_id,omitempty"` //用户ID
	Platform      string `header:"X-Auth-Platform" json:"platform,omitempty"`                //客户端所属平台
	ClientVersion string `header:"X-Auth-Version" json:"client_version,omitempty"`           //客户端版本
	ClientIP      string `header:"X-Real-Ip" json:"client_ip,omitempty"`                     //客户端IP
	DeviceID      string `header:"X-AUTH-Device" json:"device_id,omitempty"`                 //设备ID
	RequestID     string `header:"X-Request-ID" json:"request_id,omitempty"`                 //分布式追踪ID，通过header传递
	Paging        `json:"-"`
}

type Paging struct {
	PageNumber int `form:"page_number,omitempty"`
	PageRecord int `form:"page_record,omitempty"`
}

type Context struct {
	context.Context `json:"-"`
	Data            HeaderData             `json:"-"`
	Keys            map[string]interface{} `json:"keys"`

	// This mutex protect Keys map
	mu sync.RWMutex
}

func (h *HeaderData) String() string {
	s, _ := jsoniter.MarshalToString(h)
	return s
}

func (p Paging) Offset() int {
	if p.PageNumber <= 0 {
		return 0
	}
	return (p.PageNumber - 1) * p.PageRecord
}

// check 检查分页参数
func (p *Paging) check() {
	if p.PageRecord == 0 {
		p.PageRecord = 15
	}
	if p.PageRecord > 20 {
		p.PageRecord = 20
	}
	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}
}

// NewContext 创建context
func NewContext(c *gin.Context) *Context {
	var (
		headerData HeaderData
		paging     Paging
	)
	if err := c.ShouldBindHeader(&headerData); err != nil {
		headerData.UserID = os.Getenv("X_AUTH_USER")
		headerData.Platform = os.Getenv("X_AUTH_PLATFORM")
		headerData.ClientVersion = os.Getenv("X_AUTH_VERSION")
		return &Context{
			Context: c.Request.Context(),
			Data:    headerData,
			mu:      sync.RWMutex{},
		}
	}

	_ = c.ShouldBindQuery(&paging)
	paging.check()
	headerData.Paging = paging
	//requestInjectSentrySpan(c.Request, headerData)
	if v := c.GetHeader(traceIdKey); v != "" {
		headerData.RequestID = v
	}
	return &Context{
		Context: c.Request.Context(),
		Data:    headerData,
		mu:      sync.RWMutex{},
	}
}

// NewContextData 创建带headerData Context
func NewContextData(ctx context.Context, headerData *HeaderData) *Context {
	if headerData != nil && headerData.RequestID != "" {
		ctx = context.WithValue(ctx, traceIdKey, headerData.RequestID)
	}
	newCtx := &Context{
		Context: ctx,
		mu:      sync.RWMutex{},
	}
	if headerData != nil {
		newCtx.Data = *headerData
	}
	return newCtx
}

func (c *Context) Error() *logx.TraceLoggerEvent {
	return logx.NewTraceLogger(c.Context).Error()
}

func (c *Context) Message() *logx.TraceLoggerEvent {
	return logx.NewTraceLogger(c.Context).Info()
}

func (c *Context) IsValidConn() bool {
	if c.Data.UserID == "" || c.Data.Platform == "" {
		return false
	}
	return true
}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
