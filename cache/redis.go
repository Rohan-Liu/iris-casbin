package cache

import (
	"../configs"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var redisPool *redis.Pool

// 执行操作
func do(cmd string, args ...interface{}) (r string, err error) {
	conn := getRedis()
	defer conn.Close()
	r, err = redis.String(conn.Do(cmd, args...))
	return
}

// 设置缓存
func Set(key, val string, ttl time.Duration) (string, error) {
	return do("SET", key, val, "EX", uint64(ttl))
}

// 获取缓存
func Get(key string) string {
	if r, err := do("GET", key); err != nil {
		return ""
	} else {
		return r
	}
}

// 删除
func Del(key string) (string, error) {
	return do("DEL", key)
}

// PING
func Ping() (string, error) {
	return do("PING")
}

// 初始化redis
func init() {
	redisPool = &redis.Pool{
		MaxIdle:     configs.GetConfigInt("redis.maxIdle"),
		MaxActive:   configs.GetConfigInt("redis.maxOpen"),
		IdleTimeout: time.Duration(configs.GetConfigInt("redis.timeout")) * time.Minute,
		Dial: func() (redis.Conn, error) {
			//return redis.Dial("tcp", "127.0.0.1:6379", redis.DialDatabase(0))
			return redis.Dial("tcp", configs.GetConfigString("redis.host")+":"+configs.GetConfigString("redis.port"), redis.DialDatabase(0))
		},
	}

	// 测试连接
	if r, _ := Ping(); r != "PONG" {
		err := errors.New("redis connect failed")
		fmt.Print(err)
	}

}

// 获取redis连接
func getRedis() redis.Conn {
	return redisPool.Get()
}
