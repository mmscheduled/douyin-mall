package handler

import (
	"douyin/pkg/database"
	"douyin/pkg/model"
	"net/http"
	"context"
	"github.com/cloudwego/hertz/pkg/common/utils"


	"github.com/cloudwego/hertz/pkg/app"
)

func CreateOrder(c context.Context, ctx *app.RequestContext) {
	var req struct {
		UserID    uint `json:"user_id"`
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
		return
	}

	// 获取商品信息
	var product model.Product
	if err := database.DB.Where("id = ?", req.ProductID).First(&product).Error; err != nil {
		ctx.JSON(http.StatusNotFound, utils.H{"error": "product not found"})
		return
	}

	// 计算总价
	totalPrice := product.Price * float64(req.Quantity)

	// 创建订单
	order := model.Order{
		UserID:     req.UserID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: totalPrice,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to create order"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"order_id": order.ID,
	})
}

//添加订单列表
func ListOrders(c context.Context, ctx *app.RequestContext) {
    userID := ctx.Query("user_id")
    var orders []model.Order
    if err := database.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to fetch orders"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "orders": orders,
    })
}
