package httpx

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/963204765/httpclient/httplib"
	"gitlab.com/a16624741591/zero/core/mapping"
	"gitlab.com/a16624741591/zero/rest/httprouter"
	"gitlab.com/a16624741591/zero/rest/pathvar"
)

const (
	formKey           = "form"
	pathKey           = "path"
	headerKey         = "header"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

var (
	formUnmarshaler   = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
	pathUnmarshaler   = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())
	headerUnmarshaler = mapping.NewUnmarshaler(headerKey, mapping.WithStringValues(),
		mapping.WithCanonicalKeyFunc(textproto.CanonicalMIMEHeaderKey))
)

// SendHTTP 发送http请求
// Deprecated: 使用 httprouter.PerformanceRequest 代替，新的方法中集成了sentry
func SendHTTP(ctx *httprouter.Context, method, url string) *httplib.BeegoHTTPRequest {
	req := httplib.NewBeegoRequest(url, method)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	header := ctx.Data
	req.Header("Authorization", header.Authorization)
	req.Header("X-Auth-User", header.UserID)
	req.Header("X-Auth-Platform", header.Platform)
	req.Header("X-Auth-Version", header.ClientVersion)
	req.Header("X-Auth-ClientIP", header.ClientIP)
	req.Header("X-Request-ID", header.RequestID)
	return req
}

// PerformanceRequest 发起HTTP请求
func PerformanceRequest(ctx *httprouter.Context, req *http.Request) (data []byte, err error) {
	header := ctx.Data
	req.Header.Add("Authorization", header.Authorization)
	req.Header.Add("X-Auth-User", header.UserID)
	req.Header.Add("X-Auth-Platform", header.Platform)
	req.Header.Add("X-Auth-Version", header.ClientVersion)
	req.Header.Add("X-Auth-ClientIP", header.ClientIP)
	req.Header.Add("X-Request-ID", header.RequestID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err = errors.New(string(b))
		return nil, err
	}
	return b, nil
}

// Parse parses the request.
func Parse(r *http.Request, v interface{}) error {
	if err := ParsePath(r, v); err != nil {
		return err
	}

	if err := ParseForm(r, v); err != nil {
		return err
	}

	if err := ParseHeaders(r, v); err != nil {
		return err
	}

	return ParseJsonBody(r, v)
}

// ParseHeaders parses the headers request.
func ParseHeaders(r *http.Request, v interface{}) error {
	m := map[string]interface{}{}
	for k, v := range r.Header {
		if len(v) == 1 {
			m[k] = v[0]
		} else {
			m[k] = v
		}
	}

	return headerUnmarshaler.Unmarshal(m, v)
}

// ParseForm parses the form request.
func ParseForm(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}

	params := make(map[string]interface{}, len(r.Form))
	for name := range r.Form {
		formValue := r.Form.Get(name)
		if len(formValue) > 0 {
			params[name] = formValue
		}
	}

	return formUnmarshaler.Unmarshal(params, v)
}

// ParseHeader parses the request header and returns a map.
func ParseHeader(headerValue string) map[string]string {
	ret := make(map[string]string)
	fields := strings.Split(headerValue, separator)

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}

		kv := strings.SplitN(field, "=", tokensInAttribute)
		if len(kv) != tokensInAttribute {
			continue
		}

		ret[kv[0]] = kv[1]
	}

	return ret
}

// ParseJsonBody parses the post request which contains json in body.
func ParseJsonBody(r *http.Request, v interface{}) error {
	if withJsonBody(r) {
		reader := io.LimitReader(r.Body, maxBodyLen)
		return mapping.UnmarshalJsonReader(reader, v)
	}

	return mapping.UnmarshalJsonMap(nil, v)
}

// ParsePath parses the symbols reside in url path.
// Like http://localhost/bag/:name
func ParsePath(r *http.Request, v interface{}) error {
	vars := pathvar.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(ContentType), ApplicationJson)
}
