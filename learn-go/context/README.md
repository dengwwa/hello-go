## 什么是 Context？
Context（上下文）是 Go 标准库 context 包中定义的一个接口，主要用于在多个 goroutine 之间传递请求范围的值、取消信号和超时信息。

## 为什么需要 Context？
* 在并发编程中，我们经常需要：
* 控制 goroutine 的生命周期
* 设置操作超时时间
* 在多个组件间传递共享数据
* 优雅地停止一系列相关操作

## Context 接口

``` go 
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
 ```
## 创建 Context
### 1. 基础 Context
```go
package main

import (
    "context"
    "fmt"
)

func main() {
    // 创建空的根 Context - 通常作为起点
    ctx := context.Background()
    fmt.Printf("Background: %T\n", ctx)
    
    // 当不确定用哪种 Context 时使用
    ctx2 := context.TODO()
    fmt.Printf("TODO: %T\n", ctx2)
}
```
### 2. 带值的 Context
```go
package main

import (
    "context"
    "fmt"
)

func main() {
    // 创建空的根 Context - 通常作为起点
    ctx := context.Background()
    fmt.Printf("Background: %T\n", ctx)
    
    // 当不确定用哪种 Context 时使用
    ctx2 := context.TODO()
    fmt.Printf("TODO: %T\n", ctx2)
}
```
### 3. 可取消的 Context
```go
package main

import (
	"context"
	"fmt"
	"time"
)
func main() {
    // 创建可取消的 Context
    ctx, cancel := context.WithCancel(context.Background())
    
    // 启动多个工作者
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }
    
    // 让 workers 运行一段时间
    time.Sleep(3 * time.Second)
    
    // 发送取消信号 - 所有监听这个 ctx 的 goroutine 都会收到
    fmt.Println("发送取消信号...")
    cancel()
    
    // 给 goroutines 时间退出
    time.Sleep(1 * time.Second)
    fmt.Println("主程序退出")
}

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: 收到取消信号，退出\n", id)
            return
        default:
            fmt.Printf("Worker %d: 工作中...\n", id)
            time.Sleep(1 * time.Second)
        }
    }
}
```
### 4. 带超时的 Context
```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
    // 创建 3 秒超时的 Context
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel() // 确保释放资源
    
    fmt.Println("开始执行耗时操作...")
    
    // 执行可能耗时的操作
    result, err := longRunningOperation(ctx)
    if err != nil {
        fmt.Printf("操作失败: %v\n", err)
        return
    }
    
    fmt.Printf("操作成功: %s\n", result)
}

func longRunningOperation(ctx context.Context) (string, error) {
    // 模拟一个需要 5 秒的操作
    select {
    case <-time.After(5 * time.Second):
        return "操作完成", nil
    case <-ctx.Done():
        return "", ctx.Err() // 返回超时或取消的错误
    }
}
```
### 5. 带截止时间的 Context
```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
    // 设置具体的截止时间（2秒后）
    deadline := time.Now().Add(2 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()
    
    go task(ctx)
    
    // 等待任务完成或超时
    <-ctx.Done()
    fmt.Printf("最终状态: %v\n", ctx.Err())
}

func task(ctx context.Context) {
    for i := 1; i <= 5; i++ {
        select {
        case <-ctx.Done():
            fmt.Printf("任务在步骤 %d 被中断: %v\n", i, ctx.Err())
            return
        default:
            fmt.Printf("完成步骤 %d\n", i)
            time.Sleep(500 * time.Millisecond)
        }
    }
    fmt.Println("所有步骤完成！")
}
```

