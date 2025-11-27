package main

import (
	"fmt"
	"time"
)

func main() {
	const AppDeletedRetentionDuration = 7 * 24 * time.Hour // 删除应用保留时长，保留期间所有底层资源继续正常运行
	milliseconds := GetNowUTCMilli() - AppDeletedRetentionDuration.Milliseconds()
	fmt.Println(milliseconds)
}
func GetNowUTCMilli() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}
