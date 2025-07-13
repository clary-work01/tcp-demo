package main

import (
	"time"

	"github.com/go-redis/redis"
)

var rdb *redis.Client


func initRedis(address string, minIdle, maxActive int, idleTimeout time.Duration) {
	rdb = redis.NewClient(&redis.Options{
		Addr:         address,      // Redis 地址，例如 "localhost:6379"
		Password:     "",           // 密碼，如果沒有設置密碼則為空字符串
		DB:           0,            // 默認數據庫
		
		// 連接池配置
		PoolSize:        maxActive,    // 連接池大小 (對應原來的 MaxActive)
		MinIdleConns:    minIdle,      // 最小空閒連接數 (對應原來的 MaxIdle)
		IdleTimeout: idleTimeout,  // 連接最大空閒時間 (對應原來的 IdleTimeout)
		
		// 其他可選配置
		// DialTimeout:  5 * time.Second,   // 連接超時
		// ReadTimeout:  3 * time.Second,   // 讀取超時
		// WriteTimeout: 3 * time.Second,   // 寫入超時
	})
}
