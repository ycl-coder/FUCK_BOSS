# Technology Stack

## Project Type

全国公司曝光平台是一个基于 gRPC 的 Web 应用，采用前后端分离架构。后端使用 Go 语言提供 gRPC 服务，前端使用 React 构建用户界面。

## Core Technologies

### Primary Language(s)

- **后端语言**: Go 1.21+
- **前端语言**: TypeScript 5.0+ (React)
- **协议定义**: Protocol Buffers (protobuf)
- **包管理**: 
  - Go: `go mod`
  - Frontend: `npm` / `yarn` / `pnpm`

### Key Dependencies/Libraries

#### 后端 (Go)
- **gRPC-Go**: `google.golang.org/grpc` - gRPC 框架核心
- **Protocol Buffers**: `google.golang.org/protobuf` - 数据序列化
- **PostgreSQL 驱动**: `github.com/lib/pq` 或 `gorm.io/gorm` - 数据库访问
- **Redis 客户端**: `github.com/redis/go-redis/v9` - 缓存和会话管理
- **依赖注入**: `github.com/google/wire` 或 `github.com/uber-go/fx` - DDD 依赖管理
- **配置管理**: `github.com/spf13/viper` - 配置读取
- **日志**: `go.uber.org/zap` - 结构化日志
- **验证**: `github.com/go-playground/validator/v10` - 数据验证
- **HTTP 网关**: `github.com/grpc-ecosystem/grpc-gateway/v2` - gRPC 转 HTTP/REST（可选）

#### 前端 (React)
- **React**: `^18.2.0` - UI 框架
- **TypeScript**: `^5.0.0` - 类型安全
- **路由**: `react-router-dom` - 前端路由
- **状态管理**: `zustand` 或 `@reduxjs/toolkit` - 状态管理
- **HTTP 客户端**: `axios` 或 `@tanstack/react-query` - API 调用
- **gRPC Web**: `@grpc/grpc-js` + `@grpc-web/protoc-gen-grpc-web` - gRPC Web 客户端
- **UI 组件库**: `antd` 或 `@mui/material` - 组件库
- **构建工具**: `vite` 或 `create-react-app` - 构建和开发工具

### Application Architecture

#### 架构模式：领域驱动设计 (DDD)

采用 DDD 分层架构，将业务逻辑与基础设施分离：

```
┌─────────────────────────────────────┐
│   Presentation Layer (API/UI)       │
│   - gRPC Handlers                   │
│   - HTTP Gateway (可选)             │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Application Layer                  │
│   - Use Cases / Application Services │
│   - DTOs / Commands / Queries       │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Domain Layer (核心业务逻辑)        │
│   - Entities (聚合根)               │
│   - Value Objects                   │
│   - Domain Services                 │
│   - Repository Interfaces            │
│   - Domain Events                   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Infrastructure Layer               │
│   - Repository Implementations       │
│   - Database (PostgreSQL)           │
│   - Cache (Redis)                   │
│   - External Services                │
└─────────────────────────────────────┘
```

#### 领域划分（Bounded Contexts）

1. **内容领域 (Content Context)**
   - 聚合根：`Post` (曝光内容)
   - 值对象：`City`, `CompanyName`, `Content`
   - 领域服务：内容发布、内容审核

2. **搜索领域 (Search Context)**
   - 聚合根：`SearchIndex`
   - 领域服务：全文搜索、城市筛选

3. **用户领域 (User Context)** (未来版本)
   - 聚合根：`AnonymousUser`
   - 领域服务：匿名标识管理

### Data Storage

#### Primary Storage: PostgreSQL

- **版本**: PostgreSQL 14+
- **用途**: 
  - 持久化存储曝光内容
  - 存储城市、公司等基础数据
  - 支持全文搜索（PostgreSQL Full-Text Search）
- **连接池**: 使用 `pgx` 或 `gorm` 管理连接池
- **迁移工具**: `golang-migrate/migrate` 或 `gorm` 自动迁移

#### Caching: Redis

- **版本**: Redis 7.0+
- **用途**:
  - 缓存热门内容列表
  - 缓存城市列表、公司列表
  - 搜索结果的短期缓存
  - 限流和防刷机制
  - 会话存储（未来版本）
- **数据结构**:
  - `String`: 简单缓存
  - `Hash`: 对象缓存
  - `Sorted Set`: 热门内容排行
  - `Set`: 去重、标签

#### Data Formats

