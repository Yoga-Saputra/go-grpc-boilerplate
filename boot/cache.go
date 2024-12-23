package boot

import (
	"context"
	"log"

	// rcache "github.com/go-redis/cache/v8"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/config"
	"github.com/Yoga-Saputra/go-grpc-boilerplate/pkg/remu"
	rcache "github.com/go-redis/cache/v9"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredisv9 "github.com/redis/go-redis/v9"
)

// Cache driver pointer value
var Cache *cacheAdapter

type cacheAdapter struct {
	RCln *goredisv9.Client
	RCch *rcache.Cache
	Rs   *redsync.Redsync
}

// Start cache redis connection
func cacheUp(args *AppArgs) {
	// rda := redis.New(redis.Config{
	// 	Host:       config.Of.Cache.Redis.Host,
	// 	Port:       config.Of.Cache.Redis.Port,
	// 	Password:   config.Of.Cache.Redis.Password,
	// 	Database:   config.Of.Cache.Redis.Database,
	// 	MaxRetries: config.Of.Cache.Redis.MaxRetries,
	// 	PoolSize:   config.Of.Cache.Redis.PoolSize,
	// })

	rda := remu.New(remu.Config{
		Host:       config.Of.Cache.Redis.Host,
		Port:       config.Of.Cache.Redis.Port,
		Password:   config.Of.Cache.Redis.Password,
		Database:   config.Of.Cache.Redis.Database,
		MaxRetries: config.Of.Cache.Redis.MaxRetries,
		PoolSize:   config.Of.Cache.Redis.PoolSize,
	})

	// Create new redis-cache instance
	cache := rcache.New(&rcache.Options{
		Redis: rda.Client,
	})

	// create new pool & redis sync instance
	pool := goredis.NewPool(rda.Client)
	rs := redsync.New(pool)

	// Create adapter
	Cache = &cacheAdapter{
		RCln: rda.Client,
		RCch: cache,
		Rs:   rs,
	}
	printOutUp("New Cache Redis connection successfully open")
}

// Stop cache redis connection
func cacheDown() {
	printOutDown("Closing current Cache2 Redis connection...")

	if Cache.RCln != nil {
		id := Cache.RCln.ClientID(context.Background())

		if err := Cache.RCln.Close(); err != nil {
			log.Printf("ERROR - failed to close redis connection, err: %v \n", err.Error())
		}

		log.Printf("SUCCESS - Redis connection already closed, %v \n", id)
	}
}
