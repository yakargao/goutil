package global_lock

import (
	"fmt"
	"testing"
)

func TestGlobalLock_Lock(t *testing.T) {
	g := NewGlobalLock()
	userID := "user123"

	// 模拟用户执行某个操作
	go func() {
		fmt.Println("User trying to lock:", userID)
		g.Lock(userID)
		fmt.Println("User locked:", userID)
		// 模拟操作
		// time.Sleep(2 * time.Second)
		g.Unlock(userID)
		fmt.Println("User unlocked:", userID)
	}()
	// 为了演示输出，暂停主协程
	select {}
}

func TestGlobalLock_TryLock(t *testing.T) {
	g := NewGlobalLock()
	userID := "user456"
	if g.TryLock(userID) {
		fmt.Println("User locked:", userID)
		if g.TryLock(userID) {
			fmt.Println("User locked again:", userID)
		} else {
			fmt.Println("User failed to lock:", userID)
		}
		// 模拟操作
		// time.Sleep(2 * time.Second)
		g.Unlock(userID)
	}
}