## Context 的传递规则
### 1. 上下文继承
```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
    // 创建根 Context
    rootCtx := context.Background()
    
    // 添加跟踪 ID
    traceCtx := context.WithValue(rootCtx, "trace_id", "12345")
    
    // 设置超时
    timeoutCtx, cancel := context.WithTimeout(traceCtx, 5*time.Second)
    defer cancel()
    
    // 继续添加更多信息
    finalCtx := context.WithValue(timeoutCtx, "user_role", "admin")
    
    // 所有父 Context 的值都会传递下来
    processRequest(finalCtx)
}

func processRequest(ctx context.Context) {
    // 可以获取所有祖先 Context 设置的值
    traceID := ctx.Value("trace_id")
    userRole := ctx.Value("user_role")
    
    fmt.Printf("TraceID: %v, UserRole: %v\n", traceID, userRole)
    
    // 检查超时
    if deadline, ok := ctx.Deadline(); ok {
        fmt.Printf("截止时间: %v\n", deadline)
    }
}
```
### 2. 取消传播
```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    // 创建子 Context
    childCtx := context.WithValue(ctx, "level", "child")
    
    go parentProcess(ctx)
    go childProcess(childCtx)
    
    time.Sleep(2 * time.Second)
    
    // 取消父 Context，子 Context 也会被取消
    fmt.Println("取消父 Context...")
    cancel()
    
    time.Sleep(1 * time.Second)
}

func parentProcess(ctx context.Context) {
    <-ctx.Done()
    fmt.Println("父进程: 收到取消信号")
}

func childProcess(ctx context.Context) {
    <-ctx.Done()
    level := ctx.Value("level")
    fmt.Printf("子进程: 收到取消信号 (Level: %v)\n", level)
}
```

## 实际应用场景
### 场景 1：HTTP 请求处理
```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func main() {
    http.HandleFunc("/api/user", userHandler)
    fmt.Println("服务器启动在 :8080")
    http.ListenAndServe(":8080", nil)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    // 使用请求的 Context
    ctx := r.Context()
    
    // 为数据库查询设置超时
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // 模拟数据库查询
    user, err := getUserFromDB(ctx, r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusRequestTimeout)
        return
    }
    
    fmt.Fprintf(w, "用户信息: %s", user)
}

func getUserFromDB(ctx context.Context, userID string) (string, error) {
    // 模拟数据库查询耗时
    select {
    case <-time.After(3 * time.Second): // 模拟慢查询
        return fmt.Sprintf("用户 %s 的数据", userID), nil
    case <-ctx.Done():
        return "", fmt.Errorf("查询超时: %v", ctx.Err())
    }
}
```
### 场景 2：并发任务控制
```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 启动多个并发任务
    results := make(chan string, 3)
    
    go fetchData(ctx, "API-1", results)
    go fetchData(ctx, "API-2", results)
    go fetchData(ctx, "API-3", results)
    
    // 收集结果
    for i := 0; i < 3; i++ {
        select {
        case result := <-results:
            fmt.Println("收到结果:", result)
        case <-ctx.Done():
            fmt.Println("操作被取消")
            return
        }
    }
}

func fetchData(ctx context.Context, source string, results chan<- string) {
    // 模拟不同响应时间的任务
    var delay time.Duration
    switch source {
    case "API-1":
        delay = 1 * time.Second
    case "API-2":
        delay = 2 * time.Second
    case "API-3":
        delay = 3 * time.Second
    }
    
    select {
    case <-time.After(delay):
        results <- fmt.Sprintf("来自 %s 的数据", source)
    case <-ctx.Done():
        fmt.Printf("%s: 任务取消\n", source)
    }
}
```
## 最佳实践
### 1. 函数签名规范
```go
package main
import "context"

// 好的做法：Context 作为第一个参数
func ProcessUser(ctx context.Context, userID string) error {
    // 在开始处检查 Context 是否已取消
    if err := ctx.Err(); err != nil {
        return err
    }
    
    // 业务逻辑...
    return nil
}

// 不好的做法：Context 放在后面或不使用
func ProcessUserBad(userID string, ctx context.Context) error {
    // 不符合惯例
    return nil
}
```
### 2. 正确使用 WithValue
```go
package main
import "context"
import "fmt"
// 定义类型安全的 key
type contextKey string

const (
    userIDKey contextKey = "userID"
    authKey   contextKey = "authToken"
)

func main() {
    ctx := context.Background()
    
    // 使用类型安全的 key
    ctx = context.WithValue(ctx, userIDKey, 123)
    ctx = context.WithValue(ctx, authKey, "token123")
    
    // 获取值时进行类型断言
    if userID, ok := ctx.Value(userIDKey).(int); ok {
        fmt.Printf("UserID: %d\n", userID)
    }
}
```
### 3. 资源清理
```go
package main
import "context"
import "time"
import "fmt"

func properResourceCleanup() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // 确保 cancel 被调用，释放资源
    
    // 即使函数提前返回，defer 也会确保 cancel 被调用
    result, err := someOperation(ctx)
    if err != nil {
        return
    }
    
    fmt.Println("结果:", result)
}
```

