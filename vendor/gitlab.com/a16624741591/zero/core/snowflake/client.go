package sonyflake

import (
	"fmt"
	"strconv"
	"time"
)

var (
	sf *Sonyflake
)

func init() {
	var st Settings
	st.StartTime = time.Date(2020, 5, 6, 0, 0, 0, 0, time.Local)
	sf = NewSonyflake(st)
}

//GenerateID 创建分布式ID
func GenerateID() string {
	id, err := sf.NextID()
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return strconv.FormatUint(id, 10)
}

func GenerateInt64ID() uint64 {
	id, _ := sf.NextID()
	return id
}
