//////////////////////////////////////////////////////////////////////
//
// Given is some code to cache key-value pairs from a database into
// the main memory (to reduce access time). Note that golang's map are
// not entirely thread safe. Multiple readers are fine, but multiple
// writers are not. Change the code to make this thread safe.
//

package main

import (
	"container/list"
	"sync"
)

// CacheSize determines how big the cache can grow
const CacheSize = 100

// KeyStoreCacheLoader is an interface for the KeyStoreCache
type KeyStoreCacheLoader interface {
	// Load implements a function where the cache should gets it's content from
	Load(string) string
}

// KeyStoreCache is a LRU cache for string key-value pairs
type KeyStoreCache struct {
	cache sync.Map
	pages list.List
	load  func(string) string
	mutex sync.Mutex
}

// New creates a new KeyStoreCache
func New(load KeyStoreCacheLoader) *KeyStoreCache {
	return &KeyStoreCache{
		load:  load.Load,
		// Make cache a concurrent map with sync.Map
		cache: sync.Map{},
	}
}

func (k *KeyStoreCache) Len() int {
	var length int
	k.cache.Range(func(_,_ interface{}) bool {
		length++
		return true
	})
	return length
}

// Get gets the key from cache, loads it from the source if needed
func (k *KeyStoreCache) Get(key string) string {

	val , ok := k.cache.LoadOrStore(key, k.load(key))

	// Miss - load from database and save it in cache
	if !ok {
		// Make Pushing pages to front in the cache atomic
		k.mutex.Lock()
		k.pages.PushFront(key)
		k.mutex.Unlock()
		// if cache is full remove the least used item
		if k.Len() > CacheSize {
			k.cache.Delete(k.pages.Back().Value.(string))
			k.pages.Remove(k.pages.Back())
		}
	}
	return val.(string)
}

// Loader implements KeyStoreLoader
type Loader struct {
	DB *MockDB
}

// Load gets the data from the database
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}

	return val
}

func run() *KeyStoreCache {
	loader := Loader{
		DB: GetMockDB(),
	}
	cache := New(&loader)

	RunMockServer(cache)

	return cache
}

func main() {
	run()
}
