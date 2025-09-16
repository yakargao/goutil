# Go 进程间通信 (IPC) 示例

这个项目包含了使用 Go 语言实现的各种进程间通信方式的示例代码。

## 项目结构

```
ipc_demo/
├── main.go                          # 主程序入口
├── README.md                        # 项目说明
├── pipe/
│   └── pipe_demo.go                 # 管道通信示例
├── message_queue/
│   └── message_queue_demo.go        # 消息队列示例
├── shared_memory/
│   └── shared_memory_demo.go        # 共享内存示例
└── signal/
    └── signal_demo.go               # 信号通信示例
```

## 运行示例

### 方式一：使用主程序运行

```bash
# 进入项目目录
cd ipc_demo

# 运行特定示例
go run main.go pipe           # 管道通信示例
go run main.go mq             # 消息队列示例
go run main.go shm            # 共享内存示例
go run main.go signal         # 信号通信示例

# 运行所有示例
go run main.go all
```

### 方式二：直接运行单个示例

```bash
# 管道通信示例
go run pipe/pipe_demo.go

# 消息队列示例
go run message_queue/message_queue_demo.go

# 共享内存示例
go run shared_memory/shared_memory_demo.go

# 信号通信示例
go run signal/signal_demo.go
```

## 示例说明

### 1. 管道通信 (pipe/pipe_demo.go)

演示了以下管道通信方式：
- **匿名管道**: 使用 `os.Pipe()` 创建的管道，适用于父子进程间通信
- **命名管道**: 使用 `syscall.Mkfifo()` 创建的FIFO，可用于无关进程间通信
- **exec.Cmd管道**: 使用 `exec.Command` 连接多个命令的管道
- **双向管道**: 使用两个管道实现双向通信

特点：
- 简单易用，适合流式数据传输
- 有序传输，先进先出
- 适合生产者-消费者模式

### 2. 消息队列 (message_queue/message_queue_demo.go)

演示了以下消息队列实现：
- **基本消息队列**: 使用 channel 实现的简单消息队列
- **优先级消息队列**: 支持不同优先级的消息处理
- **发布订阅模式**: 支持多个订阅者的消息分发

特点：
- 异步通信，解耦生产者和消费者
- 支持多对多通信
- 可以实现负载均衡和优先级处理

### 3. 共享内存 (shared_memory/shared_memory_demo.go)

演示了以下共享内存实现：
- **内存映射文件**: 使用 `syscall.Mmap()` 实现的文件映射
- **unsafe包共享内存**: 使用 unsafe 包直接操作内存
- **共享内存模拟器**: 使用 map 和 mutex 模拟共享内存
- **共享缓冲区**: 实现生产者-消费者模式的共享缓冲区

特点：
- 高性能，直接内存访问
- 需要同步机制防止竞态条件
- 适合大量数据的高速传输

### 4. 信号通信 (signal/signal_demo.go)

演示了以下信号通信方式：
- **基本信号处理**: 捕获和处理系统信号
- **优雅关闭**: 使用信号实现程序的优雅关闭
- **进程间信号**: 向其他进程发送信号
- **信号量**: 使用 channel 模拟信号量机制
- **自定义信号处理器**: 封装的信号处理框架

特点：
- 轻量级，系统级通信
- 适合进程控制和状态通知
- 异步处理，不阻塞程序执行

## 使用场景

### 管道 (Pipe)
- 命令行工具的数据流处理
- 父子进程间的数据传递
- 简单的生产者-消费者场景

### 消息队列 (Message Queue)
- 微服务间的异步通信
- 任务调度和分发
- 事件驱动架构

### 共享内存 (Shared Memory)
- 高性能计算场景
- 大数据量的进程间传输
- 实时系统的数据共享

### 信号 (Signal)
- 进程生命周期管理
- 配置重载和状态切换
- 系统监控和告警

## 注意事项

1. **并发安全**: 共享内存需要适当的同步机制
2. **资源清理**: 及时清理管道、文件描述符等资源
3. **错误处理**: 处理各种异常情况和边界条件
4. **平台兼容性**: 某些功能可能依赖特定操作系统

## 扩展学习

- 学习 Go 的 context 包进行更好的并发控制
- 了解 sync 包中的各种同步原语
- 研究 Go 的内存模型和并发安全
- 探索分布式系统中的进程间通信

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这些示例！