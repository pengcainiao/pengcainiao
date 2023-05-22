
package main

import (
	"github.com/pengcainiao/pengcainiao/usercenter/internal"
	"github.com/pengcainiao/zero/core/sysx"
	"github.com/pengcainiao/zero/rest/httprouter"
)

func main() {
	sysx.SubSystem = "user"
	httprouter.RecoveryMainAsync(func() {
		internal.Setup()
	})
}
