package config

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" 
)

var RedisPool *redis.Pool
var db *sql.DB //连接池对象
func initDB() (err error) {
	//数据库
	//用户名:密码啊@tcp(ip:端口)/数据库的名字
	dsn := "root:123@tcp(127.0.0.1:3306)/test"
	//连接数据集
	db, err = sql.Open("mysql", dsn) 
	if err != nil {
		return
	}
	err = db.Ping() //尝试连接数据库
	if err != nil {
		return
	}
	fmt.Println("连接数据库成功~")
	//设置数据库连接池的最大连接数
	db.SetMaxIdleConns(10)
	return

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
