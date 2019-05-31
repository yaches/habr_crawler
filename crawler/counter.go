package crawler

import (
	"sync"
)

type counter struct {
	mu   *sync.Mutex
	cntr int
}

func NewCounter(i int) counter {
	return counter{
		mu:   &sync.Mutex{},
		cntr: i,
	}
}

func (c *counter) Inc() {
	c.mu.Lock()
	c.cntr++
	c.mu.Unlock()
}

func (c *counter) Dec() {
	c.mu.Lock()
	c.cntr--
	c.mu.Unlock()
}

func (c *counter) Zero() bool {
	return c.cntr == 0
}
