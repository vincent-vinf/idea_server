package redisdb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"idea_server/util"
	"sync"
	"time"
)

const (
	codeExpiration = time.Minute * 10 // 验证码有效期
	codeForbidden  = time.Minute      // 验证码重复发送时间
)

var (
	rdb  *redis.Client
	once sync.Once
)

func getInstance() *redis.Client {
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

func Close() {
	if rdb != nil {
		_ = rdb.Close()
	}
}

func InsertEmailCode(email, code, ip string) {
	rdb := getInstance()
	ctx := context.Background()
	// 插入验证码和请求的ip，防止过多请求
	rdb.Set(ctx, email, code, codeExpiration)
	rdb.Set(ctx, ip, ip, codeForbidden)
}

func IsAvailableEmailCode(email, code string) bool {
	rdb := getInstance()
	ctx := context.Background()
	re, err := rdb.Get(ctx, email).Result()
	if err != nil {
		return false
	}
	if re != code {
		return false
	}
	return true
}

func IsAllowedIP(ip string) bool {
	rdb := getInstance()
	ctx := context.Background()
	_, err := rdb.Get(ctx, ip).Result()
	if err == redis.Nil {
		return true
	} else {
		return false
	}
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
