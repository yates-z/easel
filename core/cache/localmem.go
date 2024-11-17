package cache

import (
	"sync"
	"time"
)

// MemCacheNode is a doubly linked list
type MemCacheNode[K comparable, V any] struct {
	Key        K
	Value      V
	Expiration int64
	prev       *MemCacheNode[K, V]
	next       *MemCacheNode[K, V]
}

// IsExpired check whether it has expired.
func (node *MemCacheNode[K, V]) IsExpired() bool {
	if node.Expiration <= 0 {
		return false // never expire.
	}
	return time.Now().UnixNano() > node.Expiration
}

type MemCacheShard[K comparable, V any] struct {
	capacity int
	nodes    map[K]*MemCacheNode[K, V]
	mu       sync.Mutex
	head     *MemCacheNode[K, V]
	tail     *MemCacheNode[K, V]
}

// Move a node to the front of the list
func (s *MemCacheShard[K, V]) moveToFront(node *MemCacheNode[K, V]) {
	if node == s.head {
		return
	}

	// Detach the node
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}

	// If it's the tail, update tail pointer
	if node == s.tail {
		s.tail = node.prev
	}

	// Move to the front
	node.next = s.head
	node.prev = nil
	if s.head != nil {
		s.head.prev = node
	}
	s.head = node

	// Update tail if needed
	if s.tail == nil {
		s.tail = node
	}
}

// Remove removes a key from the cache
func (s *MemCacheShard[K, V]) remove(node *MemCacheNode[K, V]) {
	delete(s.nodes, node.Key)

	// Update pointers
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}

	// Update head or tail if needed
	if node == s.head {
		s.head = node.next
	}
	if node == s.tail {
		s.tail = node.prev
	}

	// If list is empty, reset head and tail
	if s.head == nil {
		s.tail = nil
	}
}

// Remove the least recently used node (tail)
func (s *MemCacheShard[K, V]) removeTail() {
	if s.tail == nil {
		return
	}
	s.remove(s.tail)
}

// Set a key-value pair to cache.
func (s *MemCacheShard[K, V]) set(key K, value V, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if node, exists := s.nodes[key]; exists {
		// Update value and expiration
		node.Value = value
		node.Expiration = expirationTime(ttl)
		s.moveToFront(node)
		return nil
	}
	// Add new node
	newNode := &MemCacheNode[K, V]{
		Key:        key,
		Value:      value,
		Expiration: expirationTime(ttl),
	}
	s.nodes[key] = newNode
	s.moveToFront(newNode)

	// If over capacity, remove least recently used item
	if len(s.nodes) > s.capacity {
		s.removeTail()
	}
	return nil
}

// Get a value from given key
func (s *MemCacheShard[K, V]) get(key K) (value V, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	node, exists := s.nodes[key]
	if !exists {
		return value, false
	}
	if node.IsExpired() {
		s.remove(node)
		return value, false
	}
	return node.Value, true
}

// getOrSet never returns error.
func (s *MemCacheShard[K, V]) getOrSet(key K, newVal V, ttl time.Duration) (value V, loaded bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the key exists
	node, exists := s.nodes[key]
	if exists {
		if !node.IsExpired() {
			// Key exists and is valid, move to front and return the value
			s.moveToFront(node)
			return node.Value, true, nil
		}
		s.remove(node)
	}

	newNode := &MemCacheNode[K, V]{
		Key:        key,
		Value:      newVal,
		Expiration: expirationTime(ttl),
	}
	s.nodes[key] = newNode
	s.moveToFront(newNode)
	// Enforce capacity
	if len(s.nodes) > s.capacity {
		s.removeTail()
	}
	return newVal, false, nil
}

func (s *MemCacheShard[K, V]) keys() []K {
	s.mu.Lock()
	defer s.mu.Unlock()

	var keys []K
	for node := s.head; node != nil; node = node.next {
		keys = append(keys, node.Key)
	}
	return keys
}

func (s *MemCacheShard[K, V]) delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if node, exists := s.nodes[key]; exists {
		s.remove(node)
	}
}

// cleanup clean all keys expired.
func (s *MemCacheShard[K, V]) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, node := range s.nodes {
		if node.IsExpired() {
			s.remove(node)
		}
	}
}

func (s *MemCacheShard[K, V]) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, node := range s.nodes {
		s.remove(node)
	}
}

var _ Cache[string, string] = (*MemCache[string, string])(nil)

// MemCache is a simple lru cache.
type MemCache[K comparable, V any] struct {
	shards     []*MemCacheShard[K, V]
	numShards  int
	capacity   int
	cleanupDur time.Duration
	stopCh     chan struct{}
}

func NewMemCache[K comparable, V any](numShards, capacity int, cleanupInterval time.Duration) *MemCache[K, V] {
	if numShards <= 0 || capacity <= 0 {
		panic("numShards and capacity must be greater than 0")
	}
	shards := make([]*MemCacheShard[K, V], numShards)
	baseCapacity := capacity / numShards
	extraCapacity := capacity % numShards // Remaining capacity to distribute
	for i := 0; i < numShards; i++ {
		actualCapacity := baseCapacity
		if i < extraCapacity {
			actualCapacity++ // Distribute the remaining capacity
		}
		shards[i] = &MemCacheShard[K, V]{
			capacity: actualCapacity,
			nodes:    make(map[K]*MemCacheNode[K, V]),
		}
	}
	cache := &MemCache[K, V]{
		shards:     shards,
		numShards:  numShards,
		capacity:   capacity,
		cleanupDur: cleanupInterval,
		stopCh:     make(chan struct{}),
	}
	go cache.startCleanup()
	return cache
}

func (c *MemCache[K, V]) getShard(key K) *MemCacheShard[K, V] {
	hash := fnv32(key)
	return c.shards[hash%uint32(c.numShards)]
}

func (c *MemCache[K, V]) Set(key K, value V, ttl time.Duration) error {
	shard := c.getShard(key)
	return shard.set(key, value, ttl)
}

func (c *MemCache[K, V]) Get(key K) (value V, ok bool) {
	shard := c.getShard(key)
	return shard.get(key)
}

func (c *MemCache[K, V]) GetDefault(key K, _default V) V {
	shard := c.getShard(key)
	value, ok := shard.get(key)
	if !ok {
		return _default
	}
	return value
}

func (c *MemCache[K, V]) GetOrSet(key K, newVal V, ttl time.Duration) (V, bool, error) {
	shard := c.getShard(key)
	return shard.getOrSet(key, newVal, ttl)
}

func (c *MemCache[K, V]) HasKey(key K) bool {
	_, exists := c.Get(key)
	return exists
}

// Keys returns a slice of all keys in the cache.
func (c *MemCache[K, V]) Keys() []K {
	var keys []K
	for _, shard := range c.shards {
		keys = append(keys, shard.keys()...)
	}
	return keys
}

// Delete removes a specific key from the cache.
func (c *MemCache[K, V]) Delete(key K) {
	shard := c.getShard(key)
	shard.delete(key)
}

// startCleanup start to clean up.
func (c *MemCache[K, V]) startCleanup() {
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

func (c *MemCache[K, V]) Clear() {
	for _, shard := range c.shards {
		shard.clear()
	}
}

func (c *MemCache[K, V]) Stop() {
	close(c.stopCh)
}
