package controller

import (

	// "kajilab-store-backend/service"
	"net/http"
	"strconv"

	// // "strconv"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	//ProductService := service.ProductService{}

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

	// fmt.Println(limit)
	// fmt.Println(offset)

	// // DBから商品情報取得
	// fmt.Println(ProductService.GetAllProducts())

	// products := []model.Product{
	// 	{
	// 		Id: 1,
	// 		Name: "じゃがりこサラダ味",
	// 		Barcode: 45234123423,
	// 		Price: 120,
	// 		Stock: 7,
	// 		ImagePath: "public/images/jagariko.jpg",
	// 	},
	// 	{
	// 		Id: 2,
	// 		Name: "じゃがりこチーズ味",
	// 		Barcode: 45234123422,
	// 		Price: 120,
	// 		Stock: 4,
	// 		ImagePath: "public/images/jagarikotisu.jpg",
	// 	},
	// }



	// resProducts := []model.AllProductsGetResponse{}

	// for _, product := range products {
	// 	resProducts = append(resProducts, model.AllProductsGetResponse{
	// 		Id: int64(product.Id),
	// 		Name: product.Name,
	// 		Barcode: product.Barcode,
	// 		Price: product.Price,
	// 		Stock: product.Stock,
	// 		ImagePath: product.ImagePath,
	// 	})
	// }

	// c.JSON(http.StatusOK, resProducts)

	c.JSON(http.StatusOK, "limit:" + strconv.FormatInt(limit,10) + " offset:" + strconv.FormatInt(offset,10))

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
