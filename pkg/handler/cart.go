package handler

import (
    "context"
    "douyin/pkg/database"
    "douyin/pkg/model"
    "net/http"
	"github.com/cloudwego/hertz/pkg/common/utils"

    "github.com/cloudwego/hertz/pkg/app"
)

func AddToCart(c context.Context, ctx *app.RequestContext) {
    var req struct {
        UserID    uint `json:"user_id"`
        ProductID uint `json:"product_id"`
        Quantity  int  `json:"quantity"`
    }
    if err := ctx.Bind(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
        return
    }

    // 检查商品是否存在
    var product model.Product
    if err := database.DB.Where("id = ?", req.ProductID).First(&product).Error; err != nil {
        ctx.JSON(http.StatusNotFound, utils.H{"error": "product not found"})
        return
    }

    // 添加商品到购物车
    cartItem := model.CartItem{
        UserID:    req.UserID,
        ProductID: req.ProductID,
        Quantity:  req.Quantity,
    }
    if err := database.DB.Create(&cartItem).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to add to cart"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "cart_item_id": cartItem.ID,
    })
}

func GetCart(c context.Context, ctx *app.RequestContext) {
    userID := ctx.Query("user_id")
    var cartItems []model.CartItem
    if err := database.DB.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to fetch cart items"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "cart_items": cartItems,
    })
}

func UpdateCartItem(c context.Context, ctx *app.RequestContext) {
    var req struct {
        CartItemID uint `json:"cart_item_id"`
        Quantity   int  `json:"quantity"`
    }
    if err := ctx.Bind(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
        return
    }

    // 更新购物车商品数量
    if err := database.DB.Model(&model.CartItem{}).Where("id = ?", req.CartItemID).Update("quantity", req.Quantity).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to update cart item"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "message": "cart item updated",
    })
}

func DeleteCartItem(c context.Context, ctx *app.RequestContext) {
    cartItemID := ctx.Query("cart_item_id")
    if err := database.DB.Where("id = ?", cartItemID).Delete(&model.CartItem{}).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to delete cart item"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "message": "cart item deleted",
    })
}