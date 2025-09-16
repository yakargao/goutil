package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Go 进程间通信 (IPC) 示例集合")
	fmt.Println("=============================")
	fmt.Println()

	if len(os.Args) < 2 {
		showUsage()
		return
	}

	demoType := os.Args[1]

	switch demoType {
	case "pipe":
		runDemo("pipe", "管道通信示例")
	case "mq", "message_queue":
		runDemo("message_queue", "消息队列示例")
	case "shm", "shared_memory":
		runDemo("shared_memory", "共享内存示例")
	case "signal":
		runDemo("signal", "信号通信示例")
	case "all":
		runAllDemos()
	default:
		fmt.Printf("未知的示例类型: %s\n", demoType)
		showUsage()
	}
}

func showUsage() {
	fmt.Println("用法: go run main.go <demo_type>")
	fmt.Println()
	fmt.Println("可用的示例类型:")
	fmt.Println("  pipe           - 管道通信示例")
	fmt.Println("  mq             - 消息队列示例")
	fmt.Println("  shm            - 共享内存示例")
	fmt.Println("  signal         - 信号通信示例")
	fmt.Println("  all            - 运行所有示例")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run main.go pipe")
	fmt.Println("  go run main.go mq")
	fmt.Println("  go run main.go all")
}

func runDemo(demoDir, description string) {
	fmt.Printf("运行 %s\n", description)
	fmt.Println(strings.Repeat("=", len(description)+3))

	demoPath := filepath.Join(demoDir, demoDir+"_demo.go")

	cmd := exec.Command("go", "run", demoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("运行示例失败: %v\n", err)
	}

	fmt.Println()
}

func runAllDemos() {
	demos := []struct {
		dir  string
		desc string
	}{
		{"pipe", "管道通信示例"},
		{"message_queue", "消息队列示例"},
		{"shared_memory", "共享内存示例"},
		{"signal", "信号通信示例"},
	}

	for i, demo := range demos {
		if i > 0 {
			fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
		}
		runDemo(demo.dir, demo.desc)
	}
}
