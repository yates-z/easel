package cache

import (
	"fmt"
	"testing"
	"time"
)

func Test_GetSet(t *testing.T) {
	// 创建一个缓存实例，每 10 秒清理一次
	cache := *NewMemCache[string, any](1, 3, 10*time.Second)
	defer cache.Stop()

	// 设置键值对，TTL 为 5 秒
	cache.Set("foo", "bar", 5*time.Second)
	cache.Set("baz", 42, 0) // 永不过期的键值对

	// 获取键值
	if value, found := cache.Get("foo"); found {
		println("Found foo:", value.(string))
	} else {
		println("Foo not found")
	}

	// 等待 6 秒后尝试获取 foo
	time.Sleep(6 * time.Second)
	if value, found := cache.Get("foo"); found {
		println("Found foo:", value.(string))
	} else {
		println("Foo not found after expiration")
	}

	// 删除键
	cache.Delete("baz")
	if _, found := cache.Get("baz"); !found {
		println("Baz not found after delete")
	}
}

func Test_LRU(t *testing.T) {
	// 创建一个缓存实例，每 10 秒清理一次
	cache := *NewMemCache[string, string](2, 3, 10*time.Second)
	defer cache.Stop()

	// 设置键值对，TTL 为 5 秒
	cache.Set("foo", "bar", 0)
	cache.Set("fooo", "barr", 0)
	cache.Set("foooo", "barrr", 0)
	cache.Set("fooooo", "barrrr", 0)

	// 获取键值
	if value, found := cache.Get("foo"); found {
		println("Found foo:", value)
	} else {
		println("Foo not found")
	}
	// 获取键值
	if value, found := cache.Get("fooo"); found {
		println("Found fooo:", value)
	} else {
		println("Fooo not found")
	}
	fmt.Printf("%+v", cache.Keys())
}
