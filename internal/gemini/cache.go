package gemini

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// Cache 结果缓存
type Cache struct {
	config  CacheConfig
	entries map[string]*cacheEntry
	mu      sync.RWMutex
}

// cacheEntry 缓存条目
type cacheEntry struct {
	value     string
	expiresAt time.Time
}

// NewCache 创建缓存
func NewCache(cfg CacheConfig) *Cache {
	cache := &Cache{
		config:  cfg,
		entries: make(map[string]*cacheEntry),
	}

	// 启动过期清理协程
	go cache.cleanup()

	return cache
}

// Get 获取缓存
func (c *Cache) Get(key string) (string, bool) {
	if !c.config.Enabled {
		return "", false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	hash := c.hashKey(key)
	entry, exists := c.entries[hash]
	if !exists {
		return "", false
	}

	// 检查是否过期
	if time.Now().After(entry.expiresAt) {
		return "", false
	}

	return entry.value, true
}

// Set 设置缓存
func (c *Cache) Set(key, value string) {
	if !c.config.Enabled {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查是否超过最大条目数
	if len(c.entries) >= c.config.MaxEntries {
		c.evictOldest()
	}

	hash := c.hashKey(key)
	c.entries[hash] = &cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.config.TTL),
	}
}

// Delete 删除缓存
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	hash := c.hashKey(key)
	delete(c.entries, hash)
}

// Clear 清空缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*cacheEntry)
}

// Size 返回缓存大小
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.entries)
}

// hashKey 对 key 进行哈希
func (c *Cache) hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// evictOldest 驱逐最旧的条目
func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.entries {
		if oldestKey == "" || entry.expiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.expiresAt
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
	}
}

// cleanup 定期清理过期条目
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.expiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}
