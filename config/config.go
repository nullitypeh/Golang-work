package config

import (
	"fmt"
	//"database/sql"
	//"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"xorm.io/xorm"
)

var (
	DBEngine  *xorm.Engine
	RedisPool *redis.Pool
)

func InitDB() {
	var err error
	// Connect to MySQL database
	DBEngine, err = xorm.NewEngine("mysql", "root:201801@tcp(localhost:3306)/Gotest?charset=utf8")
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
