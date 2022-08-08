package metpoll

import (
	"sync"
	"sync/atomic"
)

type engine struct {
	ln         *listener
	wg         sync.WaitGroup
	once       sync.Once
	cond       *sync.Cond
	inShutdown int32
}

func (eng *engine) isInShutdown() bool {
	return atomic.LoadInt32(&eng.inShutdown) == 1
}

func (eng *engine) waitForShutdown() {
	eng.cond.L.Lock()
	eng.cond.Wait()
	eng.cond.L.Unlock()
}

func (eng *engine) signalShutdown() {
	eng.once.Do(func() {
		eng.cond.L.Lock()
		eng.cond.Signal()
		eng.cond.L.Unlock()
	})
}
