package ll

import (
	"errors"
	"sync"
)

var (
	ErrMaxActiveProcessCountReached = errors.New("Max Active Process Count reached")
)

// To keep track of the number of active requests (being processed)
type ActiveProcessCounter struct {
	max     int
	mu      sync.Mutex
	actives map[string]interface{}
}

func NewActiveProcessCounter(max int) *ActiveProcessCounter {
	var tmp *ActiveProcessCounter = &ActiveProcessCounter{}

	tmp.max = max
	tmp.actives = make(map[string]interface{})
	return tmp
}

func (apc *ActiveProcessCounter) Register(id string) error {
	apc.mu.Lock()
	if len(apc.actives) >= apc.max {
		apc.mu.Unlock()
		return ErrMaxActiveProcessCountReached
	}
	apc.actives[id] = struct{}{}
	apc.mu.Unlock()
	return nil
}

func (apc *ActiveProcessCounter) Unregister(id string) {
	apc.mu.Lock()
	delete(apc.actives, id)
	apc.mu.Unlock()
}

func (apc *ActiveProcessCounter) Count() int {
	var result int
	apc.mu.Lock()
	result = len(apc.actives)
	apc.mu.Unlock()
	return result
}
