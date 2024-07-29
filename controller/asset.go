package controller

import (
	"fmt"
	"kajilab-store-backend/model"
	"kajilab-store-backend/service"
	"net/http"
	"strconv"
	"time"

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

func GetAssetHistory(c *gin.Context) {

	day, err := strconv.ParseInt(c.Query("day"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "day is not number")
		return
	}
	if day < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "day is low")
		return
	}

	AssetService := service.AssetService{}
	ProductService := service.ProductService{}

	// DBから予算情報を取得
	dayMargin := int64(30)	// day日前に予算情報がない場合用のマージン(少なくともdayの前30日間で一つは予算情報があるはず)
	assets, err := AssetService.GetAssetHistory(day + dayMargin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get assets from DB")
		return
	}

	// DBから商品価値を取得
	productsValues, err := ProductService.GetProductsValuesByDay(day + dayMargin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get products values from DB")
		return
	}

	// データベースの予算情報をレスポンスの型へ変換
	resAssets := []model.AssetHistoryGetResponse{}
	latestMoney := int64(-1)
	latestDebt := int64(-1)
	i := 0
	for _, asset := range assets {
		currentDate := time.Now().AddDate(0, 0, 0-(int(day+dayMargin)-i))
		if(asset.Money != -1){
			// お金情報がない日
			latestMoney = asset.Money
			latestDebt = asset.Debt
		}
		resAssets = append(resAssets, model.AssetHistoryGetResponse{
			Date: currentDate,
			Money: latestMoney,
			Debt: latestDebt,
			Product: productsValues[i],
		})
		i++
	}

	c.JSON(http.StatusOK, resAssets[dayMargin:])
}

func UpdateAsset(c *gin.Context) {
	AssetService := service.AssetService{}
	UserService := service.UserService{}
	AssetUpdateRequest := model.AssetUpdateRequest{}
	err := c.Bind(&AssetUpdateRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// 財産情報更新
	totalDebt, err := UserService.GetUsersTotalDebt()
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get total debt")
		return
	}

	// リクエストの財産情報をデータベースの型へ変換
	asset := model.Asset{
		Money: AssetUpdateRequest.Money,
		Debt: totalDebt,
	}
	// DBへ保存
	err = AssetService.UpdateAsset(&asset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal update product")
		return
	}

	c.JSON(http.StatusOK, "success")
}