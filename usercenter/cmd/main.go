package main

import (
	"github.com/pengcainiao/pengcainiao/usercenter/internal"
	"github.com/pengcainiao/zero/rest/httprouter"
)

func main() {
	httprouter.RecoveryMainAsync(func() {
		internal.Setup()
	})
}
