package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// 匿名管道示例
func anonymousPipeDemo() {
	fmt.Println("=== 匿名管道示例 ===")

	// 创建管道
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal("创建管道失败:", err)
	}
	defer r.Close()
	defer w.Close()

	// 启动写入goroutine
	go func() {
		defer w.Close()
		for i := 0; i < 5; i++ {
			message := fmt.Sprintf("消息 %d", i)
			fmt.Printf("写入: %s\n", message)
			w.WriteString(message + "\n")
			time.Sleep(time.Second)
		}
	}()

	// 读取数据
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("读取: %s\n", scanner.Text())
	}
}

// 命名管道示例
func namedPipeDemo() {
	fmt.Println("\n=== 命名管道示例 ===")

	pipeName := "/tmp/demo_pipe"

	// 创建命名管道
	err := syscall.Mkfifo(pipeName, 0666)
	if err != nil && !os.IsExist(err) {
		log.Fatal("创建命名管道失败:", err)
	}
	defer os.Remove(pipeName)

	// 启动写入进程
	go func() {
		time.Sleep(100 * time.Millisecond) // 确保读取端先准备好

		pipe, err := os.OpenFile(pipeName, os.O_WRONLY, 0)
		if err != nil {
			log.Printf("打开写入管道失败: %v", err)
			return
		}
		defer pipe.Close()

		for i := 0; i < 3; i++ {
			message := fmt.Sprintf("命名管道消息 %d\n", i)
			fmt.Printf("写入命名管道: %s", message)
			pipe.WriteString(message)
			time.Sleep(time.Second)
		}
	}()

	// 读取数据
	pipe, err := os.OpenFile(pipeName, os.O_RDONLY, 0)
	if err != nil {
		log.Fatal("打开读取管道失败:", err)
	}
	defer pipe.Close()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Printf("从命名管道读取: %s\n", scanner.Text())
	}
}

// 使用exec.Cmd的管道示例
func execPipeDemo() {
	fmt.Println("\n=== exec.Cmd 管道示例 ===")

	// 创建命令管道: echo "hello world" | wc -w
	cmd1 := exec.Command("echo", "hello world from golang pipe")
	cmd2 := exec.Command("wc", "-w")

	// 连接管道
	pipe, err := cmd1.StdoutPipe()
	if err != nil {
		log.Fatal("创建管道失败:", err)
	}

	cmd2.Stdin = pipe

	// 启动命令
	err = cmd1.Start()
	if err != nil {
		log.Fatal("启动cmd1失败:", err)
	}

	output, err := cmd2.Output()
	if err != nil {
		log.Fatal("执行cmd2失败:", err)
	}

	err = cmd1.Wait()
	if err != nil {
		log.Fatal("等待cmd1失败:", err)
	}

	fmt.Printf("管道输出结果: %s", output)
}

// 双向管道示例
func bidirectionalPipeDemo() {
	fmt.Println("\n=== 双向管道示例 ===")

	// 创建两个管道实现双向通信
	r1, w1, _ := os.Pipe() // 父进程写，子进程读
	r2, w2, _ := os.Pipe() // 子进程写，父进程读

	defer r1.Close()
	defer w1.Close()
	defer r2.Close()
	defer w2.Close()

	// 子goroutine
	go func() {
		defer r1.Close()
		defer w2.Close()

		// 读取父进程消息
		scanner := bufio.NewScanner(r1)
		for scanner.Scan() {
			received := scanner.Text()
			fmt.Printf("子进程收到: %s\n", received)

			// 回复消息
			reply := fmt.Sprintf("回复: %s", received)
			w2.WriteString(reply + "\n")
		}
	}()

	// 父进程
	go func() {
		defer w1.Close()

		messages := []string{"消息1", "消息2", "消息3"}
		for _, msg := range messages {
			fmt.Printf("父进程发送: %s\n", msg)
			w1.WriteString(msg + "\n")
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// 父进程读取回复
	scanner := bufio.NewScanner(r2)
	for scanner.Scan() {
		fmt.Printf("父进程收到回复: %s\n", scanner.Text())
	}
}

func main() {
	fmt.Println("Go 管道通信示例")
	fmt.Println("================")

	anonymousPipeDemo()
	namedPipeDemo()
	execPipeDemo()
	bidirectionalPipeDemo()
}
