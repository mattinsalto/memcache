#Memcache for go

Memcache is a memory cache to store a collection of any kind of struct for a given duration. Cached items are indexed by a unique string (cacheID) and it has the option to sliding expiration. If it's on, the expiration of the cached items will be renewed each time the items are requested. 


**Inizialization:**

    myCache := memcache.New(slidingexp, expcallback)

*slidingexp* bool: indicates if cached items expiration will be renewed each time a cached item is requested.

*expcallback* function( cacheID string): is optional and could be any void function with a string parameter (cacheID). It will be called each time a cache item expires.