package sysx

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func CallerFunctionName(skipCaller int) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	pc, file, line, _ := runtime.Caller(skipCaller)
	f := runtime.FuncForPC(pc)
	var funcName = f.Name()[strings.LastIndex(f.Name(), "/"):]
	return path.Base(file) + funcName + "：L" + strconv.Itoa(line)
}
