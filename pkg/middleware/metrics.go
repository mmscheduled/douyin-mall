package middleware

import (
	"douyin/pkg/metrics"
	"github.com/cloudwego/hertz/pkg/app"
	"time"
	"context"
)

// MetricsMiddleware 监控中间件
func MetricsMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 记录请求开始时间
		start := time.Now()

		// 继续处理请求
		ctx.Next(c)

		// 将 []byte 转换为 string
		method := string(ctx.Method())
		path := string(ctx.Path())

		// 记录请求耗时和总数
		duration := time.Since(start).Seconds()
		metrics.RequestCount.WithLabelValues(method, path).Inc()
		metrics.RequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}