package main

import (
	"github.com/pengcainiao2/usercenter/internal"
	"github.com/pengcainiao2/zero/rest/httprouter"
)

func main() {
	httprouter.RecoveryMainAsync(func() {
		internal.Setup()
	})
}
