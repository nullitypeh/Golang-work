package utils

import (
	/*
	"time"
	*/
	"github.com/gomodule/redigo/redis"
	"Golangxy/config"
)

func AcquireLock(key string, timeout int) (bool, error) {
	conn := config.RedisPool.Get()
	defer conn.Close()
	
	result, err := redis.String(conn.Do("SET", key, "locked", "NX", "EX", timeout))
	if err != nil {
		return false, err
	}
	return result == "OK", nil
}

func ReleaseLock(key string) error {
	conn := config.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}
