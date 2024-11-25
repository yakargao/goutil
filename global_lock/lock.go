package global_lock

import "sync"

type GlobalLock struct {
	locks map[string]*sync.Mutex
	mu    sync.Mutex
}

func NewGlobalLock() *GlobalLock {
	return &GlobalLock{
		locks: make(map[string]*sync.Mutex),
	}
}
func (g *GlobalLock) TryLock(key string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.locks[key]; !ok {
		g.locks[key] = &sync.Mutex{}
	}
	return g.locks[key].TryLock()
}

func (g *GlobalLock) Lock(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.locks[key]; !ok {
		g.locks[key] = &sync.Mutex{}
	}
	g.locks[key].Lock()
}
func (g *GlobalLock) Unlock(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.locks[key]; ok {
		g.locks[key].Unlock()
	}
}
