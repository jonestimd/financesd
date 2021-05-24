package server

import (
	"fmt"
	"sync"
	"time"
)

type IDSupplier struct {
	mu     sync.Mutex
	base   int64
	suffix int64
}

func NewIDSupplier() *IDSupplier {
	return &IDSupplier{base: time.Now().Unix(), suffix: 0}
}

func (id *IDSupplier) Next() string {
	id.mu.Lock()
	defer id.mu.Unlock()
	id.suffix = id.suffix + 1
	return fmt.Sprintf("%d:%d", id.base, id.suffix)
}