- **API 通信**: Protocol Buffers (protobuf)
- **配置**: YAML / JSON
- **日志**: JSON (结构化日志)

### External Integrations

#### APIs
- **内容审核服务** (未来版本): 第三方内容审核 API
- **地理位置服务** (可选): IP 定位服务

#### Protocols
- **gRPC**: 后端服务间通信
- **gRPC-Web**: 前端与后端通信
- **HTTP/REST**: 可选，通过 grpc-gateway 提供 RESTful API

#### Authentication
- **匿名发布**: 无需认证（第一版本）
- **未来版本**: JWT Token 或 Session 认证

## Development Environment

### Build & Development Tools

#### 后端
- **Build System**: `go build` / `make`
- **开发工具**: 
  - `air` - 热重载
  - `golangci-lint` - 代码检查
- **API 生成**: `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc`
- **代码生成**: `wire` (依赖注入代码生成)

#### 前端
- **Build System**: `vite` 或 `webpack`
- **开发工具**: 
  - `vite` HMR (热模块替换)
  - `eslint` - 代码检查
  - `prettier` - 代码格式化

### Code Quality Tools

#### 后端
- **Static Analysis**: `golangci-lint` (集成多种 linter)
- **Formatting**: `gofmt` / `goimports`
- **Testing Framework**: 
  - `testing` (标准库) - 单元测试
  - `testify` - 测试断言和 mock
  - `gomock` - 接口 mock 生成
- **Coverage**: `go test -cover`
- **Documentation**: `godoc` / `swagger` (如果提供 REST API)

## Go 代码规范

为形成统一的 Go 编码风格，以保障项目代码的易维护性和编码安全性，特制定本规范。

每项规范内容，给出了要求等级，其定义为：

- **[火箭] 必须（Mandatory）**：用户必须采用
- **[火] 推荐（Preferable）**：用户理应采用，但如有特殊情况，可以不采用
- **[灯泡] 可选（Optional）**：用户可参考，自行决定是否采用
- **[旗子] 建议（Suggestion）**：最佳实践建议

### 代码风格

#### 2.1 【必须】格式化

- **[火箭]** 代码都必须用 `gofmt` 格式化

#### 2.2 【推荐】换行

- **[火]** 建议一行代码不要超过 120 列，超过的情况，使用合理的换行方法换行

#### 2.3 【必须】括号和空格

- **[火箭]** 遵循 `gofmt` 的逻辑
- **[火]** 运算符和操作数之间要留空格
- **[灯泡]** 作为输入参数或者数组下标时，运算符和运算数之间不需要空格，紧凑展示

#### 2.4 【必须】import 规范

- **[火箭]** 使用 `goimports` 自动格式化引入的包名
- **[火]** 引入单个包，也使用括号包裹
  ```go
  // 应该采用如下格式：
  import ("fmt")
  // 而不是这样：
  import "fmt"
  ```
- **[火箭]** 如果你引入了多种类型的包，必须对包进行分组管理，将包分为标准库包、程序内部包、第三方包，并将标准库作为第一组，三组包用空行间隔
- **[火]** `goimports` 或者 `gofmt` 会自动把依赖包按首字母排序
- **[灯泡]** 匿名包的引用必须使用一个新的分组引入
- **[旗子]** 不要使用相对路径引入包
  ```go
  // 不要采用这种方式
  import (
    "../net"
  )
  // 应该使用完整的路径引入包：
  import (
    "xxxx.com/proj/net"
  )
  ```
- **[火箭]** 包名和 git 路径名不一致时，使用别名代替
  ```go
  import (
    opentracing "github.com/opentracing/opentracing-go"
  )
  ```
- **[火箭]** 【推荐】在匿名引入的每个包上推荐写上注释说明
- **[火]** 应该采用如下方式进行组织你的包：
  ```go
  import (
    // standard package
    "encoding/json"
    "strings"
  
    // inner package
    "myproject/models"
    "myproject/controller"
  
    // third-party package
    "git.obc.im/obc/utils"
    "git.obc.im/dep/beego"
    "git.obc.im/dep/mysql"
  
    // alias package
    opentracing "github.com/opentracing/opentracing-go"
  
    // anonymous import package
    // import filesystem storage driver
    _ "github.com/org/repo/pkg/storage/filesystem"
  )
  ```

#### 2.5 【必须】错误处理

##### 2.5.1 【必须】error 处理

