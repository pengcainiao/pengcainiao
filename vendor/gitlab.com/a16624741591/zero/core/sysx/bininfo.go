package sysx

import "fmt"

var (
	BuildDate      string //编译时间
	BuildGoVersion string //编译使用的Go版本
	GitCommitLog   string //git 日志
	AppVersion     string //版本号
	SubSystem      string
)

func init() {
	//if BuildDate != "" {
	fmt.Printf("\t\t\t------------ 编译信息 ------------ \n BuildDate:%30s \n BuildGoVersion:%39s \n GitCommitLog:%50s\n AppVersion:%19s\n----------------------------------------------------\n", BuildDate, BuildGoVersion, GitCommitLog, AppVersion)
	//}
}
