package handlers

import (
	"Golangxy/config"
	"Golangxy/models"
	"Golangxy/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

// 获取时区
func getTimeZone(appLocal string) *time.Location {
	var loc *time.Location
	switch appLocal {
	case "uk":
		loc, _ = time.LoadLocation("Europe/London")
	case "jp":
		loc, _ = time.LoadLocation("Asia/Tokyo")
	case "ru":
		loc, _ = time.LoadLocation("Europe/Moscow")
	default:
		loc = time.UTC
	}
	return loc
}

// @warn WithLock没必要体现在如此上层的方法名字上
func CreateItemWithLock(w http.ResponseWriter, r *http.Request, localCache *utils.Cache) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lockKey := "create_item_lock"
	locked, err := utils.AcquireLock(lockKey, 10)
	if err != nil {
		http.Error(w, "Failed to acquire lock", http.StatusInternalServerError)
		return
	}
	if !locked {
		http.Error(w, "Resource is locked", http.StatusConflict)
		return
	}
	defer utils.ReleaseLock(lockKey)

	// 将创建时间存储为当前时间
	item.CreatedAt = time.Now()

	// 将信息存入数据库
	_, err = config.DBEngine.Insert(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// @warn 缺少错误组件封装
	// 获取请求头中的 app_local 字段
	response := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
		"data": map[string]interface{}{
			"item_info": item,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateItemWithLock(w http.ResponseWriter, r *http.Request, localCache *utils.Cache) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["item_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var itemxy models.Item
	if err := json.NewDecoder(r.Body).Decode(&itemxy); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	itemxy.Item_id = id

	lockKey := fmt.Sprintf("update_item_lock_%d", id)
	locked, err := utils.AcquireLock(lockKey, 10)
	if err != nil {
		http.Error(w, "Failed to acquire lock", http.StatusInternalServerError)
		return
	}
	if !locked {
		http.Error(w, "Resource is locked", http.StatusConflict)
		return
	}
	defer utils.ReleaseLock(lockKey)

	// 将更新时间存储为当前时间
	itemxy.UpdatedAt = time.Now()

	// @warn 最好抽象单独的数据处理层
	_, err = config.DBEngine.ID(id).Update(&itemxy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
		"data": map[string]interface{}{
			"store_info": itemxy,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func GetItemWithCache(w http.ResponseWriter, r *http.Request, localCache *utils.Cache) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["item_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("item_%d", id)
	var item models.Item

	// 检查本地缓存
	if cachedItem, found := localCache.Get(cacheKey); found {
		item = cachedItem.(models.Item)
		//@warn 提前结束而不是else，可以增加接口可读性
	} else {
		// 检查 Redis 缓存
		conn := config.RedisPool.Get()
		defer conn.Close()
		data, err := redis.Bytes(conn.Do("GET", cacheKey))
		if err == nil {
			json.Unmarshal(data, &item)
		} else {
			// 从数据库获取
			has, err := config.DBEngine.ID(id).Get(&item)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !has {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			// 更新 Redis 缓存
			// @warn 丢弃 error不是一个好的例子
			data, _ := json.Marshal(item)
			conn.Do("SET", cacheKey, data)
		}
		// 更新本地缓存
		localCache.Set(cacheKey, item, 10*time.Minute)
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
		"data": map[string]interface{}{
			"store_info": item,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteItem(w http.ResponseWriter, r *http.Request, localCache *utils.Cache) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["item_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 从数据库删除
	var item models.Item
	_, err = config.DBEngine.ID(id).Delete(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 从 Redis 删除缓存
	cacheKey := fmt.Sprintf("item_%d", id)
	conn := config.RedisPool.Get()
	defer conn.Close()
	conn.Do("DEL", cacheKey)

	// 从本地缓存删除
	localCache.Delete(cacheKey)

	// 获取请求头中的 app_local 字段
	appLocal := r.Header.Get("app_local")
	loc := getTimeZone(appLocal)
	deleteTime := time.Now().In(loc).Format("2006-01-02 15:04:05")

	response := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
		"data": map[string]interface{}{
			"delete_time": deleteTime,
		},
	}
	json.NewEncoder(w).Encode(response)
}
