package service

import (
	"fmt"
	"kajilab-store-backend/model"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TagMapService struct{
}

// タグ情報を登録
func (TagMapService) CreateTagMap(tagMap *model.TagMap) (*model.TagMap, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	q := db

	// 商品をDBへ登録
	result := q.Create(tagMap)
	if result.Error != nil {
		fmt.Printf("タグ登録失敗 %v", result.Error)
		return nil, result.Error
	}
	return tagMap, nil
}