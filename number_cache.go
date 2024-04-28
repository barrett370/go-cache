package cache

import (
	"fmt"
	"runtime"
	"time"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type NumberCache[T Number] struct {
	*cache[T]
}

func NewNumberCache[T Number](defaultExpiration, cleanupInterval time.Duration) *NumberCache[T] {
	items := make(map[string]Item[T])
	c := newCache(defaultExpiration, items)
	// This trick ensures that the janitor goroutine (which--granted it
	// was enabled--is running DeleteExpired on c forever) does not keep
	// the returned C object from being garbage collected. When it is
	// garbage collected, the finalizer stops the janitor goroutine, after
	// which c can be collected.
	C := &Cache[T]{c}
	if cleanupInterval > 0 {
		runJanitor(c, cleanupInterval)
		runtime.SetFinalizer(C, stopJanitor[T])
	}
	return &NumberCache[T]{
		cache: c,
	}
}

func (c *NumberCache[T]) Increment(k string, n T) (T, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.mu.Unlock()
		var t T
		return t, fmt.Errorf("Item %s not found", k)
	}
	v.Object += n
	c.items[k] = v
	c.mu.Unlock()
	return v.Object, nil
}

func (c *NumberCache[T]) Decrement(k string, n T) (T, error) {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.mu.Unlock()
		var t T
		return t, fmt.Errorf("Item %s not found", k)
	}
	v.Object -= n
	c.items[k] = v
	c.mu.Unlock()
	return v.Object, nil
}
