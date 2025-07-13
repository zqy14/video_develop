package middilware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
）
func RouTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UnixMilli()

		fmt.Println("Runtime运行前")

		c.Next()

		end := time.Now().UnixMilli()

		fmt.Sprintf("代码运行时间%d", end-start)
	}
}
	
// RedisLatencyRecorder 基于Redis的请求耗时记录器
type RedisLatencyRecorder struct {
	rdb        *redis.Client
	metricsKey string        // Redis中存储指标的key
	maxRecords int64         // 最大保存记录数
	expireTime time.Duration // 数据过期时间
}

// Middleware 返回Gin中间件函数
func (r *RedisLatencyRecorder) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求前
		fmt.Println("请求开始处理")

		// 执行后续中间件和处理器
		c.Next()

		// 请求处理完成后
		end := time.Now()
		latency := end.Sub(start).Milliseconds()

		// 记录到Redis
		go r.recordLatency(c, latency)

		// 打印日志
		fmt.Printf("请求处理完成，耗时 %d ms\n", latency)
	}
}

// recordLatency 记录耗时到Redis
func (r *RedisLatencyRecorder) recordLatency(c *gin.Context, latency int64) {
	ctx := context.Background()

	// 使用Pipeline提高性能
	pipe := r.rdb.Pipeline()

	// 1. 记录当前请求耗时
	pipe.ZAdd(ctx, r.metricsKey, redis.Z{
		Score:  float64(time.Now().Unix()), // 使用时间戳作为score
		Member: fmt.Sprintf("%s:%s:%d", c.Request.Method, c.Request.URL.Path, latency),
	})

	// 2. 限制集合大小
	pipe.ZRemRangeByRank(ctx, r.metricsKey, 0, -r.maxRecords-1)

	// 3. 设置过期时间
	pipe.Expire(ctx, r.metricsKey, r.expireTime)

	// 执行Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		fmt.Printf("记录请求耗时到Redis失败: %v\n", err)
	}
}
