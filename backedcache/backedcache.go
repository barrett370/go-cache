package backedcache

import (
	"context"
	"time"

	"github.com/barrett370/go-cache/v2"
	"tailscale.com/util/singleflight"
)

type CacheKey interface {
	comparable
	String() string
	Expiration() time.Duration
}

type GetterFn[K CacheKey, V any] func(ctx context.Context, key K) (V, error)

type Getter[K CacheKey, V any] interface {
	Get(context.Context, K) (V, error)
}

type Cacher[K CacheKey, V any] interface {
	Getter[K, V]
}

type Cache[K CacheKey, V any] struct {
	cache    cache.Cache[V]
	getterFn GetterFn[K, V]
	sf       singleflight.Group[K, V]
}

func (c *Cache[K, V]) Get(ctx context.Context, key K) (V, error) {
	keyString := key.String()
	if v, ok := c.cache.Get(keyString); ok {
		return v, nil
	}

	v, err, _ := c.sf.Do(key, func() (V, error) {
		v, err := c.getterFn(ctx, key)
		c.cache.Set(keyString, v, key.Expiration())
		return v, err
	})

	if err != nil {
		var z V
		return z, err
	}

	return v, nil
}
