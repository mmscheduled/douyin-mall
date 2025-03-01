package main

import (
	"context"
	"douyin/pkg/config"
	"douyin/pkg/database"
	"douyin/pkg/handler"
	"douyin/pkg/registry" 
	"douyin/pkg/middleware"
	"douyin/pkg/metrics" 
	"douyin/pkg/logger" 
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	// 加载配置文件
	config.LoadConfig()

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 启动 Prometheus 指标服务器
	metrics.StartMetricsServer(":9090")

	// 创建 Hertz 服务器
	h := server.Default(server.WithHostPorts(":8080"))

	// 注册全局中间件
	h.Use(middleware.AuthMiddleware())

	// 注册路由
	registerRoutes(h)

	// 创建 Consul 客户端
	consulClient, err := registry.NewConsulClient()
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	// 注册服务到 Consul
	if err := consulClient.RegisterService("douyin-service", 8080); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// 测试服务发现
	serviceAddress, err := consulClient.DiscoverService("douyin-service")
	if err != nil {
		log.Fatalf("Failed to discover service: %v", err)
	}
	log.Printf("Discovered service address: %s", serviceAddress)
	
	// 启动服务器
	if err := h.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// 刷新日志缓冲区
	defer logger.Sync()
}

// 注册路由
func registerRoutes(h *server.Hertz) {
	// 应用监控中间件
	h.Use(middleware.MetricsMiddleware())
	
	// 用户服务路由
	userGroup := h.Group("/api/user")
	{
		userGroup.POST("/register", handler.UserRegister)
		userGroup.POST("/login", handler.UserLogin)
		userGroup.GET("/info", handler.GetUserInfo)
		userGroup.POST("/logout", handler.UserLogout)
		userGroup.DELETE("/delete", handler.DeleteUser)
		userGroup.PUT("/update", handler.UpdateUser)
	}

	// 商品模块路由
	productGroup := h.Group("/api/product")
    {
        productGroup.POST("/add", handler.AddProduct)
        productGroup.GET("/list", handler.ListProducts)
        productGroup.GET("/detail", handler.GetProduct)
		productGroup.PUT("/update", handler.UpdateProduct)
		productGroup.DELETE("/delete", handler.DeleteProduct)
	}

	// 订单模块路由
	orderGroup := h.Group("/api/order")
    {
        orderGroup.POST("/create", handler.CreateOrder)
        orderGroup.GET("/list", handler.ListOrders)
    }
	//支付模块路由
	paymentGroup := h.Group("/api/payment")
    {
        paymentGroup.POST("/create", handler.CreatePayment)
    }
	//认证服务路由
	authGroup := h.Group("/api/auth")
	{
		authGroup.POST("/refresh", handler.RefreshToken)
	}
	//购物车路由
	cartGroup := h.Group("/api/cart")
    {
        cartGroup.POST("/add", handler.AddToCart)
        cartGroup.GET("/list", handler.GetCart)
        cartGroup.PUT("/update", handler.UpdateCartItem)
        cartGroup.DELETE("/delete", handler.DeleteCartItem)
    }


	// 健康检查路由
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})
}