- **[火箭]** error 作为函数的值返回，必须对 error 进行处理，或将返回值赋值给明确忽略
- **[火]** error 作为函数的值返回且有多个返回值的时候，error 必须是最后一个参数
  ```go
  // 不要采用这种方式
  func do() (error, int) {
  }
  // 要采用下面的方式
  func do() (int, error) {
  }
  ```
- **[火箭]** 错误描述不需要标点结尾
- **[火]** 采用独立的错误流进行处理
  ```go
  // 不要采用这种方式
  if err != nil {
    // error handling
  } else {
    // normal code
  }
  
  // 而要采用下面的方式
  if err != nil {
    // error handling
    return // or continue, etc.
  }
  // normal code
  ```
- 如果返回值需要初始化，则采用下面的方式：
  ```go
  x, err := f()
  if err != nil {
    // error handling
    return // or continue, etc.
  }
  // use x
  ```
- 错误返回的判断独立处理，不与其他变量组合逻辑判断
  ```go
  // 不要采用这种方式：
  x, y, err := f()
  if err != nil || y == nil {
    return err   // 当y与err都为空时，函数的调用者会出现错误的调用逻辑
  }
  
  // 应当使用如下方式：
  x, y, err := f()
  if err != nil {
    return err
  }
  if y == nil {
    return fmt.Errorf("some error")
  }
  ```
- **[火]** 【推荐】建议 go1.13 以上，error 生成方式为：`fmt.Errorf("module xxx: %w", err)`

##### 2.5.2 【必须】panic 处理

- **[火箭]** 在业务逻辑处理中禁止使用 panic
- **[火]** 在 main 包中只有当完全不可运行的情况可使用 panic，例如：文件无法打开，数据库无法连接导致程序无法正常运行
- **[灯泡]** 对于其它的包，可导出的接口不能有 panic，只能在包内使用
- **[旗子]** 建议在 main 包中使用 `log.Fatal` 来记录错误，这样就可以由 log 来结束程序，或者将 panic 抛出的异常记录到日志文件中，方便排查问题
- **[火箭]** panic 捕获只能到 goroutine 最顶层，每个自行启动的 goroutine，必须在入口处捕获 panic，并打印详细堆栈信息或进行其它处理

##### 2.5.3 【必须】recover 处理

- **[火箭]** recover 用于捕获 runtime 的异常，禁止滥用 recover
- **[火]** 必须在 defer 中使用，一般用来捕获程序运行期间发生异常抛出的 panic 或程序主动抛出的 panic
  ```go
  package main
  import (
    "log"
  )
  func main() {
    defer func() {
        if err := recover(); err != nil {
            // do something or record log
            log.Println("exec panic error: ", err)
            // log.Println(debug.Stack())
        }
    }()
    
    getOne()
    
    panic(11) //手动抛出panic
  }
  // getOne 模拟slice越界 runtime运行时抛出的panic
  func getOne() {
    defer func() {
        if err := recover(); err != nil {
            // do something or record log
            log.Println("exec panic error: ", err)
            // log.Println(debug.Stack())
        }
    }()
    
    var arr = []string{"a", "b", "c"}
    log.Println("hello,", arr[4])
  }
  ```

#### 2.6 【必须】单元测试

- **[火箭]** 单元测试文件名命名规范为 `example_test.go`
- **[火]** 测试用例的函数名称必须以 `Test` 开头，例如 `TestExample`
- **[灯泡]** 每个重要的可导出函数都要首先编写测试用例，测试用例和正规代码一起提交方便进行回归测试

#### 2.7 【必须】类型断言失败处理

- type assertion 的单个返回值形式针对不正确的类型将产生 panic。因此，请始终使用"comma ok"的惯用法
  ```go
  // 不要采用这种方式
  t := i.(string)
  // 而要采用下面的方式
  t, ok := i.(string)
  if !ok {
    // 优雅地处理错误
  }
  ```

### 注释

1. 在编码阶段同步写好变量、函数、包注释，注释可以通过 `godoc` 导出生成文档
2. 注释必须是完整的句子，以需要注释的内容作为开头，句点作为结尾
3. 程序中每一个被导出的(大写的)名字，都应该有一个文档注释

#### 3.1 【必须】包注释

- **[火箭]** 每个包都应该有一个包注释
- **[火]** 包如果有多个 go 文件，只需要出现在一个 go 文件中（一般是和包同名的文件）即可，格式为："// Package 包名 包信息描述"
  ```go
  // Package math provides basic constants and mathematical functions.
  package math
  // 或者
  /*
  Package template implements data-driven templates for generating textual
  output such as HTML.
  ....
  */
  package template
  ```

