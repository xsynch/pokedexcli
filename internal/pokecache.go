package pokecache

import (
	"log"
	"sync"
	"time"
)

type PokeCache struct {
	Entry      map[string]CacheEntry
	mutex      *sync.Mutex
	expiration time.Duration
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) PokeCache {
	newCache := PokeCache{expiration: interval, Entry: make(map[string]CacheEntry), mutex: &sync.Mutex{}}

	go newCache.Readloop(interval)

	return newCache
}

func (p *PokeCache) Add(key string, v []byte) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// newCacheEntry := make(map[string]CacheEntry)
	// newCacheEntry[key] = CacheEntry{CreatedAt: time.Now(), Val: v}
	p.Entry[key] = CacheEntry{CreatedAt: time.Now(), Val: v}

}

func (p *PokeCache) Get(key string) (value []byte, b bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, ok := p.Entry[key]
	if ok {
		return p.Entry[key].Val, true
	}

	return nil, false
}

func (p *PokeCache) Readloop(interval time.Duration) {

	ticker := time.NewTicker(interval)
	for range ticker.C {
		p.deleteCacheEntry(interval)
	}

}

func (p *PokeCache) deleteCacheEntry(interval time.Duration) {
	expirationTime := time.Now().Add(-interval)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for idx, val := range p.Entry {

		log.Printf("Removed: %v at %v", idx, expirationTime)
		if val.CreatedAt.Before(expirationTime) {
			delete(p.Entry, idx)

		}
	}
}
