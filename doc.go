// This package provide a simple way to use cache drivers.
//
// Example Usage
//
// The following is a simple example using memcached driver:
//    import (
//      "fmt"
//      "github.com/fabiorphp/cachego"
//      "github.com/bradfitz/gomemcache/memcache"
//    )
//
//    func main() {
//
//      cache := cachego.NewMemcached(
//          memcached.New("localhost:11211"),
//      )
//
//      cache.Save("foo", "bar")
//
//      fmt.Println(cache.Fetch("foo"))
//    }
package cachego
