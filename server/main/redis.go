package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

/*
	连接池
*/

// 定义全局pool
var pool *redis.Pool

// 初始化连接池
func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:         maxIdle,	// 最大空闲连接数
		MaxActive:       maxActive,	// 与数据库最大连接数，0为没有限制
		IdleTimeout:     idleTimeout, // 最大空闲时间
		Dial: func() (redis.Conn, error) {	// 初始化连接的代码，连接哪个ip的redis
			return redis.Dial("tcp", address)

		},
		//TestOnBorrow:    nil,
		//Wait:            false,
		//MaxConnLifetime: 0,
	}
}
