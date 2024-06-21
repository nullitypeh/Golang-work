package config

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"github.com/gomodule/redigo/redis"
)

var (
	DBEngine *xorm.Engine
	RedisPool *redis.Pool
)

func InitDB() {
	var err error
	DBEngine, err = xorm.NewEngine("mysql", "user:password@/dbname?charset=utf8")
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}
}

func InitRedis() {
	RedisPool = &redis.Pool{
		MaxIdle:   3,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}
