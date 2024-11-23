package timex

import (
	"fmt"
	"math/rand"
	"time"
)

// ReprOfDuration returns the string representation of given duration in ms.
func ReprOfDuration(duration time.Duration) string {
	return fmt.Sprintf("%.1fms", float32(duration)/float32(time.Millisecond))
}

//RandomExpireSeconds 在基准值的基础上加2小时随机时间
func RandomExpireSeconds(base time.Duration) time.Duration {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1800)

	return time.Second*time.Duration(n) + base
}
