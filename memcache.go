/*Package memcache is a memory cache to store a collection of any kind of struct in memory for a given duration.
Cached items are indexed by a unique string.
It has the option to sliding expiration. If it's on, the expiration of the cached items
will be renewed each time the item is requested. */
package memcache

import (
	"fmt"
	"time"
)

//Cache item wrapper is the struct stored in memory cache.
//It holds the struct to be stored plus needed properties for cache item expiration
type cacheitmwrpr struct {
	cacheitmID string
	cacheitm   interface{}
	duration   time.Duration
	exptimer   *time.Timer
	expiration time.Time
}

/*
Memcache is a memory cache which stores cache items indexed by ID
index:  is the map storing references to cache item wraper structs
slidingexp: indicates if the expiration of the cache is going to be renewed
            each time the cache is requested
expcallback (optional): is one or various functions that will be called when a cache item expires
and is removed from the memory cache
*/
type Memcache struct {
	index       map[string]*cacheitmwrpr
	slidingexp  bool
	expcallback []func(string)
}

//New inizializes MemoryCache
func New(slidingExp bool, expCallback ...func(string)) (memcache *Memcache) {

	memcache = new(Memcache)
	memcache.index = make(map[string]*cacheitmwrpr)
	memcache.slidingexp = slidingExp
	memcache.expcallback = expCallback
	return
}

//Set a cache item and its expiration
func (memcache *Memcache) Set(cacheID string, cacheitm interface{}, duration time.Duration) {

	cw := &cacheitmwrpr{
		cacheitmID: cacheID,
		cacheitm:   cacheitm,
		duration:   duration,
		exptimer:   time.NewTimer(duration),
		expiration: time.Now().Add(duration)}

	cw.ttl(memcache)

	memcache.index[cacheID] = cw
}

/*
Get a cache item by cacheID
If sliding expiration of the memory cache is true,
cache item expiration will be renewed
*/
func (memcache *Memcache) Get(cacheID string) (cacheitm interface{}, err error) {

	cw, err := getCacheitmwrpr(cacheID, memcache)
	if err != nil {
		return
	}
	if memcache.slidingexp {
		cw.ttl(memcache)
	}

	return cw.cacheitm, err
}

/*
TTL sets automatic deletion of a cache item when duration expires
This function is only for renewing expiration on non sliding expiration cache items.
When setting a cache item, the expiration is set and in case
memory cache sliding expiration is true, cache item expiration is renewed
automatically every time a cache item is requested. So there is no need
to call this function unless you want to renew the expiration
of non sliding expiration cache items.
*/
func (memcache *Memcache) TTL(cacheID string, duration time.Duration) (err error) {

	cw, err := getCacheitmwrpr(cacheID, memcache)
	if err != nil {
		return
	}
	cw.duration = duration
	cw.ttl(memcache)
	return
}

//Expire a cache item from memory cache
func (memcache *Memcache) Expire(cacheID string) (err error) {

	cw, err := getCacheitmwrpr(cacheID, memcache)
	if err != nil {
		return
	}
	cw.expire(memcache)
	return
}

//Expiration indicates cached item expiration date
func (memcache *Memcache) Expiration(cacheID string) (expdate time.Time, err error) {

	cw, err := getCacheitmwrpr(cacheID, memcache)
	if err != nil {
		return
	}

	return cw.expiration, nil
}

//Gets a cache item wraper by cacheID from memory cache
func getCacheitmwrpr(cacheID string, memcache *Memcache) (cw *cacheitmwrpr, err error) {

	cw, exists := memcache.index[cacheID]
	if !exists {
		err = fmt.Errorf("There is no cache in memory with ID: %s", cacheID)
		return nil, err
	}
	return cw, nil
}

//Sets automatic deletion of a cache item wraper from memory cache when duration expires
func (cw *cacheitmwrpr) ttl(memcache *Memcache) {

	cw.expiration = time.Now().Add(cw.duration)

	go func() {
		cw.exptimer.Reset(cw.duration)
		<-cw.exptimer.C
		cw.expire(memcache)
	}()
}

//Removes cache item wraper from memory cache
func (cw *cacheitmwrpr) expire(memcache *Memcache) {
	delete(memcache.index, cw.cacheitmID)
	if len(memcache.expcallback) > 0 {
		memcache.expcallback[0](cw.cacheitmID)
	}
}
