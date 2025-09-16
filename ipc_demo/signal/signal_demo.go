package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// 基本信号处理示例
func basicSignalDemo() {
	fmt.Println("=== 基本信号处理示例 ===")

	// 创建信号通道
	sigChan := make(chan os.Signal, 1)

	// 注册要捕获的信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

	fmt.Println("程序运行中，发送信号进行测试:")
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Println("可以使用以下命令测试:")
	fmt.Printf("kill -USR1 %d  # 发送SIGUSR1信号\n", os.Getpid())
	fmt.Printf("kill -USR2 %d  # 发送SIGUSR2信号\n", os.Getpid())
	fmt.Printf("kill -TERM %d  # 发送SIGTERM信号\n", os.Getpid())
	fmt.Println("或者按 Ctrl+C 发送SIGINT信号")

	// 启动一个goroutine定期发送自定义信号
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("自动发送SIGUSR1信号...")
		syscall.Kill(os.Getpid(), syscall.SIGUSR1)

		time.Sleep(2 * time.Second)
		fmt.Println("自动发送SIGUSR2信号...")
		syscall.Kill(os.Getpid(), syscall.SIGUSR2)
	}()

	// 处理信号
	for i := 0; i < 3; i++ {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT:
			fmt.Println("收到SIGINT信号 (Ctrl+C)")
		case syscall.SIGTERM:
			fmt.Println("收到SIGTERM信号")
		case syscall.SIGUSR1:
			fmt.Println("收到SIGUSR1信号 - 执行自定义操作1")
		case syscall.SIGUSR2:
			fmt.Println("收到SIGUSR2信号 - 执行自定义操作2")
		default:
			fmt.Printf("收到未知信号: %v\n", sig)
		}

		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			fmt.Println("程序即将退出...")
			break
		}
	}
}

// 优雅关闭示例
func gracefulShutdownDemo() {
	fmt.Println("\n=== 优雅关闭示例 ===")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 模拟工作任务
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("执行工作任务...")
			case <-ctx.Done():
				fmt.Println("工作任务收到停止信号，正在清理...")
				time.Sleep(1 * time.Second) // 模拟清理工作
				fmt.Println("工作任务清理完成")
				return
			}
		}
	}()

	// 信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("优雅关闭演示运行中 (PID: %d)\n", os.Getpid())
	fmt.Println("按 Ctrl+C 或发送SIGTERM信号来测试优雅关闭")

	// 自动发送信号进行演示
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("自动发送SIGTERM信号进行演示...")
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	// 等待信号
	sig := <-sigChan
	fmt.Printf("收到信号 %v，开始优雅关闭...\n", sig)

	// 取消上下文，通知所有goroutine停止
	cancel()

	// 等待所有goroutine完成
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// 设置超时
	select {
	case <-done:
		fmt.Println("所有任务已完成，程序正常退出")
	case <-time.After(5 * time.Second):
		fmt.Println("超时，强制退出")
	}
}

// 进程间信号通信示例
func interProcessSignalDemo() {
	fmt.Println("\n=== 进程间信号通信示例 ===")

	// 创建子进程
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	if err != nil {
		log.Fatal("启动子进程失败:", err)
	}

	childPID := cmd.Process.Pid
	fmt.Printf("子进程PID: %d\n", childPID)

	// 等待一段时间后向子进程发送信号
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Printf("向子进程 %d 发送SIGUSR1信号\n", childPID)
		syscall.Kill(childPID, syscall.SIGUSR1)

		time.Sleep(2 * time.Second)
		fmt.Printf("向子进程 %d 发送SIGTERM信号\n", childPID)
		syscall.Kill(childPID, syscall.SIGTERM)
	}()

	// 等待子进程结束
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("子进程退出: %v\n", err)
	} else {
		fmt.Println("子进程正常退出")
	}
}

// 信号量模拟
type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

// 信号量示例
func semaphoreDemo() {
	fmt.Println("\n=== 信号量示例 ===")

	// 创建容量为2的信号量
	sem := NewSemaphore(2)
	var wg sync.WaitGroup

	// 启动5个goroutine，但只允许2个同时运行
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			fmt.Printf("Goroutine %d 等待获取信号量...\n", id)
			sem.Acquire()
			fmt.Printf("Goroutine %d 获取到信号量，开始工作\n", id)

			// 模拟工作
			time.Sleep(2 * time.Second)

			fmt.Printf("Goroutine %d 工作完成，释放信号量\n", id)
			sem.Release()
		}(i)
	}

	wg.Wait()
	fmt.Println("所有任务完成")
}

// 自定义信号处理器
type SignalHandler struct {
	handlers map[os.Signal]func()
	sigChan  chan os.Signal
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewSignalHandler() *SignalHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &SignalHandler{
		handlers: make(map[os.Signal]func()),
		sigChan:  make(chan os.Signal, 10),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (sh *SignalHandler) Register(sig os.Signal, handler func()) {
	sh.handlers[sig] = handler
	signal.Notify(sh.sigChan, sig)
}

func (sh *SignalHandler) Start() {
	sh.wg.Add(1)
	go func() {
		defer sh.wg.Done()

		for {
			select {
			case sig := <-sh.sigChan:
				if handler, exists := sh.handlers[sig]; exists {
					fmt.Printf("处理信号: %v\n", sig)
					handler()
				}
			case <-sh.ctx.Done():
				return
			}
		}
	}()
}

func (sh *SignalHandler) Stop() {
	sh.cancel()
	sh.wg.Wait()
}

// 自定义信号处理器示例
func customSignalHandlerDemo() {
	fmt.Println("\n=== 自定义信号处理器示例 ===")

	handler := NewSignalHandler()

	// 注册信号处理函数
	handler.Register(syscall.SIGUSR1, func() {
		fmt.Println("执行SIGUSR1处理逻辑：重新加载配置")
	})

	handler.Register(syscall.SIGUSR2, func() {
		fmt.Println("执行SIGUSR2处理逻辑：打印统计信息")
	})

	handler.Register(syscall.SIGINT, func() {
		fmt.Println("执行SIGINT处理逻辑：准备退出")
		handler.Stop()
	})

	// 启动信号处理器
	handler.Start()

	fmt.Printf("自定义信号处理器运行中 (PID: %d)\n", os.Getpid())

	// 自动发送信号进行演示
	go func() {
		time.Sleep(1 * time.Second)
		syscall.Kill(os.Getpid(), syscall.SIGUSR1)

		time.Sleep(1 * time.Second)
		syscall.Kill(os.Getpid(), syscall.SIGUSR2)

		time.Sleep(1 * time.Second)
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()

	// 等待处理器停止
	handler.wg.Wait()
	fmt.Println("自定义信号处理器已停止")
}

func main() {
	fmt.Println("Go 信号通信示例")
	fmt.Println("===============")

	basicSignalDemo()
	gracefulShutdownDemo()
	interProcessSignalDemo()
	semaphoreDemo()
	customSignalHandlerDemo()
}