## 常见错误和陷阱
1. 忘记取消 Context：如果 Context 没有被取消，可能会导致资源泄漏。
   ```go
   func memoryLeakExample() {
    // 错误：没有调用 cancel，可能导致内存泄漏
    ctx, cancel := context.WithCancel(context.Background())
    _ = ctx // 使用 ctx
    // 忘记调用 cancel()!
    
    // 正确做法：使用 defer
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()}
   ```
2. 在结构体中存储 Context
    ```go
    // 不好的做法
    type Service struct {
    ctx context.Context // 不应该存储 Context
    }
    
    // 好的做法
    type Service struct {
    // 不存储 Context，在方法中传递
    }
    
    func (s *Service) Process(ctx context.Context, data string) error {
    // 使用传入的 Context
    return nil
    }
   ```
不要存储 Context 在结构体中的主要原因：
* **生命周期问题** - 存储的 Context 可能过期，导致后续所有操作失败
* **灵活性丧失** - 无法为不同的操作设置不同的超时和取消策略
* **测试困难** - 难以模拟不同的 Context 场景
* **语义混淆** - 不清楚存储的 Context 代表什么含义
* **并发安全问题** - 多个 goroutine 可能同时使用同一个存储的 Context

**正确做法**：
* 将 Context 作为函数的第一个参数传递
* 每个请求/操作使用独立的 Context
* 在结构体中只存储真正的配置数据，而不是 Context

**例外情况**：
* 对象生命周期与 Context 完全绑定时
    ```go
    // 可以接受的情况：对象的生命周期与 Context 完全绑定
    type Worker struct {
        ctx    context.Context
        cancel context.CancelFunc
        done   chan struct{}
    }
    
    func NewWorker() *Worker {
        ctx, cancel := context.WithCancel(context.Background())
        w := &Worker{
            ctx:    ctx,
            cancel: cancel,
            done:   make(chan struct{}),
        }
        go w.run()
        return w
    }
    
    func (w *Worker) run() {
        defer close(w.done)
        
        for {
            select {
            case <-w.ctx.Done():
                return // Worker 停止
            case <-time.After(time.Second):
                // 执行工作
            }
        }
    }
    
    func (w *Worker) Stop() {
        w.cancel()
        <-w.done
    }
    ```
* 仅存储 base Context 用于创建派生 Context 时
  ```go
  type RequestProcessor struct {
    baseCtx context.Context
  }
    
  func NewRequestProcessor(baseCtx context.Context) *RequestProcessor {
  // 存储 baseCtx 用于创建派生 Context 是可以的
  return &RequestProcessor{baseCtx: baseCtx}
  }
    
  func (p *RequestProcessor) ProcessRequest(req *Request) error {
  // 为每个请求创建新的 Context，基于存储的 baseCtx
  ctx, cancel := context.WithTimeout(p.baseCtx, req.Timeout)
  defer cancel()
    
      return p.process(ctx, req)
  }
  ```

## 总结
Context 是 Go 并发编程的核心组件，主要用途：
1. 传播取消信号 - 优雅停止相关操作
2. 设置超时和截止时间 - 防止操作无限期等待
3. 传递请求范围数据 - 在调用链中共享数据
4. 控制 goroutine 生命周期 - 管理并发执行

记住这些关键点：
* Context 应该是不可变的
* 总是作为第一个参数传递
* 及时调用 cancel 函数
* 不要存储 Context 在结构体中
* 使用类型安全的 key 来存储值