#### 3.2 【必须】结构体注释

- **[火箭]** 每个需要导出的自定义结构体或者接口都必须有注释说明
- **[火]** 注释对结构进行简要介绍，放在结构体定义的前一行
- **[灯泡]** 格式为："// 结构体名 结构体信息描述"
- **[旗子]** 结构体内的每个需要导出的成员变量都要有说明，该说明放在成员变量的前一行
  ```go
  // User 用户结构定义了用户基础信息
  type User struct {
    // UserName 用户名
    UserName string
    // Email 邮箱
    Email string
  }
  ```

#### 3.3 【必须】方法注释

- **[火箭]** 每个需要导出的函数或者方法（结构体或者接口下的函数称为方法）都必须有注释
- **[火]** 注释描述函数或方法功能、调用方等信息
- **[灯泡]** 格式为："// 函数名 函数信息描述"
  ```go
  // NewAttrModel 是属性数据层操作类的工厂方法
  func NewAttrModel(ctx common.Context) AttrModel {
    // TODO
  }
  ```

#### 3.4 【必须】变量注释

- **[火箭]** 每个需要导出的常量和变量都必须有注释说明
- **[火]** 该注释对常量或变量进行简要介绍，放在常量或者变量定义的前一行
- **[灯泡]** 格式为："// 变量名 变量信息描述"
  ```go
  // FlagConfigFile 配置文件的命令行参数名
  const FlagConfigFile = "--config"
  // FullName 返回指定用户名的完整名称
  var FullName = func(username string) string {
    return fmt.Sprintf("fake-%s", username)
  }
  ```

#### 3.5 【必须】类型注释

- **[火箭]** 每个需要导出的类型定义(type definition)和类型别名(type aliases)都必须有注释说明
- **[火]** 该注释对类型进行简要介绍，放在定义的前一行
- **[灯泡]** 格式为："// 类型名 类型信息描述"
  ```go
  // StorageClass 存储类型
  type StorageClass string
  // FakeTime 标准库时间的类型别名
  type FakeTime = time.Time
  ```

### 命名规范

命名是代码规范中很重要的一部分，统一的命名规范有利于提高代码的可读性，好的命名仅仅通过命名就可以获取到足够多的信息。

#### 4.1 【推荐】包命名

- **[火箭]** 保持 package 的名字和目录一致
- **[火]** 尽量采取有意义、简短的包名，尽量不要和标准库冲突
- **[灯泡]** 包名应该为小写单词，不要使用下划线或者混合大小写，使用多级目录来划分层级
- **[旗子]** 项目名可以通过中划线来连接多个单词
- **[火箭]** 简单明了的包命名，如：time、list、http
- **[火]** 不要使用无意义的包名，如：util、common、misc

#### 4.2 【必须】文件命名

- **[火箭]** 采用有意义，简短的文件名
- **[火]** 文件名应该采用小写，并且使用下划线分割各个单词

#### 4.3 【必须】结构体命名

- **[火箭]** 采用驼峰命名方式，首字母根据访问控制采用大写或者小写
- **[火]** 结构体名应该是名词或名词短语，如 `Customer`、`WikiPage`、`Account`、`AddressParser`
- **[灯泡]** 避免使用 `Manager`、`Processor`、`Data`、`Info` 这样的结构体名，结构体名不应当是动词
- **[旗子]** 结构体的申明和初始化格式采用多行，例如：
  ```go
  // User 多行申明
  type User struct {
    // Username 用户名
    UserName string
    // Email 电子邮件地址
    Email string
  }
  // 多行初始化
  u := User{
    UserName: "john",
    Email:    "john@example.com",
  }
  ```

#### 4.4 【推荐】接口命名

- **[火箭]** 命名规则基本保持和结构体命名规则一致
- **[火]** 单个函数的结构名以 "er" 作为后缀，例如 `Reader`，`Writer`
  ```go
  // Reader 字节数组读取接口
  type Reader interface {
    // Read 读取整个给定的字节数据并返回读取的长度
    Read(p []byte) (n int, err error)
  }
  ```
- **[火箭]** 两个函数的接口名综合两个函数名
- **[火]** 三个以上函数的接口名，类似于结构体名
  ```go
  // Car 小汽车结构申明
  type Car interface {
    // Start ...
    Start([]byte)
    // Stop ...
    Stop() error
    // Recover ...
    Recover()
  }
  ```

