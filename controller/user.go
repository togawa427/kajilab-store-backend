package controller

import (
	"fmt"
	"kajilab-store-backend/model"
	"kajilab-store-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	UserService := service.UserService{}
	UserCreateRequest := model.UserCreateRequest{}
	err := c.Bind(&UserCreateRequest)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// リクエストの商品情報をデータベースの型へ変換
	user := model.User{
		Name: UserCreateRequest.Name,
		Debt: 0,
		Barcode: UserCreateRequest.Barcode,
	}
	
	err = UserService.CreateUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create product")
		return
	}

	c.JSON(http.StatusOK, "success")
}