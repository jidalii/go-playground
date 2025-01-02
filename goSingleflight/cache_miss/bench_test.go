package cache_miss_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
    // "time"
    "log"
    "context"
    
    "golang.org/x/sync/singleflight"

    "goSingleflight/cache_miss"
)


func benchmarkQueryUserGenerate(b *testing.B, num float64) {
	var callCount atomic.Int64
    var cache cache_miss.Cache
    cache.SetupRedis() 
    cache.Client.FlushAll(context.Background())

	getUser := func(id string) (interface{}, error) {
		key := fmt.Sprintf("user:%s", id)
        data, err := cache.GetUser(key)

		if err == nil && len(data) > 0 {
            return data, nil
		}

        callCount.Add(1)
		user, _ := cache_miss.MockDBUserQuery(key)
        err = cache.SetUser(key, user)
        if err != nil {
            panic(err)
        }

		return user, nil
	}

	wg := sync.WaitGroup{}

	b.ResetTimer()
	for i := 0; i < int(num); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getUser("1")
		}()
	}
    
	wg.Wait()
	b.StopTimer()

    cache.Client.FlushAll(context.Background())
	log.Printf("\nDB query was called %d times\n", callCount.Load())
}

func benchmarkQueryUserSingleflightGenerate(b *testing.B, num float64) {
	var callCount atomic.Int64
	g := singleflight.Group{}

    var cache cache_miss.Cache
    cache.SetupRedis() 
    cache.Client.FlushAll(context.Background())

	getUserSingleFlight := func(g *singleflight.Group, id string) (interface{}, error) {
		key := fmt.Sprintf("user:%s", id)
        data, err := cache.GetUser(key)

		if err == nil && len(data) > 0 {
            return data, nil
		}
        
		v, _, _ := g.Do(key, func() (interface{}, error) {
			callCount.Add(1)
			return cache_miss.MockDBUserQuery(key)
		})
        err = cache.SetUser(key, v.(string))
        if err != nil {
            panic(err)
        }
		return v, nil
	}

	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < int(num); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getUserSingleFlight(&g, "1")
		}()
	}
	wg.Wait()
	b.StopTimer()
    cache.Client.FlushAll(context.Background())
	log.Printf("\nDB query was called %d times\n", callCount.Load())
}

// func BenchmarkQueryUserSingleflight_1000(b *testing.B) {
//     // time.Sleep(500 * time.Millisecond)
// 	benchmarkQueryUserSingleflightGenerate(b, 1000)
// }

// func BenchmarkQueryUser_1000(b *testing.B) {
//     // time.Sleep(500 * time.Millisecond)
// 	benchmarkQueryUserGenerate(b, 1000)
// }

func BenchmarkQueryUserSingleflight_10000(b *testing.B) {
    // time.Sleep(1 * time.Second)
	benchmarkQueryUserSingleflightGenerate(b, 10000)
}

func BenchmarkQueryUser_10000(b *testing.B) {
    // time.Sleep(1 * time.Second)
	benchmarkQueryUserGenerate(b, 10000)
}