#### 4.5 【必须】变量命名

- **[火箭]** 变量名必须遵循驼峰式，首字母根据访问控制决定使用大写或小写
- **[火]** 特有名词时，需要遵循以下规则：
  - 如果变量为私有，且特有名词为首个单词，则使用小写，如 `apiClient`
  - 其他情况都应该使用该名词原有的写法，如 `APIClient`、`repoID`、`UserID`
  - 错误示例：`UrlArray`，应该写成 `urlArray` 或者 `URLArray`
- **[火箭]** 若变量类型为 bool 类型，则名称应以 `Has`，`Is`，`Can` 或者 `Allow` 开头
- **[火]** 私有全局变量和局部变量规范一致，均以小写字母开头

#### 4.6 【必须】常量命名

- 常量均需遵循驼峰式
  ```go
  // AppVersion 应用程序版本号定义
  const AppVersion = "1.0.0"
  ```
- 如果是枚举类型的常量，需要先创建相应类型：
  ```go
  // Scheme 传输协议
  type Scheme string
  const (
    // HTTP 表示HTTP明文传输协议
    HTTP Scheme = "http"
    // HTTPS 表示HTTPS加密传输协议
    HTTPS Scheme = "https"
  )
  ```
- 私有全局常量和局部变量规范一致，均以小写字母开头
  ```go
  const appVersion = "1.0.0"
  ```

### 控制结构

#### 5.1 【推荐】if

- if 接受初始化语句，约定如下方式建立局部变量：
  ```go
  if err := file.Chmod(0664); err != nil {
    return err
  }
  ```

#### 5.2 【推荐】for

- 采用短声明建立局部变量：
  ```go
  sum := 0
  for i := 0; i < 10; i++ {
    sum += 1
  }
  ```

#### 5.3 【必须】range

- 如果只需要第一项（key），就丢弃第二个：
  ```go
  for key := range m {
    if key.expired() {
      delete(m, key)
    }
  }
  ```
- 如果只需要第二项，则把第一项置为下划线：
  ```go
  sum := 0
  for _, value := range array {
    sum += value
  }
  ```

#### 5.4 【必须】switch

- 要求必须有 default：
  ```go
  switch os := runtime.GOOS; os {
  case "darwin":
    fmt.Println("OS X.")
  case "linux":
    fmt.Println("Linux.")
  default:
    // freebsd, openbsd,
    // plan9, windows...
    fmt.Printf("%s.\n", os)
  }
  ```

#### 5.5 【推荐】return

- 尽早 return，一旦有错误发生，马上返回：
  ```go
  f, err := os.Open(name)
  if err != nil {
    return err
  }
  d, err := f.Stat()
  if err != nil {
    f.Close()
    return err
  }
  codeUsing(f, d)
  ```

#### 5.6 【必须】goto

- 业务代码禁止使用 goto，其他框架或底层源码推荐尽量不用

### 函数

#### 6.1 【推荐】函数参数

- 函数返回相同类型的两个或三个参数，或者如果从上下文中不清楚结果的含义，使用命名返回，其它情况不建议使用命名返回
  ```go
  // Parent1 ...
  func (n Node) Parent1() Node
  // Parent2 ...
  func (n Node) Parent2() (Node, error)
  // Location ...
  func (f *Foo) Location() (lat, long float64, err error)
  ```
- **[火箭]** 传入变量和返回变量以小写字母开头
- **[火]** 参数数量均不能超过 5 个
- **[灯泡]** 尽量用值传递，非指针传递
- **[旗子]** 传入参数是 map，slice，chan，interface 不要传递指针

#### 6.2 【必须】defer

- **[火箭]** 当存在资源管理时，应紧跟 defer 函数进行资源的释放
- **[火]** 判断是否有错误发生之后，再 defer 释放资源
  ```go
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  // 如果操作成功，再defer Close()
  defer resp.Body.Close()
  ```
- 禁止在循环中使用 defer，举例如下：
  ```go
  // 不要这样使用
  func filterSomething(values []string) {
    for _, v := range values {
      fields, err := db.Query(v) // 示例，实际不要这么查询，防止sql注入
      if err != nil {
        // xxx
      }
      defer fields.Close()
      // 继续使用fields
    }
  }
  // 应当使用如下的方式：
  func filterSomething(values []string) {
    for _, v := range values {
      func() {
        fields, err := db.Query(v) // 示例，实际不要这么查询，防止sql注入
        if err != nil {
          // ...
        }
        defer fields.Close()
        // 继续使用fields
      }()
    }
  }
  ```

