package service

import (
	"errors"
	"fmt"
	"kajilab-store-backend/model"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TagService struct{
}

// 全てのタグを取得する
func (TagService) GetTags() ([]model.Tag, error){
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

	tags := make([]model.Tag, 0)
	q := db

	result := q.Find(&tags)
	if result.Error != nil {
		fmt.Printf("タグ取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return tags, nil
}

// タグIDからタグ情報を取得
func (TagService) GetTagById(id int64) (*model.Tag, error){
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

	tag := model.Tag{}
	q := db

	result := q.First(&tag, id)
	if result.Error != nil {
		fmt.Printf("タグ取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return &tag, nil
}

// タグ名からタグ情報を取得
func (TagService) GetTagByName(name string) (*model.Tag, error){
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

	tag := model.Tag{}
	q := db

	result := q.Where("name = ?", name).Take(&tag)
	if result.Error != nil {
		fmt.Printf("タグ取得失敗 %v", result.Error)
		return nil, result.Error
	}
	return &tag, nil
}

// タグ情報を登録
func (TagService) CreateTag(tag *model.Tag) error {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	q := db
	// タグ名が既に登録済みかをチェック
	existTag := model.Tag{}
	result := q.Where("name = ?", tag.Name).First(&existTag)
	if result.Error == nil {
		err := errors.New("the barcode is existing")
		fmt.Printf("%v",err)
		return err
	}

	// 商品をDBへ登録
	result = q.Create(tag)
	if result.Error != nil {
		fmt.Printf("タグ登録失敗 %v", result.Error)
		return result.Error
	}
	return nil
}