package controller

import (
	"fmt"
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

func UpdateAsset(c *gin.Context) {
	AssetService := service.AssetService{}
	AssetUpdateRequest := model.AssetUpdateRequest{}
	err := c.Bind(&AssetUpdateRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 財産情報更新
	// リクエストの財産情報をデータベースの型へ変換
	asset := model.Asset{
		Money: AssetUpdateRequest.Money,
		Debt: AssetUpdateRequest.Debt,
	}
	// DBへ保存
	err = AssetService.UpdateAsset(&asset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal update product")
		return
	}

	c.JSON(http.StatusOK, "success")
}