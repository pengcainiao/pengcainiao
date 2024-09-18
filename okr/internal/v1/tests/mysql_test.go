package tests

import (
	"fmt"
	"github.com/pengcainiao2/zero/core/env"
	"github.com/pengcainiao2/zero/tools/syncer"
	"testing"
)

type Okr struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// TestTry s
func TestTry(t *testing.T) {
	env.DbDSN = "penglonghui:Nrtg1X-syTXF@tcp(119.29.5.54:3306)/okr?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	//var (
	//	Ob   Okr
	//	asql = `select id,name from okr.okr where id = 1`
	//)
	//c := context.Background()
	if syncer.MySQL() == nil {
		fmt.Println("111")
	}
	fmt.Println(syncer.MySQL())
	//err := syncer.MySQL().Get(c, &Ob, asql)
	//if err != nil {
	//	logx.NewTraceLogger(c).Err(err).Msg("tests Mysql fail")
	//	return
	//}

}
