package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileNode[V any] struct {
	Value      V
	Expiration int64
}

// IsExpired check whether it has expired.
func (node *FileNode[V]) IsExpired() bool {
	if node.Expiration <= 0 {
		return false // never expire.
	}
	return time.Now().UnixNano() > node.Expiration
}

func encodeFileNode[V any](node *FileNode[V]) (compressed []byte, err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err = encoder.Encode(node); err != nil {
		return compressed, fmt.Errorf("failed to encode cache item: %w", err)
	}
	compressed, err = compress(buf.Bytes())
	if err != nil {
		return compressed, fmt.Errorf("failed to compress cache data: %w", err)
	}
	return compressed, nil
}

func decodeFileNode[V any](data []byte) (node *FileNode[V], err error) {
	decompressed, err := decompress(data)
	if err != nil {
		return node, fmt.Errorf("failed to decompress cache data: %w", err)
	}

	decoder := gob.NewDecoder(bytes.NewReader(decompressed))
	if err = decoder.Decode(&node); err != nil {
		return node, fmt.Errorf("failed to decode cache node: %w", err)
	}
	return node, nil
}

// FileCacheShard represents a single shard of the cache
type FileCacheShard[K comparable, V any] struct {
	path       string
	mu         sync.Mutex
	cleanupDur time.Duration
}

// filePath generates a safe file path for a given key
func (s *FileCacheShard[K, V]) filePath(key K) (string, error) {
	hexKey, err := EncodeToHex(key)
	if err != nil {
		return "", fmt.Errorf("failed to encode cache key: %w", err)
	}
	return filepath.Join(s.path, fmt.Sprintf("%s.cache", hexKey)), nil
}

func (s *FileCacheShard[K, V]) getFiles() (files []os.DirEntry) {
	files, _ = os.ReadDir(s.path)
	return
}

func (s *FileCacheShard[K, V]) set(key K, value V, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath, err := s.filePath(key)
	if err != nil {
		return err
	}
	node := FileNode[V]{Expiration: expirationTime(ttl), Value: value}
	compressed, err := encodeFileNode[V](&node)
	if err != nil {
		return err
	}
	if err = os.WriteFile(filePath, compressed, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}
	return nil
}

// Get retrieves a value from the file-based cache
func (s *FileCacheShard[K, V]) get(key K) (value V, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath, err := s.filePath(key)
	if err != nil {
		return value, false
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return value, false
	}
	node, err := decodeFileNode[V](data)
	if err != nil {
		return value, false
	}

	// Check expiration
	if node.IsExpired() {
		_ = os.Remove(filePath) // Clean up expired cache file
		return value, false
	}
	return node.Value, true
}

// getOrSet attempts to get a value for the given key.
// If the key does not exist or is expired, it sets the value using the provided function and returns it.
func (s *FileCacheShard[K, V]) getOrSet(key K, newVal V, ttl time.Duration) (value V, loaded bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath, err := s.filePath(key)

	// Try to read and return existing value
	data, err := os.ReadFile(filePath)
	if err == nil {
		if node, err := decodeFileNode[V](data); err == nil {
			// Check expiration
			if !node.IsExpired() {
				return node.Value, true, nil
			}
		}
	}

	// If the key doesn't exist or is expired, compute a new value
	newNode := FileNode[V]{Expiration: expirationTime(ttl), Value: newVal}
	compressed, err := encodeFileNode[V](&newNode)
	if err != nil {
		return value, false, err
	}

	if err = os.WriteFile(filePath, compressed, 0644); err != nil {
		return value, false, err
	}

	return newVal, false, nil
}

func (s *FileCacheShard[K, V]) keys() []K {
	s.mu.Lock()
	defer s.mu.Unlock()

	var keys []K
	for _, file := range s.getFiles() {
		name := strings.TrimSuffix(file.Name(), ".cache")
		b, err := DecodeFromHex[K](name)
		if err != nil {
			continue
		}
		keys = append(keys, b)
	}

	return keys
}

