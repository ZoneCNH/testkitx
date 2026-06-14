package fake

import (
	"fmt"
	"sync"
)

// Reader mirrors the configx.Reader interface consumed by FoundationX modules.
// It is defined locally so testkitx can implement fakes without depending on
// the full configx module.
type Reader interface {
	Get(key string) any
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

// configImpl is a deterministic, in-memory Reader backed by a map.
type configImpl struct {
	mu     sync.RWMutex
	values map[string]any
}

// Compile-time contract: *configImpl implements Reader.
var _ Reader = (*configImpl)(nil)

// FakeConfig creates a deterministic configx.Reader from the supplied values.
// Pass nil for an empty config (all Get calls return zero values).
func FakeConfig(values map[string]any) Reader {
	cp := make(map[string]any, len(values))
	for k, v := range values {
		cp[k] = v
	}
	return &configImpl{values: cp}
}

func (c *configImpl) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.values == nil {
		return nil
	}
	return c.values[key]
}

func (c *configImpl) GetString(key string) string {
	v := c.Get(key)
	s, _ := v.(string)
	return s
}

func (c *configImpl) GetInt(key string) int {
	v := c.Get(key)
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	case string:
		var i int
		_, _ = fmt.Sscanf(n, "%d", &i)
		return i
	}
	return 0
}

func (c *configImpl) GetBool(key string) bool {
	v := c.Get(key)
	b, _ := v.(bool)
	return b
}
