package main

import (
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	cache := newTTLCache[string](time.Minute)
	cache.set("key", "value")
	got, ok := cache.get("key")
	if !ok || got != "value" {
		t.Errorf("get() = %q, %v; want value, true", got, ok)
	}
}

func TestCacheMiss(t *testing.T) {
	cache := newTTLCache[string](time.Minute)
	if _, ok := cache.get("missing"); ok {
		t.Error("get(missing) = true, want false")
	}
}

func TestCacheExpiry(t *testing.T) {
	cache := newTTLCache[string](time.Millisecond)
	cache.set("key", "value")
	time.Sleep(5 * time.Millisecond)
	if _, ok := cache.get("key"); ok {
		t.Error("get(expired) = true, want false")
	}
}

func TestCacheExpiredEntryIsRemoved(t *testing.T) {
	cache := newTTLCache[string](time.Millisecond)
	cache.set("key", "value")
	time.Sleep(5 * time.Millisecond)
	cache.get("key")
	cache.mu.Lock()
	if len(cache.values) != 0 {
		t.Errorf("expired entries not removed, len = %d", len(cache.values))
	}
	cache.mu.Unlock()
}

func TestCacheOverwrite(t *testing.T) {
	cache := newTTLCache[int](time.Minute)
	cache.set("key", 1)
	cache.set("key", 2)
	got, ok := cache.get("key")
	if !ok || got != 2 {
		t.Errorf("get() = %d, %v; want 2, true", got, ok)
	}
}