func (s *FileCacheShard[K, V]) delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if filePath, err := s.filePath(key); err == nil {
		_ = os.Remove(filePath)
	}
}

// cleanup removes expired cache files from all shards
func (s *FileCacheShard[K, V]) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, file := range s.getFiles() {
		filePath := filepath.Join(s.path, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		node, err := decodeFileNode[V](data)
		if err != nil {
			continue
		}

		if node.IsExpired() {
			_ = os.Remove(filePath)
		}
	}
}

func (s *FileCacheShard[K, V]) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, file := range s.getFiles() {
		filePath := filepath.Join(s.path, file.Name())
		_ = os.Remove(filePath)
	}
}

var _ Cache[string, string] = (*FileCache[string, string])(nil)

// FileCache represents the file-based shard cache system
type FileCache[K comparable, V any] struct {
	numShards  int
	shards     []*FileCacheShard[K, V]
	cleanupDur time.Duration
	stopCh     chan struct{}
}

// NewFileCache initializes a new sharded file-based cache.
func NewFileCache[K comparable, V any](numShards int, baseDir string, cleanupInterval time.Duration) (*FileCache[K, V], error) {
	if numShards <= 0 {
		return nil, errors.New("number of shards must be greater than 0")
	}

	cache := &FileCache[K, V]{
		numShards:  numShards,
		shards:     make([]*FileCacheShard[K, V], numShards),
		cleanupDur: cleanupInterval,
		stopCh:     make(chan struct{}),
	}

	// Create shard directories
	for i := 0; i < numShards; i++ {
		shardDir := filepath.Join(baseDir, fmt.Sprintf("shard-%d", i))
		if err := os.MkdirAll(shardDir, 0755); err != nil {
			return nil, err
		}
		cache.shards[i] = &FileCacheShard[K, V]{
			path:       shardDir,
			cleanupDur: cleanupInterval,
		}
	}

	// Start cleanup routine
	go cache.startCleanup()

	return cache, nil
}

// getShard returns the shard for a given key
func (c *FileCache[K, V]) getShard(key K) *FileCacheShard[K, V] {
	shardIdx := fnv32(key) % uint32(c.numShards)
	return c.shards[shardIdx]
}

func (c *FileCache[K, V]) Set(key K, value V, ttl time.Duration) error {
	shard := c.getShard(key)
	return shard.set(key, value, ttl)
}

// Get retrieves a value from the file-based cache
func (c *FileCache[K, V]) Get(key K) (value V, ok bool) {
	shard := c.getShard(key)
	return shard.get(key)
}

func (c *FileCache[K, V]) GetDefault(key K, _default V) V {
	shard := c.getShard(key)
	value, ok := shard.get(key)
	if !ok {
		return _default
	}
	return value
}

func (c *FileCache[K, V]) GetOrSet(key K, newVal V, ttl time.Duration) (V, bool, error) {
	shard := c.getShard(key)
	return shard.getOrSet(key, newVal, ttl)
}

func (c *FileCache[K, V]) HasKey(key K) bool {
	_, exists := c.Get(key)
	return exists
}

func (c *FileCache[K, V]) Keys() []K {
	var keys []K
	for _, shard := range c.shards {
		keys = append(keys, shard.keys()...)
	}
	return keys
}

func (c *FileCache[K, V]) Delete(key K) {
	shard := c.getShard(key)
	shard.delete(key)
}

func (c *FileCache[K, V]) Clear() {
	for _, shard := range c.shards {
		shard.clear()
	}
}

// startCleanup starts a background cleanup goroutine
func (c *FileCache[K, V]) startCleanup() {
	ticker := time.NewTicker(c.cleanupDur)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, shard := range c.shards {
				shard.cleanup()
			}
		case <-c.stopCh:
			return
		}
	}
}

func (c *FileCache[K, V]) Stop() {
	close(c.stopCh)
}