#### 6.3 【必须】方法的接收器

- **[火箭]** 接收器的命名在函数超过 20 行的时候不要用单字符
- **[火]** 命名不能采用 me，this，self 这类易混淆名称

#### 6.4 【推荐】代码行数

- **[火箭]** 【必须】文件长度不能超过 800 行
- **[火]** 【推荐】函数长度不能超过 80 行

#### 6.5 【必须】嵌套

- 嵌套深度不能超过 4 层：
  ```go
  // AddArea 添加成功或出错
  func (s *BookingService) AddArea(areas ...string) error {
    s.Lock()
    defer s.Unlock()
    
    for _, area := range areas {
      for _, has := range s.areas {
        if area == has {
          return srverr.ErrAreaConflict
        }
      }
      s.areas = append(s.areas, area)
      s.areaOrders[area] = new(order.AreaOrder)
    }
    return nil
  }
  
  // 建议调整为这样：
  // AddArea 添加成功或出错
  func (s *BookingService) AddArea(areas ...string) error {
    s.Lock()
    defer s.Unlock()
    
    for _, area := range areas {
      if s.HasArea(area) {
        return srverr.ErrAreaConflict
      }
      s.areas = append(s.areas, area)
      s.areaOrders[area] = new(order.AreaOrder)
    }
    return nil
  }
  
  // HasArea ...
  func (s *BookingService) HasArea(area string) bool {
    for _, has := range s.areas {
      if area == has {
        return true
      }
    }
    return false
  }
  ```

#### 6.6 【推荐】变量声明

- 变量声明尽量放在变量第一次使用前面，就近原则

#### 6.7 【必须】魔法数字

- 如果魔法数字出现超过 2 次，则禁止使用
  ```go
  func getArea(r float64) float64 {
    return 3.14 * r * r
  }
  
  func getLength(r float64) float64 {
    return 3.14 * 2 * r
  }
  
  // PI ...
  const PI = 3.14
  func getArea(r float64) float64 {
    return PI * r * r
  }
  
  func getLength(r float64) float64 {
    return PI * 2 * r
  }
  ```

### 依赖管理

#### 7.1 【必须】go modules

- **[火箭]** go1.11 以上必须使用 go modules 模式：
  ```bash
  go mod init github.com/group/myrepo
  ```

#### 7.2 【推荐】代码提交

- **[火箭]** 建议所有不对外开源的工程的 module name 使用 `github.com/group/repo`，方便他人直接引用
- **[火]** 建议使用 go modules 作为依赖管理的项目不提交 vendor 目录
- **[灯泡]** 建议使用 go modules 管理依赖的项目将 go.sum 文件不添加到忽略提交规则中

### 应用服务

#### 8.1 【推荐】应用服务接口建议有 readme.md

- 其中建议包括服务基本描述、使用方法、部署时的限制与要求、基础环境依赖（例如最低 go 版本、最低外部通用包版本）等

#### 8.2 【必须】应用服务必须要有接口测试

### 常用工具

go 语言本身在代码规范性这方面也做了很多努力，很多限制都是强制语法要求，例如左大括号不换行，引用的包或者定义的变量不使用会报错，此外 go 还是提供了很多好用的工具帮助我们进行代码的规范。

- **[火箭]** `gofmt`，大部分的格式问题可以通过 `gofmt` 解决，`gofmt` 自动格式化代码，保证所有的 go 代码与官方推荐的格式保持一致，于是所有格式有关问题，都以 `gofmt` 的结果为准
- **[火]** `goimports`，此工具在 `gofmt` 的基础上增加了自动删除和引入包
- **[灯泡]** `go vet`，vet 工具可以帮我们静态分析我们的源码存在的各种问题，例如多余的代码，提前 return 的逻辑，struct 的 tag 是否符合标准等。编译前先执行代码静态分析
- **[旗子]** `golint`，类似 javascript 中的 `jslint` 的工具，主要功能就是检测代码中不规范的地方

#### 前端
- **Static Analysis**: `eslint` + `typescript-eslint`
- **Formatting**: `prettier`
- **Testing Framework**:
  - `vitest` / `jest` - 单元测试
  - `@testing-library/react` - React 组件测试
  - `playwright` / `cypress` - E2E 测试
