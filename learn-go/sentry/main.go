package main

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
)

// https://developer.aliyun.com/article/951536
func main() {
	err := sentry.Init(sentry.ClientOptions{
		// 在此处设置您的 DSN 或设置 SENTRY_DSN 环境变量。
		Dsn: "https://3bc9042c5c6c32b8ac6a01209c1ceb57@o4510230836871168.ingest.us.sentry.io/4510230843883520",
		// 可以在这里设置 environment 和 release，
		// 也可以设置 SENTRY_ENVIRONMENT 和 SENTRY_RELEASE 环境变量。
		Environment: "dev",
		Release:     "v-1.0.0",
		// 允许打印 SDK 调试消息。
		// 入门或尝试解决某事时很有用。
		Debug:      true,
		EnableLogs: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	sentry.CaptureException(errors.New("erro -test"))
	sentry.CaptureMessage("Hello, sentry!")

	//sentry.Flush(5 * time.Second)

	hello1()
}

func hello1() {
	defer hello2()
	fmt.Println("hello1 world")
	hello3()
}

func hello2() {
	fmt.Println("hello2 world")
}

func hello3() {
	fmt.Println("hello3 world")
}
