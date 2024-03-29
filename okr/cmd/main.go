package main

import (
	"github.com/pengcainiao/pengcainiao/okr/internal"
	"github.com/pengcainiao/zero/core/sysx"
	"github.com/pengcainiao/zero/rest/httprouter"
)

func main() {
	sysx.SubSystem = "okr"
	httprouter.RecoveryMainAsync(func() {
		internal.Setup()
	})
}
