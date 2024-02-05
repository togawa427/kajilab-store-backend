package controller

import (
	"kajilab-store-backend/model"
	"kajilab-store-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAsset(c *gin.Context) {
	AssetService := service.AssetService{}

	// DBから予算情報を取得
	asset, err := AssetService.GetAsset()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get products from DB")
		return
	}

	// データベースの予算情報をレスポンスの型へ変換
	jsonAsset := model.AssetGetResponse{
		Money: asset.Money,
		Debt: asset.Debt,
	}

	c.JSON(http.StatusOK, jsonAsset)
}