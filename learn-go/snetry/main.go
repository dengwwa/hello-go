package main

import (
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

// https://developer.aliyun.com/article/951536
func main() {
	err := sentry.Init(sentry.ClientOptions{
		// 在此处设置您的 DSN 或设置 SENTRY_DSN 环境变量。
		Dsn: "https://examplePublicKey@o0.ingest.sentry.io/0",
		// 可以在这里设置 environment 和 release，
		// 也可以设置 SENTRY_ENVIRONMENT 和 SENTRY_RELEASE 环境变量。
		Environment: "dev",
		Release:     "v-1.0.0",
		// 允许打印 SDK 调试消息。
		// 入门或尝试解决某事时很有用。
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// 在程序终止之前刷新缓冲事件。
	// 将超时设置为程序能够等待的最大持续时间。
	defer sentry.Flush(2 * time.Second)
}
