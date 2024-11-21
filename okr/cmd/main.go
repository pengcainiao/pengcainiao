package main

import (
	"gitlab.com/a16624741591/zero/rest/httprouter"
	"pp/okr/internal"
)

func main() {
	httprouter.RecoveryMainAsync(func() {
		internal.Setup()
	})
}
