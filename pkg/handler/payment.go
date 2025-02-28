package handler

import (
    "douyin/pkg/database"
    "douyin/pkg/model"
    "net/http"
	"context"
	"github.com/cloudwego/hertz/pkg/common/utils"

    "github.com/cloudwego/hertz/pkg/app"
)

func CreatePayment(c context.Context, ctx *app.RequestContext) {
    var req struct {
        OrderID uint `json:"order_id"`
    }
    if err := ctx.Bind(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
        return
    }

    // 获取订单信息
    var order model.Order
    if err := database.DB.Where("id = ?", req.OrderID).First(&order).Error; err != nil {
        ctx.JSON(http.StatusNotFound, utils.H{"error": "order not found"})
        return
    }

    // 创建支付记录
    payment := model.Payment{
        OrderID: req.OrderID,
        Amount:  order.TotalPrice,
    }
    if err := database.DB.Create(&payment).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to create payment"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "payment_id": payment.ID,
    })
}