- **Type Checking**: TypeScript compiler

### Version Control & Collaboration

- **VCS**: Git
- **Branching Strategy**: Git Flow
  - `main`: 生产环境
  - `develop`: 开发分支
  - `feature/*`: 功能分支
  - `fix/*`: 修复分支
- **Code Review Process**: 
  - 所有代码必须经过 Code Review
  - 至少一人 Approve 才能合并
  - CI/CD 检查通过后才能合并

## Deployment & Distribution

### Target Platform(s)
- **后端**: Linux (Docker 容器)
- **前端**: 静态文件托管 (Nginx / CDN)
- **数据库**: PostgreSQL (云服务或自托管)
- **缓存**: Redis (云服务或自托管)

### Distribution Method
- **容器化**: Docker + Docker Compose (开发环境)
- **Kubernetes**: 生产环境 (可选)
- **CI/CD**: GitHub Actions / GitLab CI

### Installation Requirements
- **开发环境**:
  - Go 1.21+
  - Node.js 18+
  - PostgreSQL 14+
  - Redis 7.0+
  - Docker & Docker Compose (推荐)

### Update Mechanism
- **数据库迁移**: 通过 migration 工具自动执行
- **服务更新**: 滚动更新 (零停机)

## Technical Requirements & Constraints

### Performance Requirements

- **API 响应时间**: 
  - 列表查询: < 200ms (有缓存)
  - 详情查询: < 100ms (有缓存)
  - 发布操作: < 500ms
- **并发支持**: 支持 1000+ 并发请求
- **数据库连接池**: 最大 100 连接
- **Redis 连接池**: 最大 50 连接

### Compatibility Requirements

- **Go 版本**: >= 1.21
- **PostgreSQL 版本**: >= 14.0
- **Redis 版本**: >= 7.0
- **Node.js 版本**: >= 18.0
- **浏览器支持**: Chrome, Firefox, Safari, Edge (最新 2 个版本)

### Security & Compliance

#### Security Requirements
- **数据加密**: 
  - 传输层: TLS 1.3
  - 存储层: 敏感数据加密存储
- **SQL 注入防护**: 使用参数化查询，禁止拼接 SQL
- **XSS 防护**: 前端输入输出转义
- **CSRF 防护**: Token 验证 (未来版本)
- **限流**: Redis 实现接口限流，防止恶意刷接口
- **匿名保护**: 
  - 不记录用户 IP (或加密存储)
  - 不存储任何可追踪信息

#### Compliance Standards
- **数据隐私**: 遵循 GDPR 原则（最小化数据收集）
- **内容合规**: 建立内容审核机制，防止违法内容

#### Threat Model
- **主要威胁**:
  - 恶意刷接口
  - 虚假/恶意内容发布
  - 数据泄露
  - DDoS 攻击
- **防护措施**:
  - 接口限流
  - 内容审核（人工+自动）
  - 数据加密
  - CDN + WAF

### Scalability & Reliability

#### Expected Load
- **初期**: 1000 DAU, 10000 条内容
- **中期**: 10000 DAU, 100000 条内容
- **长期**: 100000+ DAU, 1000000+ 条内容

#### Availability Requirements
- **目标可用性**: 99.9% (月度)
- **故障恢复时间**: < 30 分钟
- **数据备份**: 每日自动备份，保留 30 天

#### Growth Projections
- **水平扩展**: 
  - 无状态服务，支持多实例部署
  - 数据库读写分离（未来）
  - Redis 集群（未来）
- **垂直扩展**: 
  - 数据库索引优化
  - 缓存策略优化

## Technical Decisions & Rationale

### Decision Log

1. **选择 gRPC 而非 REST API**
   - **原因**: 
     - 类型安全（protobuf）
     - 性能更好（二进制协议）
     - 支持流式传输（未来实时功能）
     - 代码生成，减少手写代码
   - **权衡**: 需要前端使用 gRPC-Web，增加复杂度

2. **采用 DDD 架构**
   - **原因**:
     - 业务逻辑清晰，易于维护
     - 领域模型与基础设施解耦
     - 便于测试（依赖注入）
     - 支持未来功能扩展
   - **权衡**: 初期开发复杂度较高，但长期收益大

3. **PostgreSQL + Redis 组合**
   - **原因**:
     - PostgreSQL: 关系型数据，支持复杂查询和全文搜索
     - Redis: 高性能缓存，支持复杂数据结构
     - 两者互补，满足不同场景需求
   - **权衡**: 需要维护两套存储系统

