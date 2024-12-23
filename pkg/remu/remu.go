package remu

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Remu is Redis interface that is implemented by app
type Remu struct {
	Client *redis.Client

	nameSpace string
}

// New creates a new redis client
func New(config ...Config) *Remu {
	// Set default config
	cfg := configDefault(config...)

	// Create new redis options
	var options *redis.Options
	var err error

	if cfg.URL != "" {
		options, err = redis.ParseURL(cfg.URL)
		if err != nil {
			panic(err)
		}
	} else {
		options = &redis.Options{
			Addr:       fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
			DB:         cfg.Database,
			Username:   cfg.Username,
			Password:   cfg.Password,
			MaxRetries: cfg.MaxRetries,
			PoolSize:   cfg.PoolSize,
		}
	}

	cln := redis.NewClient(options)

	// Test connection
	if err := cln.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &Remu{
		Client:    cln,
		nameSpace: cfg.NameSpace,
	}
}

// Close a redis (remu) connection
func (r *Remu) Close(showLog bool) {
	if r.Client != nil {
		id := r.Client.ClientID(context.Background())

		if err := r.Client.Close(); err != nil {
			if showLog {
				log.Printf("ERROR - failed to close redis (remu) connection, err: %v \n", err.Error())
			}
		}

		if showLog {
			log.Printf("SUCCESS - Redis (remu) connection already closed, %v \n", id)
		}
	}
}

// Get value by key
func (r *Remu) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}

	val, err := r.Client.Get(context.Background(), r.getFullKey(key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}

	return val, err
}

// Set key with value
func (r *Remu) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	return r.Client.Set(context.Background(), r.getFullKey(key), val, exp).Err()
}

// Delete key by key
func (r *Remu) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}

	return r.Client.Del(context.Background(), r.getFullKey(key)).Err()
}

// CheckUniqueIsExists to check is key already exists
// If not, set that key with ttl
func (r *Remu) CheckUniqueIsExists(uniqueKey string, exp time.Duration) (bool, error) {
	// Check given key
	if len(uniqueKey) <= 0 {
		return false, errors.New("unique key of redis (remu) cannot be empty")
	}
	uniqueKey = "uk::" + uniqueKey
	uniqueKey = r.getFullKey(uniqueKey)

	// Do get
	exists, err := r.Get(uniqueKey)
	if err != nil {
		return false, err
	}

	// If not exists, do Set and return false
	if exists == nil {
		loc, err := time.LoadLocation("Asia/Manila")
		if err != nil {
			loc = time.UTC
		}
		now := time.Now().In(loc).String()
		if err := r.Set(uniqueKey, []byte(now), exp); err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// GetScan given key pattern and return key:value pair
func (r *Remu) GetScan(pattern string) (result map[string]interface{}, err error) {
	// Create context and variable
	ctx := context.Background()
	res := make(map[string]interface{})
	var (
		cursor uint64
		keys   []string
		values []interface{}
		errs   []error
	)

	// Iterate based on cursor
	for {
		var ks []string
		var err error
		ks, cursor, err = r.Client.Scan(ctx, cursor, pattern, 0).Result()
		if err != nil {
			errs = append(errs, err)
		}
		keys = append(keys, ks...)

		if cursor == 0 {
			break
		}
	}

	// If got an error while do scan
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	// Do MGet to retrive all values of the keys
	values, err = r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return
	}

	// Pairing keys and value together, happily ever after :*
	for i, v := range keys {
		res[v] = values[i]
	}
	result = res

	return
}
