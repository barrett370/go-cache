// backedcache wraps github.com/barrett370/go-cache#Cache[V] with a backend,
// requests for keys not in the underlying cache will fetch
// (via a singleflight group to prevent stampeding) via a user-defined GetterFn[K,V]
package backedcache
