package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"Golangxy/config"
	"Golangxy/models"
	"Golangxy/utils"

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
func handlePanic(w http.ResponseWriter) {
	if err := recover(); err != nil {
		utils.Logger.Printf("Panic recovered: %v", err)
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
	}
}

func CreateItemWithLock(w http.ResponseWriter, r *http.Request, localCache *utils.Cache) {
	defer handlePanic(w)
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lockKey := "create_item_lock"
	locked, err := utils.AcquireLock(lockKey, 10)
	if err != nil {
		panic("Failed to acquire lock")
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
		panic(err)
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
		"data": map[string]interface{}{
			"item_info": item,
		},
	}
	json.NewEncoder(w).Encode(response)
	utils.Logger.Println("Item created successfully")
}

func UpdateItemWithLock(w http.ResponseWriter, r *http.Request, localCache *utils.Cache) {
	defer handlePanic(w) // 用于处理异常

	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["item_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		utils.Logger.Println("Invalid item_id:", err)
		return
	}

	var itemxy models.Item
	if err := json.NewDecoder(r.Body).Decode(&itemxy); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		utils.Logger.Println("Failed to decode request body:", err)
		return
	}
	upid := itemxy.Item_id
	utils.Logger.Println("id updated successfully:", id)
	utils.Logger.Println("Item_id updated successfully:", itemxy.Item_id)

	lockKey := fmt.Sprintf("update_item_lock_%d", id)
	locked, err := utils.AcquireLock(lockKey, 10)
	if err != nil {
		http.Error(w, "Failed to acquire lock", http.StatusInternalServerError)
		utils.Logger.Println("Failed to acquire lock:", err)
		return
	}
	if !locked {
		http.Error(w, "Resource is locked", http.StatusConflict)
		utils.Logger.Println("Resource is locked:", id)
		return
	}
	defer utils.ReleaseLock(lockKey)

	// 更新数据库中的记录，允许更改 item_id
	_, err = config.DBEngine.ID(upid).Update(&itemxy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		utils.Logger.Println("Failed to update item:", err)
		return
	}

	// 如果 item_id 发生变化，处理缓存
	if id != itemxy.Item_id {
		oldCacheKey := fmt.Sprintf("item_%d", id)
		newCacheKey := fmt.Sprintf("item_%d", itemxy.Item_id)

		// 删除旧缓存
		localCache.Delete(oldCacheKey)
		conn := config.RedisPool.Get()
		defer conn.Close()
		conn.Do("DEL", oldCacheKey)

		// 更新新缓存
		data, _ := json.Marshal(itemxy)
		localCache.Set(newCacheKey, itemxy, 10*time.Minute)
		conn.Do("SET", newCacheKey, data)
	} else {
		cacheKey := fmt.Sprintf("item_%d", id)
		data, _ := json.Marshal(itemxy)
		localCache.Set(cacheKey, itemxy, 10*time.Minute)
		conn := config.RedisPool.Get()
		defer conn.Close()
		conn.Do("SET", cacheKey, data)
	}

	response := map[string]interface{}{
		"code": 0,
		"msg":  "成功",
		"data": map[string]interface{}{
			"store_info": itemxy,
		},
	}
	json.NewEncoder(w).Encode(response)
	utils.Logger.Println("Item updated successfully:", itemxy)
}

func GetItemWithCache(w http.ResponseWriter, r *http.Request, localCache *utils.Cache) {
	defer handlePanic(w)
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
				panic(err)
			}
			if !has {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			// 更新 Redis 缓存
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
	utils.Logger.Println("Item retrieved successfully")
}

func DeleteItem(w http.ResponseWriter, r *http.Request, localCache *utils.Cache) {
	defer handlePanic(w)
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
		panic(err)
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
	utils.Logger.Println("Item deleted successfully")
}
