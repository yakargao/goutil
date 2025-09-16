package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// 消息结构
type Message struct {
	ID      int
	Content string
	From    string
	To      string
}

// 简单消息队列实现
type MessageQueue struct {
	queue   chan Message
	workers int
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// 创建消息队列
func NewMessageQueue(bufferSize, workers int) *MessageQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &MessageQueue{
		queue:   make(chan Message, bufferSize),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// 发送消息
func (mq *MessageQueue) Send(msg Message) error {
	select {
	case mq.queue <- msg:
		return nil
	case <-mq.ctx.Done():
		return fmt.Errorf("消息队列已关闭")
	default:
		return fmt.Errorf("消息队列已满")
	}
}

// 启动消费者
func (mq *MessageQueue) StartConsumers(handler func(Message)) {
	for i := 0; i < mq.workers; i++ {
		mq.wg.Add(1)
		go func(workerID int) {
			defer mq.wg.Done()
			fmt.Printf("消费者 %d 启动\n", workerID)

			for {
				select {
				case msg := <-mq.queue:
					fmt.Printf("消费者 %d 处理消息: %+v\n", workerID, msg)
					handler(msg)
				case <-mq.ctx.Done():
					fmt.Printf("消费者 %d 停止\n", workerID)
					return
				}
			}
		}(i)
	}
}

// 关闭消息队列
func (mq *MessageQueue) Close() {
	mq.cancel()
	close(mq.queue)
	mq.wg.Wait()
}

// 基本消息队列示例
func basicMessageQueueDemo() {
	fmt.Println("=== 基本消息队列示例 ===")

	// 创建消息队列
	mq := NewMessageQueue(10, 2)

	// 消息处理函数
	handler := func(msg Message) {
		fmt.Printf("处理消息: ID=%d, Content=%s, From=%s, To=%s\n",
			msg.ID, msg.Content, msg.From, msg.To)
		time.Sleep(500 * time.Millisecond) // 模拟处理时间
	}

	// 启动消费者
	mq.StartConsumers(handler)

	// 发送消息
	messages := []Message{
		{1, "Hello", "Producer1", "Consumer"},
		{2, "World", "Producer1", "Consumer"},
		{3, "Go", "Producer2", "Consumer"},
		{4, "IPC", "Producer2", "Consumer"},
		{5, "Demo", "Producer1", "Consumer"},
	}

	for _, msg := range messages {
		err := mq.Send(msg)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
		} else {
			fmt.Printf("发送消息: %+v\n", msg)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// 等待处理完成
	time.Sleep(3 * time.Second)
	mq.Close()
}

// 优先级消息队列
type PriorityMessage struct {
	Message
	Priority int
}

type PriorityQueue struct {
	high   chan PriorityMessage
	medium chan PriorityMessage
	low    chan PriorityMessage
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewPriorityQueue(bufferSize int) *PriorityQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &PriorityQueue{
		high:   make(chan PriorityMessage, bufferSize),
		medium: make(chan PriorityMessage, bufferSize),
		low:    make(chan PriorityMessage, bufferSize),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (pq *PriorityQueue) Send(msg PriorityMessage) error {
	var targetQueue chan PriorityMessage

	switch {
	case msg.Priority >= 8:
		targetQueue = pq.high
	case msg.Priority >= 5:
		targetQueue = pq.medium
	default:
		targetQueue = pq.low
	}

	select {
	case targetQueue <- msg:
		return nil
	case <-pq.ctx.Done():
		return fmt.Errorf("优先级队列已关闭")
	default:
		return fmt.Errorf("优先级队列已满")
	}
}

func (pq *PriorityQueue) StartConsumer(handler func(PriorityMessage)) {
	pq.wg.Add(1)
	go func() {
		defer pq.wg.Done()

		for {
			select {
			case msg := <-pq.high:
				fmt.Printf("处理高优先级消息: %+v\n", msg)
				handler(msg)
			case msg := <-pq.medium:
				fmt.Printf("处理中优先级消息: %+v\n", msg)
				handler(msg)
			case msg := <-pq.low:
				fmt.Printf("处理低优先级消息: %+v\n", msg)
				handler(msg)
			case <-pq.ctx.Done():
				return
			}
		}
	}()
}

func (pq *PriorityQueue) Close() {
	pq.cancel()
	close(pq.high)
	close(pq.medium)
	close(pq.low)
	pq.wg.Wait()
}

// 优先级消息队列示例
func priorityMessageQueueDemo() {
	fmt.Println("\n=== 优先级消息队列示例 ===")

	pq := NewPriorityQueue(5)

	handler := func(msg PriorityMessage) {
		fmt.Printf("处理消息: Priority=%d, Content=%s\n", msg.Priority, msg.Content)
		time.Sleep(300 * time.Millisecond)
	}

	pq.StartConsumer(handler)

	// 发送不同优先级的消息
	messages := []PriorityMessage{
		{Message{1, "低优先级消息1", "Producer", "Consumer"}, 2},
		{Message{2, "高优先级消息1", "Producer", "Consumer"}, 9},
		{Message{3, "中优先级消息1", "Producer", "Consumer"}, 6},
		{Message{4, "低优先级消息2", "Producer", "Consumer"}, 1},
		{Message{5, "高优先级消息2", "Producer", "Consumer"}, 10},
		{Message{6, "中优先级消息2", "Producer", "Consumer"}, 5},
	}

	for _, msg := range messages {
		err := pq.Send(msg)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
		} else {
			fmt.Printf("发送消息: Priority=%d, Content=%s\n", msg.Priority, msg.Content)
		}
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)
	pq.Close()
}

// 发布订阅模式
type PubSub struct {
	subscribers map[string][]chan Message
	mutex       sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewPubSub() *PubSub {
	ctx, cancel := context.WithCancel(context.Background())
	return &PubSub{
		subscribers: make(map[string][]chan Message),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (ps *PubSub) Subscribe(topic string, bufferSize int) <-chan Message {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ch := make(chan Message, bufferSize)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *PubSub) Publish(topic string, msg Message) {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if subscribers, exists := ps.subscribers[topic]; exists {
		for _, ch := range subscribers {
			select {
			case ch <- msg:
			case <-ps.ctx.Done():
				return
			default:
				fmt.Printf("订阅者缓冲区已满，丢弃消息: %+v\n", msg)
			}
		}
	}
}

func (ps *PubSub) Close() {
	ps.cancel()
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	for _, subscribers := range ps.subscribers {
		for _, ch := range subscribers {
			close(ch)
		}
	}
}

// 发布订阅示例
func pubSubDemo() {
	fmt.Println("\n=== 发布订阅模式示例 ===")

	pubsub := NewPubSub()
	defer pubsub.Close()

	// 订阅者1
	sub1 := pubsub.Subscribe("news", 5)
	go func() {
		for msg := range sub1 {
			fmt.Printf("订阅者1收到新闻: %s\n", msg.Content)
		}
	}()

	// 订阅者2
	sub2 := pubsub.Subscribe("news", 5)
	go func() {
		for msg := range sub2 {
			fmt.Printf("订阅者2收到新闻: %s\n", msg.Content)
		}
	}()

	// 订阅者3订阅体育新闻
	sub3 := pubsub.Subscribe("sports", 5)
	go func() {
		for msg := range sub3 {
			fmt.Printf("订阅者3收到体育新闻: %s\n", msg.Content)
		}
	}()

	// 发布消息
	time.Sleep(100 * time.Millisecond)

	pubsub.Publish("news", Message{1, "重要新闻：Go 1.21发布", "NewsPublisher", ""})
	pubsub.Publish("news", Message{2, "技术新闻：新的并发模式", "TechPublisher", ""})
	pubsub.Publish("sports", Message{3, "体育新闻：世界杯开始", "SportsPublisher", ""})
	pubsub.Publish("news", Message{4, "经济新闻：股市上涨", "EconPublisher", ""})

	time.Sleep(2 * time.Second)
}

func main() {
	fmt.Println("Go 消息队列通信示例")
	fmt.Println("==================")

	basicMessageQueueDemo()
	priorityMessageQueueDemo()
	pubSubDemo()
}
