# prepare2go

> prepare2go 一个非常适合 phper 转 gopher 的web开发脚手架，封装了
> [iris](https://github.com/kataras/iris)、
> [cobra](https://github.com/spf13/cobra)、
> [gorm](https://github.com/go-gorm/gorm)、
> [go-redis](https://github.com/redis/go-redis)、
> [go-resty](https://github.com/go-resty/resty)、
> [zap](https://github.com/uber-go/zap)
> 等优秀开源包，开箱即用、快速部署。

## 起步

> 以下命令需要你已经安装好 **golang** 编译器，如果你已经有了 **golang** 的开发经验并且本地已经配置好 **golang**
> 开发环境请直接跳到 ```step4```

- step1. 确认已经配置好 ```GOPATH```:

```shell
# 查看当前GOPATH的配置
go env GOPATH

# 如果没有输出，需要先设置 go env 中的 GOPATH，然后写入系统 PATH
go env -w GOPATH=<path to your go-path>
```

- step2. 将 ```GOPATH``` 写入系统 ```PATH``` 以便使用命令行工具

```shell
echo -e "\nexport GOPATH=$(go env GOPATH)\n\nexport PATH=\$PATH:$GOPATH/bin" >> ~/.bash_profile
```

- step3. 建议开启: ```GO111MODULE```，并配置go环境变量:

```shell
go env -w GO111MODULE=on
go env -w GOPROXY="https://goproxy.cn,direct"
```

- step4. 配置环境变量

```shell
echo 'export RUNTIME_ENV=local' >> ~/.bash_profile && source ~/.bash_profile
```

- step5. 安装命令行工具 aurora

```shell
go install github.com/stubborn-gaga/aurora@latest
```

- 已经为你编写好了命令行工具 [aurora](https://github.com/stubborn-gaga-0805/aurora/blob/main/README.md)
  你可以直接使用她来管理和运行你的项目:

```shell
# 将./configs/config.local.yaml配置文件中的配置修改为你自己的

aurora init
aurora run

# 输出一下信息表示服务成功启动：
# HttpServer 启动成功，监听地址: 0.0.0.0:8800
```

- 启动服务后可以通过curl命令测试服务是否可用：

```shell
curl http://127.0.0.1:8800/ping
# 输出：[2023-05-22 00:29:13] ...pong! 表示成功
```

- 默认的二进制文件为项目的根目录下 ```./bin/server``` 文件。你也可以手动编译项目、指定二进制文件的输出目录:

```shell
aurora build -o ./bin/server  # 自定义编译二进制文件
```

> PS: 命令行工具 [aurora](https://github.com/stubborn-gaga-0805/aurora/blob/main/README.md)
> 提供了丰富的项目管理功能。可以运行 ```aurora -h``` 来了解帮助信息

## 目录结构

- 核心目录如下所示，项目层级以及核心代码的实现原理将单独在doc中展示

```shell
.
├── api # 请求控制层，这一层主要进行参数校验，组装请求结构体给service
│   ├── ecode # 定义错误码和报错信息
│   ├── middleware  # 拦截器（中间件）
│   ├── request # 请求参数的结构体, 包含参数校验逻辑
│   ├── response  # 返回逻辑代码
│   └── router  # 路由文件
├── bin # 二进制文件目录(.gitignore)
├── cmd   # 命令行实现代码
├── configs # 配置文件目录
│   ├── conf  # 配置项的结构体
│   ├── config.local.yaml # 不同环境的配置文件
├── internal  # 业务逻辑层
│   ├── consts  # 常量文件目录
│   ├── controller  # 控制器文件目录
│   ├── job # 任务文件目录
│   ├── repo  # 数据访问层
│   │   ├── entity  # 数据实体文件目录
│   │   ├── orm # models文件目录
│   ├── service # service层
│   └── server  # 服务文件实现目录
├── logs  # 本地日志文件目录(.gitignore)
├── main.go # main文件
└── pkg # 第三方包目录

```

## 架构说明

### api

> **api** 层主要制定路由规则并处理前端请求, 路由到的请求在这一层里转发到 **controller** 实现处理。
> **中间件**、**拦截器** 也是在这一层实现。
> 请求、返回的结构体，异常错误码也在这一层里来定义
 ---------

### controller

> **controller** 层主要是对于参数，为 **server** 的流入数据进行规范化处理、同样对于 **server** 的返回。
> 在这一层，前端传过来的参数将进行校验，然后传给 **service** 中进行业务逻辑处理。
> 同样，对于 **service** 的运行结果，也需要在这一层进行判断然后将成功的运行结果或者失败的异常处理返回到前端。
 ---------

### service

> **service** 层主要是对实现业务逻辑处理。
> 在这一层，主要应当关系业务逻辑的流转。
> 与第三方组建的交互也应当在这一层处理。
> 业务数据逻辑化处理后在这一层被抽象成结构化数据，同样、结构化数据在这一层被转换成业务数据返回给上游
---------

### repo

> **repo** 层主要与持久化数据交互
> 在这一层，业务数据被转换成结构化数据实现持久化。这一层里需要维护数据驱动（如：mysql、redis、es的连接）
---------