package redisdb

import (
	"github.com/go-redis/redis/v8"
	"idea_server/util"
	"sync"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func GetInstance() *redis.Client {
	if rdb == nil {
		once.Do(func() {
			cfg := util.LoadRedisCfg()
			rdb = redis.NewClient(&redis.Options{
				Addr:     cfg.Address,
				Password: cfg.Passwd,
				DB:       cfg.DB,
			})
		})
	}
	return rdb
}

//func ExampleClient() {

//
//	val, err := rdb.Get(ctx, "key").Result()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("key", val)
//
//	val2, err := rdb.Get(ctx, "key2").Result()
//	if err == redis.Nil {
//		fmt.Println("key2 does not exist")
//	} else if err != nil {
//		panic(err)
//	} else {
//		fmt.Println("key2", val2)
//	}
//	// Output: key value
//	// key2 does not exist
//}