4. **使用 TypeScript 而非 JavaScript**
   - **原因**:
     - 类型安全，减少运行时错误
     - 更好的 IDE 支持
     - 代码可读性更好
   - **权衡**: 需要编译步骤，但收益远大于成本

5. **前端使用 React + 现代构建工具**
   - **原因**:
     - React 生态成熟，组件库丰富
     - Vite 构建速度快，开发体验好
     - 支持 SSR（未来版本）
   - **权衡**: 学习曲线，但团队熟悉度高

## Testing & Verification Strategy

### 测试金字塔

```
        /\
       /E2E\         少量端到端测试
      /────\
     /Integration\  集成测试（数据库、Redis）
    /────────────\
   /   Unit Test   \ 大量单元测试
  /────────────────\
```

### 测试要求

#### 1. 单元测试
- **覆盖率要求**: >= 70% (核心业务逻辑 >= 90%)
- **测试范围**:
  - Domain Layer: 所有实体、值对象、领域服务
  - Application Layer: 所有 Use Cases
  - Infrastructure Layer: Repository 实现
- **Mock 策略**: 
  - 使用 `gomock` 生成接口 mock
  - 测试时隔离外部依赖（数据库、Redis）

#### 2. 集成测试
- **测试范围**:
  - Repository 与数据库集成
  - Redis 缓存集成
  - gRPC 服务端集成
- **测试环境**: 
  - 使用 Docker Compose 启动测试数据库和 Redis
  - 每个测试用例使用独立的事务，测试后回滚

#### 3. 端到端测试
- **测试范围**:
  - 完整的用户流程（发布、查看、搜索）
  - API 端到端测试
- **工具**: 
  - 后端: `httptest` + gRPC 测试客户端
  - 前端: Playwright / Cypress

### 三方组件验证流程

**重要原则**: 设计三方组件时，必须保证调通后再进入下一步。

#### 验证步骤

1. **组件设计阶段**
   - 定义接口和依赖
   - 编写接口测试用例

2. **实现阶段**
   - 实现组件功能
   - 编写单元测试

3. **集成验证阶段** ⚠️ **必须完成此步骤**
   - 启动真实依赖（PostgreSQL、Redis）
   - 编写集成测试，验证组件与依赖的交互
   - **确保所有测试通过**
   - 验证错误处理和边界情况

4. **文档更新**
   - 更新组件使用文档
   - 记录已知问题和限制

5. **进入下一步**
   - 只有集成验证通过后，才能继续开发依赖该组件的功能

#### 示例：Redis 缓存组件验证

```go
// 1. 定义接口
type CacheRepository interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

// 2. 实现
type RedisCache struct {
    client *redis.Client
}

// 3. 集成测试（必须通过）
func TestRedisCache_Integration(t *testing.T) {
    // 启动真实 Redis（Docker）
    client := setupRedis(t)
    cache := NewRedisCache(client)
    
    // 测试 Get/Set
    err := cache.Set(ctx, "test", "value", time.Minute)
    require.NoError(t, err)
    
    val, err := cache.Get(ctx, "test")
    require.NoError(t, err)
    require.Equal(t, "value", val)
    
    // 测试 TTL
    // 测试错误处理
    // ...
}

// 4. 只有测试通过后，才能在其他地方使用
```

### 持续集成

- **CI Pipeline**:
  1. 代码检查 (lint)
  2. 单元测试
  3. 集成测试（需要 Docker）
  4. 构建
  5. 部署到测试环境

- **测试环境**: 
  - 使用 Docker Compose 提供 PostgreSQL 和 Redis
  - 每次 CI 运行都使用干净的环境

## Known Limitations

1. **第一版本不支持用户认证**
   - **影响**: 无法追踪发布者，无法实现个人中心
   - **解决方案**: 未来版本引入匿名用户系统

2. **内容审核依赖人工**
   - **影响**: 审核效率低，可能存在延迟
   - **解决方案**: 未来引入自动审核 + 人工复核

3. **搜索功能基于 PostgreSQL 全文搜索**
   - **影响**: 性能和功能有限
   - **解决方案**: 未来引入 Elasticsearch 或 Meilisearch

4. **单数据库架构**
   - **影响**: 高并发时可能成为瓶颈
   - **解决方案**: 未来实现读写分离、分库分表

