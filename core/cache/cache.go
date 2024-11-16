package cache

import "time"

type Cache[K comparable, V any] interface {
	// Set a value in the cache if the key does not already exist. If
	// timeout is given, use that timeout for the key;
	Set(key K, value V, ttl time.Duration)

	// Get a given key from the cache.
	Get(key K) (value V, ok bool)

	// GetDefault a given key from the cache. If the key does not exist, return
	// default, which itself defaults to nil.
	GetDefault(key K, _default V) V

	// GetOrSet fetch a given key from the cache. If the key does not exist,
	// add the key and set it to the new value. If ttl is given, use that
	// timeout for the key.
	// The value result is the value for the existing key or the given new value.
	// The loaded result is true if the value was loaded, false if stored.
	GetOrSet(key K, newVal V, ttl time.Duration) (value V, loaded bool)

	// HasKey returns True if the key is in the cache and has not expired.
	HasKey(key K) bool

	// Keys returns a slice of all keys in the cache in order of usage.
	// note: Not all keys are exist because they may have expired and not yet been cleaned up.
	Keys() []K

	// Delete a key from the cache and return whether it succeeded, failing
	// silently.
	Delete(key K)

	// Clear *all* values from the cache at once.
	Clear()
}
