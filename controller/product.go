package controller

import (

	// "kajilab-store-backend/service"

	"kajilab-store-backend/model"
	"kajilab-store-backend/service"
	"net/http"
	"strconv"

	// // "strconv"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	ProductService := service.ProductService{}

	// limitとoffsetの取得
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "limit is not number")
		return
	}
	offset, err := strconv.ParseInt(c.Query("offset"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "offset is not number")
		return
	}

	// DBから商品情報取得
	products, err := ProductService.GetAllProducts(limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "fetal get products from DB")
		return
	}

	resProducts := []model.AllProductsGetResponse{}

	for _, product := range products {
		resProducts = append(resProducts, model.AllProductsGetResponse{
			Id: int64(product.ID),
			Name: product.Name,
			Barcode: product.Barcode,
			Price: product.Price,
			Stock: product.Stock,
			TagId: product.TagId,
			ImagePath: product.ImagePath,
		})
	}

	c.JSON(http.StatusOK, resProducts)
}

// func GetTagsByCommunityIdHandler(c *gin.Context) {
// 	communityId, err := strconv.ParseInt(c.Param("communityId"), 10, 64) // string -> int64
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, "Type is not number")
// 	}

// 	TagService := service.TagService{}

// 	// DBからどこのコミュニティにも該当するタグを持ってくる
// 	publicTags, err := TagService.GetTagsByCommunityId(-1)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get public tags"})
// 		return
// 	}
// 	// DBからコミュニティのタグネームを持ってくる
// 	communityTags, err := TagService.GetTagsByCommunityId(communityId)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get private tags"})
// 		return
// 	}
// 	tags := append(publicTags, communityTags...)

// 	tagsResponse := []model.TagsGetResponse{}

// 	for _, tag := range tags {

// 		tagsResponse = append(tagsResponse, model.TagsGetResponse{
// 			Id:   int64(tag.ID),
// 			Name: tag.Name,
// 		})
// 	}

// 	c.JSON(http.StatusOK, tagsResponse)
// }
