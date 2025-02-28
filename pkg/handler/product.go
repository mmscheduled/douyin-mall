package handler

import (
	"context"
	"douyin/pkg/database"
	"douyin/pkg/model"
	"net/http"

	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/app"
)

//创建商品
func AddProduct(c context.Context, ctx *app.RequestContext) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
		return
	}

	product := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	if err := database.DB.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to create product"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"product_id": product.ID,
	})
}

//商品列表
func ListProducts(c context.Context, ctx *app.RequestContext) {
    var products []model.Product
    if err := database.DB.Find(&products).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to fetch products"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "products": products,
    })
}

//添加商品详情
func GetProduct(c context.Context, ctx *app.RequestContext) {
    productID := ctx.Query("product_id")
    var product model.Product
    if err := database.DB.Where("id = ?", productID).First(&product).Error; err != nil {
        ctx.JSON(http.StatusNotFound, utils.H{"error": "product not found"})
        return
    }

    ctx.JSON(http.StatusOK, utils.H{
        "product": product,
    })
}

func UpdateProduct(c context.Context, ctx *app.RequestContext) {
	var req struct {
		ProductID   uint    `json:"product_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.H{"error": "invalid request"})
		return
	}

	// 更新商品信息
	if err := database.DB.Model(&model.Product{}).Where("id = ?", req.ProductID).Updates(map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
		"stock":       req.Stock,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to update product"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"message": "product updated successfully",
	})
}

func DeleteProduct(c context.Context, ctx *app.RequestContext) {
	productID := ctx.Query("product_id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, utils.H{"error": "product_id is required"})
		return
	}

	// 删除商品
	if err := database.DB.Where("id = ?", productID).Delete(&model.Product{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"message": "product deleted successfully",
	})
}