package service

import (
	"errors"
	"fmt"
	"kajilab-store-backend/model"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserService struct{}

// バーコードからユーザ情報取得
func (UserService) GetUserByBarcode (barcode string) (model.User, error) {
	user := model.User{}
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 取得
	result := db.Where(&model.User{Barcode: barcode}).First(&user)
	if result.Error != nil {
		fmt.Printf("ユーザ取得失敗 %v", result.Error)
		return user, result.Error
	}
	return user, nil
}

// ユーザ情報を登録
func (UserService) CreateUser(user *model.User) error {
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

	// バーコードのかぶりチェック
	existUser := model.User{}
	result := db.Where("barcode = ?", user.Barcode).First(&existUser)
	if result.Error == nil {
		err := errors.New("the barcode is existing")
		fmt.Printf("%v",err)
		return err
	}

	// ユーザをDBへ登録
	result = db.Create(user)
	if result.Error != nil {
		fmt.Printf("ユーザ登録失敗 %v", result.Error)
		return result.Error
	}
	return nil
}

// ユーザ情報編集
func (UserService) UpdateUser(id int64, user *model.User) error {
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

	// ユーザ情報をDBへ登録
	result := db.Model(&model.User{}).Where("id = ?", id).Updates(&user)
	if result.Error != nil {
		fmt.Printf("ユーザ情報更新失敗 %v", result.Error)
		return result.Error
	}
	return nil
}