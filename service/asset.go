package service

import (
	"fmt"
	"kajilab-store-backend/model"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AssetService struct{}


// お金を増減させる
func (AssetService) IncreaseMoney(money int64) (error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 現在のお金を取得
	asset := model.Asset{}
	result := db.First(&asset)
	if result.Error != nil {
		fmt.Printf("財産取得失敗 %v", result.Error)
		return result.Error
	}

	// DBへ更新
	result = db.Model(&model.Asset{}).Where("id = 1").Update("money", asset.Money + money)
	if result.Error != nil {
		fmt.Printf("財産更新失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// 現在の予算を取得
func (AssetService) GetAsset() (model.Asset, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return model.Asset{}, err
	}

	// 現在の予算をDBから取得
	asset := model.Asset{}
	result := db.First(&asset)
	if result.Error != nil {
		fmt.Printf("財産取得失敗 %v", result.Error)
		return model.Asset{}, result.Error
	}

	return asset, nil
}