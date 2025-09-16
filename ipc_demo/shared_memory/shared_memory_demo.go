package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// 共享内存结构
type SharedData struct {
	Counter int64
	Message [256]byte
	Flag    int32
}

// 内存映射文件示例
func mmapDemo() {
	fmt.Println("=== 内存映射文件示例 ===")

	filename := "/tmp/shared_memory_demo"
	size := 4096

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("创建文件失败:", err)
	}
	defer os.Remove(filename)

	// 设置文件大小
	err = file.Truncate(int64(size))
	if err != nil {
		log.Fatal("设置文件大小失败:", err)
	}

	// 内存映射
	data, err := syscall.Mmap(int(file.Fd()), 0, size,
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Fatal("内存映射失败:", err)
	}
	defer syscall.Munmap(data)
	defer file.Close()

	// 写入数据
	message := "Hello from shared memory!"
	copy(data, message)

	fmt.Printf("写入共享内存: %s\n", message)

	// 启动另一个goroutine读取数据
	go func() {
		time.Sleep(100 * time.Millisecond)

		// 重新打开文件并映射
		file2, err := os.Open(filename)
		if err != nil {
			log.Printf("打开文件失败: %v", err)
			return
		}
		defer file2.Close()

		data2, err := syscall.Mmap(int(file2.Fd()), 0, size,
			syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			log.Printf("内存映射失败: %v", err)
			return
		}
		defer syscall.Munmap(data2)

		// 读取数据
		readMessage := string(data2[:len(message)])
		fmt.Printf("从共享内存读取: %s\n", readMessage)
	}()

	time.Sleep(500 * time.Millisecond)
}

// 使用unsafe包的共享内存示例
func unsafeSharedMemoryDemo() {
	fmt.Println("\n=== unsafe包共享内存示例 ===")

	// 创建共享数据结构
	sharedData := &SharedData{
		Counter: 0,
		Flag:    0,
	}

	// 初始化消息
	message := "Shared memory with unsafe"
	copy(sharedData.Message[:], message)

	var wg sync.WaitGroup

	// 启动多个goroutine操作共享内存
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				// 原子操作增加计数器
				ptr := (*int64)(unsafe.Pointer(&sharedData.Counter))
				for {
					old := *ptr
					if syscall.CompareAndSwapInt64((*int64)(unsafe.Pointer(ptr)), old, old+1) {
						break
					}
				}

				fmt.Printf("Goroutine %d: Counter = %d\n", id, sharedData.Counter)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	// 读取最终结果
	fmt.Printf("最终计数器值: %d\n", sharedData.Counter)
	fmt.Printf("共享消息: %s\n", string(sharedData.Message[:len(message)]))
}

// 基于channel的共享内存模拟
type SharedMemorySimulator struct {
	data   map[string]interface{}
	mutex  sync.RWMutex
	notify chan string
}

func NewSharedMemorySimulator() *SharedMemorySimulator {
	return &SharedMemorySimulator{
		data:   make(map[string]interface{}),
		notify: make(chan string, 100),
	}
}

func (sms *SharedMemorySimulator) Write(key string, value interface{}) {
	sms.mutex.Lock()
	defer sms.mutex.Unlock()

	sms.data[key] = value

	// 通知有数据更新
	select {
	case sms.notify <- key:
	default:
	}
}

func (sms *SharedMemorySimulator) Read(key string) (interface{}, bool) {
	sms.mutex.RLock()
	defer sms.mutex.RUnlock()

	value, exists := sms.data[key]
	return value, exists
}

func (sms *SharedMemorySimulator) Subscribe() <-chan string {
	return sms.notify
}

func (sms *SharedMemorySimulator) Close() {
	close(sms.notify)
}

// 共享内存模拟器示例
func sharedMemorySimulatorDemo() {
	fmt.Println("\n=== 共享内存模拟器示例 ===")

	sms := NewSharedMemorySimulator()
	defer sms.Close()

	// 启动监听器
	go func() {
		for key := range sms.Subscribe() {
			if value, exists := sms.Read(key); exists {
				fmt.Printf("检测到数据更新: %s = %v\n", key, value)
			}
		}
	}()

	var wg sync.WaitGroup

	// 启动写入者
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 3; j++ {
				key := fmt.Sprintf("writer_%d_data", id)
				value := fmt.Sprintf("数据来自写入者%d，第%d次写入", id, j+1)

				sms.Write(key, value)
				fmt.Printf("写入者%d写入: %s = %s\n", id, key, value)

				time.Sleep(200 * time.Millisecond)
			}
		}(i)
	}

	// 启动读取者
	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(500 * time.Millisecond)

		for i := 0; i < 3; i++ {
			key := fmt.Sprintf("writer_%d_data", i)
			if value, exists := sms.Read(key); exists {
				fmt.Printf("读取者读取: %s = %v\n", key, value)
			}
		}
	}()

	wg.Wait()
	time.Sleep(500 * time.Millisecond)
}

// 生产者消费者模式的共享缓冲区
type SharedBuffer struct {
	buffer   []interface{}
	size     int
	head     int
	tail     int
	count    int
	mutex    sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

func NewSharedBuffer(size int) *SharedBuffer {
	sb := &SharedBuffer{
		buffer: make([]interface{}, size),
		size:   size,
	}
	sb.notEmpty = sync.NewCond(&sb.mutex)
	sb.notFull = sync.NewCond(&sb.mutex)
	return sb
}

func (sb *SharedBuffer) Put(item interface{}) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()

	// 等待缓冲区不满
	for sb.count == sb.size {
		sb.notFull.Wait()
	}

	sb.buffer[sb.tail] = item
	sb.tail = (sb.tail + 1) % sb.size
	sb.count++

	sb.notEmpty.Signal()
}

func (sb *SharedBuffer) Get() interface{} {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()

	// 等待缓冲区不空
	for sb.count == 0 {
		sb.notEmpty.Wait()
	}

	item := sb.buffer[sb.head]
	sb.head = (sb.head + 1) % sb.size
	sb.count--

	sb.notFull.Signal()
	return item
}

// 共享缓冲区示例
func sharedBufferDemo() {
	fmt.Println("\n=== 共享缓冲区示例 ===")

	buffer := NewSharedBuffer(5)
	var wg sync.WaitGroup

	// 生产者
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				item := fmt.Sprintf("生产者%d-项目%d", id, j)
				buffer.Put(item)
				fmt.Printf("生产者%d生产: %s\n", id, item)
				time.Sleep(200 * time.Millisecond)
			}
		}(i)
	}

	// 消费者
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				item := buffer.Get()
				fmt.Printf("消费者%d消费: %v\n", id, item)
				time.Sleep(300 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}

func main() {
	fmt.Println("Go 共享内存通信示例")
	fmt.Println("==================")

	mmapDemo()
	unsafeSharedMemoryDemo()
	sharedMemorySimulatorDemo()
	sharedBufferDemo()
}
