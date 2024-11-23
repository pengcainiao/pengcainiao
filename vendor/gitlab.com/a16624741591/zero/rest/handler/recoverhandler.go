package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// RecoverHandler returns a middleware that recovers if panic happens.
func RecoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if result := recover(); result != nil {
				// 妈的，这里太坑，看不出错误
				//internal.Error(r, fmt.Sprintf("%v\n%s", result, debug.Stack()))
				fmt.Println(result)
				fmt.Println(string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
