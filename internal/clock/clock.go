package clock

import "sync"

type Clock interface{ Now() int64 }

type Fixed struct {
	mu    sync.RWMutex
	value int64
}

func NewFixed(value int64) *Fixed { return &Fixed{value: value} }

func (c *Fixed) Now() int64 { c.mu.RLock(); defer c.mu.RUnlock(); return c.value }

func (c *Fixed) Set(value int64) { c.mu.Lock(); c.value = value; c.mu.Unlock() }

func (c *Fixed) Advance(delta int64) int64 {
	c.mu.Lock()
	c.value += delta
	value := c.value
	c.mu.Unlock()
	return value
}

func Sequence(start, step, count int64) []int64 {
	if count < 0 {
		count = 0
	}
	result := make([]int64, 0, count)
	for i := int64(0); i < count; i++ {
		result = append(result, start+i*step)
	}
	return result
}
