#Memcache for go (golang)

Memcache is a memory cache to store a collection of any kind of struct for a given duration (TTL). Cached items are indexed by a unique string (cacheID) and it has the option for sliding expiration. If it's on, the expiration of the cached items will be renewed each time the items are requested. 


**Inizialization:**

    package main
    
    import(
	    "fmt"
	    "time"
	    "github.com/mattinsalto/memcahe"
    )

    function main(){
    	    myCache := memcache.New(slidingexp, expcallback)
    	}
	
	func expcallback(cacheID string) {
		fmt.Printf("Cache item with ID: %s has expired at %s \n", cacheID, time.Now())
	}

 - *slidingexp* bool: indicates if cached items expiration will be renewed each time a cached item is requested.
 
 - *expcallback* function( cacheID string): is optional and could be any void function with a string parameter (cacheID). It will be called
   each time a cache item expires.

**Set a cache item:**

    myCache.Set(cacheID, cacheitm, duration)

 - *cacheID* string: uniquely identifies a cached item

 - *cacheitm* interface{}: any kind of struct you want to store in memcache.

 - *duration* time.Duration: duration for cache item expiration
 
**Get a cache item:**

    myCacheitm := myCache.Get(cacheID).(myCacheitmType)

- *cacheID* string: cached item ID

Type assertion is needed to convert empty interface in the cached item type.


**Renew cache item expiration manually:**

This function is for renewing cached item expiration when sliding expiration is off.

    myCacheitm.TTL(cacheID, duration)

 - *cacheID* string: cached item ID
 - *duration* time.Duration: duration for cache item expiration

**Expire a cache item:**
   Expires a cache item immediately.

    myCacheitm.Expire(cacheID)

 - *cacheID* string: cached item ID

**Get expiration datetime of a cached item**

    expDate := myCacheitm.Expiration(cacheID)

 - *cacheID* string: cached item ID

 - *expDate* time.Time: expiration date of a cached item


----------
There is a working example in:
[example.go](https://github.com/mattinsalto/memcache/blob/master/example/example.go)