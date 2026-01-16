package lru

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

type LRUCache struct {
	store    map[string]string
	buffer   *CircularBuffer
	capacity int
}

func NewLruCache(cap int) *LRUCache {
	return &LRUCache{
		store:    make(map[string]string),
		buffer:   NewCircularBuffer(cap),
		capacity: cap,
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	val, ok := c.store[key]
	if !ok {
		Logger.Warn("Cache miss", slog.String("key", key))
		return "", false
	}
	c.buffer.Add(key)
	Logger.Info("Cache hit", slog.String("key", key), slog.String("value", val))
	return val, true
}

func (c *LRUCache) Put(key, value string) {
	if _, ok := c.store[key]; ok {
		c.store[key] = value
		c.buffer.Add(key)
		Logger.Info("Updated key", slog.String("key", key), slog.String("value", value))
		return
	}
	removed := c.buffer.Add(key)
	delete(c.store, removed)
	c.store[key] = value
	if removed != "" {
		Logger.Info("Put with eviction",
			slog.String("key", key),
			slog.String("value", value),
			slog.String("removed", removed))
	} else {
		Logger.Info("Put new key", slog.String("key", key), slog.String("value", value))
	}
}
