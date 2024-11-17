package cache

import (
	"fmt"
	"testing"
	"time"
)

func Test_EncodeAny(t *testing.T) {
	data, err := EncodeToHex("hello")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data, err)

	origin, err := DecodeFromHex[string](data)
	fmt.Println(origin, err)
}

func Test_FileCache(t *testing.T) {
	cache, _ := NewFileCache[string, string](2, "./cache", 10*time.Second)
	_ = cache.Set("key1", "value1", 30*time.Second)
	value, ok := cache.Get("key1")
	if !ok {
		fmt.Println("key1 not exist")
		return
	}
	fmt.Println("Found key1: " + value)
}

func Test_FileCache_Get(t *testing.T) {
	cache, _ := NewFileCache[string, string](2, "./cache", 10*time.Second)
	value, ok := cache.Get("key1")
	if !ok {
		fmt.Println("key1 not exist")
		return
	}
	fmt.Println("Found key1: " + value)
}

func Test_FileCache_Keys(t *testing.T) {
	cache, _ := NewFileCache[string, string](2, "./cache", 10*time.Second)
	fmt.Println(cache.Keys())
	fmt.Println(cache.HasKey("key1"))
}

func Test_FileCache_Cleanup(t *testing.T) {
	cache, _ := NewFileCache[string, string](2, "./cache", 10*time.Second)
	_ = cache.Set("key1", "value1", 5*time.Second)
	time.Sleep(6 * time.Second)
	value, ok := cache.Get("key1")
	if !ok {
		fmt.Println("key1 not exist")
		return
	}
	fmt.Println("Found key1: " + value)
}

func Test_FileCache_Delete(t *testing.T) {
	cache, _ := NewFileCache[string, string](2, "./cache", 10*time.Second)
	cache.Delete("key1")
	fmt.Println(cache.HasKey("key1"))
}

func Test_FileCache_GetOrSet(t *testing.T) {
	cache, _ := NewFileCache[string, string](2, "./cache", 10*time.Second)
	fmt.Println(cache.GetOrSet("key1", "value1", 30*time.Second))
	fmt.Println(cache.HasKey("key1"))
}
