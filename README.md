# uid

## 介绍

uid 是一个采用snowflake算法的uid生成库，借鉴了百度的uid-generte库和fastid库
+ 插件设计，业务也可自己注册IdGenerate接口实现
+ 默认实现采用lock-free

## 安装方式

```go
go get github.com/lanceryou/uid
```
## 开始使用

生成ID

```go
import (
  "fmt"
  "github.com/lanceryou/uid"
)

func ExampleNextID() {
  id := NextID()
  fmt.Printf("id generated: %v", ParseID(id))
}

func ExampleNextID1() {
	// 可注册IdGenerate实现，通过name查找调用
  id := IG("default").NextID()
  fmt.Printf("id generated: %v", ParseID(id))
}
```

## Benchmarks

### Benchmark Settings
+ 28 bits timestamp
+ 22 bits worker ID
+ 13 bits sequence number

```go
go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/lanceryou/uid
BenchmarkGenID-12        1218966              1043 ns/op
BenchmarkGenIDP-12        572725              2002 ns/op
PASS
ok      github.com/lanceryou/uid        3.995s

```