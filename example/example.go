/*
Memcache has two kinds of memory caches, both with TTL:

1. Memory cache with sliding expiration:
Cached items will expire on time unless
they are requested. They will renew expiration time by duration
each time they are requested.

2. Non sliding expiration memory cache:
Cached items will expire on time.
If you want to renew expiration time, you must do it manually
calling Memcache.TTL(cacheID string, duration time.Duration)

For this example we will store session structs in the cache. It could be any type.
But memcache is specially useful for sessions in a web app or web api.
We can store the session in the database when its created,
but use the cached one each time we receive a request with sessionID or token,
avoiding unnecessary calls to the database. When cached item expires
(session expires), we can log the logout date in previously saved
session register in the database with the help of the expcallback function
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mattinsalto/memcache"
)

type session struct {
	user      string
	name      string
	lastname  string
	profileID int
	//... any other info we need in session
}

/*
This function will be called when a cache item has expired.
Can be any void function with a string parámeter.
It's optional, you can inizialize Memcache without it
*/
func expcallback(cacheitmID string) {
	fmt.Printf("%s has expired at %s \n", cacheitmID, time.Now())
}

func main() {

	sessionOne := session{
		user:      "mattinsalto",
		name:      "Martin",
		lastname:  "Garmendia",
		profileID: 1}

	sessionTwo := session{
		user:      "gopher",
		name:      "Gopher",
		lastname:  "golang",
		profileID: 5}

	//For this example we will use both types of memcache to show the differences
	slidingExpCache := memcache.New(true, expcallback)
	nonslidingExpCache := memcache.New(false, expcallback)

	/*
		We will set same duration for two sessions, but will store sessionOne in slidingExpCache
		and sessionTwo in nonslidingExpCache. So sessionOne will renew its expiration every time
		is requested and sessionTwo will die. Despite in this example we will store only one struct
		in each memcache, you could store as many as you want.
	*/

	fmt.Println("------- Program started at: ", time.Now())

	//We set two cache items
	err := slidingExpCache.Set("sessionOne", sessionOne, time.Second*10)
	if err != nil {
		println("Error setting the cache item: ", err.Error())
	}

	nonslidingExpCache.Set("sessionTwo", sessionTwo, time.Second*10)

	//Expiration indicates cache item expiration date and time
	expDate, err := slidingExpCache.Expiration("sessionOne")
	if err != nil {
		fmt.Println("sessionOne expiration: ", err.Error())
	} else {
		fmt.Println("sessionOne expiration: ", expDate)
	}

	expDate, err = nonslidingExpCache.Expiration("sessionTwo")
	if err != nil {
		fmt.Println("sessionTwo expiration: ", err)
	} else {
		fmt.Println("sessionTwo expiration: ", expDate)
	}

	go func() {
		t := time.NewTimer(time.Second * 5)
		<-t.C
		fmt.Println("------- 5 seconds elapsed: ", time.Now())

		//We request sessionOne, so its expiration time will be renewed by 10 seconds
		sessionInfo("sessionOne", slidingExpCache)

		//We request sessionTwo, but as sliding expiration is false, it will die in 3 seconds
		sessionInfo("sessionTwo", nonslidingExpCache)
	}()

	go func() {
		t := time.NewTimer(time.Second * 10)
		<-t.C
		fmt.Println("------- 10 seconds elapsed: ", time.Now())
		//sessionOne is alive and we are extending its life for 10 seconds more
		sessionInfo("sessionOne", slidingExpCache)

		//sessionTwo has expired so doesn't exist hence will see an error.
		//We could renew sessionTwo expiration by calling:
		//nonslidingExpCache.TTL("654321", time.Second * 10) before its expiration
		sessionInfo("sessionTwo", nonslidingExpCache)
	}()

	//Wait for user input to terminate
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

//Gets session from memcache and shows its info
func sessionInfo(cacheID string, slidingExpCache *memcache.Memcache) {
	//Get session by cacheID
	cacheitm, err := slidingExpCache.Get(cacheID)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s: ", cacheID), err.Error())
	} else {
		fmt.Println(fmt.Sprintf("%s: ", cacheID), cacheitm.(session))
		expDate, err := slidingExpCache.Expiration(cacheID)
		if err != nil {
			fmt.Println(fmt.Sprintf("%s expiration: ", cacheID), err)
		} else {
			fmt.Println(fmt.Sprintf("%s expiration: ", cacheID), expDate)
		}
	}
}
