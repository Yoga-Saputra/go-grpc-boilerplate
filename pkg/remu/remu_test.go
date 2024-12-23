package remu

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

func TestMutex(t *testing.T) {

	rda := New(Config{
		Host:       "127.0.0.1",
		Port:       6379,
		Password:   "",
		Database:   1,
		MaxRetries: 3,
		PoolSize:   9,
	})

	// create new pool & redis sync instance
	pool := goredis.NewPool(rda.Client)
	rs := redsync.New(pool)
	_ = rs

	c := &Counter{}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			for j := 0; j < 1000; j++ {
				mutex := rs.NewMutex(fmt.Sprintf("add-account:{%v}", j))
				_ = mutex.Lock()

				defer func() {
					_, _ = mutex.Unlock()
				}()

				c.Increment()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	// Check the value
	if c.Value() != 100000 {
		t.Errorf("Expected 100000, got %d", c.Value())
	}

	log.Println(c.Value())
}

type Counter struct {
	value int
}

func (c *Counter) Increment() {
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}
