package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type wrappedResponseWriter struct {
	gin.ResponseWriter
	writer http.ResponseWriter
}

func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *wrappedResponseWriter) WriteString(s string) (n int, err error) {
	return w.writer.Write([]byte(s))
}

// An http.Handler that passes on calls to downstream middlewares
type nextRequestHandler struct {
	c *gin.Context
}

// Run the next request in the middleware chain and return
func (h *nextRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.c.Writer = &wrappedResponseWriter{h.c.Writer, w}
	h.c.Next()
}

func WrapHttpHandler(hh ...func(h http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, f := range hh {
			f(&nextRequestHandler{c}).ServeHTTP(c.Writer, c.Request)
		}
	}
}
