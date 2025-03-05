package controller

import (
	"fmt"
	"kajilab-store-backend/model"
	"kajilab-store-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTags(c *gin.Context) {
	TagService := service.TagService{}

	// DBからタグ情報取得
	tags, err := TagService.GetTags()
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get tags from DB")
		return
	}

	resTags := []model.TagGetResponseTag{}

	for _, tag := range tags {
		resTags = append(resTags, model.TagGetResponseTag{
			Name: tag.Name,
		})
	}

	response := model.TagGetResponse{
		Logs: resTags,
	}

	c.JSON(http.StatusOK, response)
}

// タグ登録API
func CreateTag(c *gin.Context){
	TagService := service.TagService{}
	TagCreateRequest := model.TagPostRequest{}
	err := c.Bind(&TagCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "request is not correct")
		return
	}

	// リクエストの商品情報をデータベースの型へ変換
	tag := model.Tag{
		Name: TagCreateRequest.Tag.Name,
	}

	err = TagService.CreateTag(&tag)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal create tag")
		return
	}

	c.JSON(http.StatusOK, "success")
}