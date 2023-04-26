package main

func main() {
	httprouter.RecoveryMainAsync(func() {
		internal.Setup()
	})
}
