/*Package memcache is a memory cache to store a collection of any kind of struct for a given duration.
Cached items are indexed by a unique string and it has the option to sliding expiration.
If it's on, the expiration of the cached items will be renewed each time the items are requested. */
package memcache

import (
	"fmt"
	"sync"
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
Memcache memory cache.
index:  is the map storing references to cache items
slidingexp: indicates if the expiration of the cache is going to be renewed
            each time the cache is requested
expcallback (optional): is one or various functions that will be called when a cache item expires
*/
type Memcache struct {
	sync.Mutex
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
func (memcache *Memcache) Set(cacheID string, cacheitm interface{}, duration time.Duration) (err error) {

	memcache.Lock()
	defer memcache.Unlock()

	//Verify cacheID is not already set
	_, err = getCacheitmwrpr(cacheID, memcache)
	if err == nil {
		return fmt.Errorf("The cacheID: %s already exists", cacheID)
	}

	cw := &cacheitmwrpr{
		cacheitmID: cacheID,
		cacheitm:   cacheitm,
		duration:   duration,
		exptimer:   time.NewTimer(duration),
		expiration: time.Now().Add(duration)}

	cw.ttl(memcache)

	memcache.index[cacheID] = cw
	return nil
}

/*
Get a cache item by cacheID
If sliding expiration is true,
expiration will be renewed
*/
func (memcache *Memcache) Get(cacheID string) (cacheitm interface{}, err error) {
	memcache.Lock()
	defer memcache.Unlock()

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
TTL sets automatic deletion of a cache item when duration expires.
This function is only for renewing expiration when sliding expiration is off.
When setting a cache item, the expiration is set and in case
memory cache sliding expiration is true, cache item expiration is renewed
automatically every time a item is requested. So there is no need
to call this function unless you want to renew the expiration
of non sliding expiration items.
*/
func (memcache *Memcache) TTL(cacheID string, duration time.Duration) (err error) {
	memcache.Lock()
	defer memcache.Unlock()

	cw, err := getCacheitmwrpr(cacheID, memcache)
	if err != nil {
		return
	}
	cw.duration = duration
	cw.ttl(memcache)
	return
}

//Expire a cache item
func (memcache *Memcache) Expire(cacheID string) (err error) {

	memcache.Lock()
	defer memcache.Unlock()

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

//Gets a cache item wraper by cacheID
func getCacheitmwrpr(cacheID string, memcache *Memcache) (cw *cacheitmwrpr, err error) {

	cw, exists := memcache.index[cacheID]
	if !exists {
		err = fmt.Errorf("There is no cache in memory with ID: %s", cacheID)
		return nil, err
	}
	return cw, nil
}

//Sets automatic deletion of a cache item wraper when duration expires
func (cw *cacheitmwrpr) ttl(memcache *Memcache) {

	cw.expiration = time.Now().Add(cw.duration)

	go func() {
		cw.exptimer.Reset(cw.duration)
		<-cw.exptimer.C
		cw.expire(memcache)
	}()
}

//Expires a cahe item wraper
func (cw *cacheitmwrpr) expire(memcache *Memcache) {
	delete(memcache.index, cw.cacheitmID)
	if len(memcache.expcallback) > 0 {
		memcache.expcallback[0](cw.cacheitmID)
	}
}
