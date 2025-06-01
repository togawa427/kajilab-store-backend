package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"kajilab-store-backend/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserService struct{}

// IDからユーザ情報取得
func (UserService) GetUserById(id int64) (model.User, error) {
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
	result := db.First(&user, id)
	if result.Error != nil {
		fmt.Printf("ユーザ取得失敗 %v", result.Error)
		return user, result.Error
	}
	return user, nil
}

// バーコードからユーザ情報取得
func (UserService) GetUserByBarcode(barcode string) (model.User, error) {
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

// ユーザ保有残高の総額を取得
func (UserService) GetUsersTotalDebt() (int64, error) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE_NAME")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer sqlDB.Close()

	// 取得
	users := make([]model.User, 0)
	result := db.Find(&users)
	if result.Error != nil {
		fmt.Printf("ユーザ取得失敗 %v", result.Error)
		return 0, result.Error
	}
	totalDebt := int64(0)
	for _, user := range users {
		totalDebt += user.Debt
	}
	return totalDebt, nil
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
		fmt.Printf("%v", err)
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

	fmt.Println("ユーザIDは　" + strconv.Itoa(int(id)))
	fmt.Println("ユーザ名は　" + user.Name)

	// ログに出力
	go createCloudLog(user.Name + "(ID: " + strconv.Itoa(int(id)) + ")の残高\n" + "更新後：" + strconv.Itoa(int(user.Debt)))

	return nil
}

// 残高を増減させる
func (UserService) IncreaseUserDebt(userId int64, debt int64) error {
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

	// 現在の残高を取得
	user := model.User{}
	result := db.First(&user, userId)
	if result.Error != nil {
		fmt.Printf("ユーザ残高取得失敗 %v", result.Error)
		return result.Error
	}

	// DBへ更新
	result = db.Model(&model.User{}).Where("id = ?", userId).Update("debt", user.Debt+debt)
	// result = db.Create(&model.Asset{Money: .Money+money, Debt: asset.Debt})
	if result.Error != nil {
		fmt.Printf("ユーザ残高更新失敗 %v", result.Error)
		return result.Error
	}

	// ログに出力
	go createCloudLog(user.Name + "(ID: " + strconv.Itoa(int(userId)) + ")" + "の残高\n" + "支払い前：" + strconv.Itoa(int(user.Debt)) + "\n支払い後：" + strconv.Itoa(int(user.Debt+debt)))

	return nil
}

func (UserService) IsEnoughUserDebt(userId int64, price int64) error {
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

	// 現在の残高を取得
	user := model.User{}
	result := db.First(&user, userId)
	if result.Error != nil {
		fmt.Printf("ユーザ残高取得失敗 %v", result.Error)
		return result.Error
	}

	// 残高更新後マイナスになる場合エラー
	if user.Debt-price < 0 {
		fmt.Println("ユーザ残高が足りません")
		err := errors.New("the debt is low")
		return err
	}
	return nil
}

func (UserService) CreateKajilabPayLog(userId int64, prevDebt int64, currentDebt int64) {
	// DBに更新前の残高，更新後の残高，UserId，を保存

	//
}

// 残高を増減させる
func (UserService) IncreaseKajilabpayDebt(userId int64, paymentId int64, debt int64, content string) error {
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

	fmt.Println("1 IncreaseKajilabpayDebt")

	// 現在の残高を取得
	user := model.User{}
	result := db.First(&user, userId)
	if result.Error != nil {
		fmt.Printf("ユーザ残高取得失敗 %v", result.Error)
		return result.Error
	}

	fmt.Println("2 IncreaseKajilabpayDebt")

	// DBのユーザ残高情報を更新
	result = db.Model(&model.User{}).Where("id = ?", userId).Update("debt", user.Debt+debt)
	// result = db.Create(&model.Asset{Money: .Money+money, Debt: asset.Debt})
	if result.Error != nil {
		fmt.Printf("ユーザ残高更新失敗 %v", result.Error)
		return result.Error
	}

	fmt.Println("3 IncreaseKajilabpayDebt")

	// 梶研Payの残高履歴を更新
	// UserId，更新前の残高，更新後の残高，PaymentId，Contentを保存
	kajilabpayLog := model.KajilabpayLog{
		UserId:      userId,
		PrevDebt:    user.Debt,
		CurrentDebt: user.Debt + debt,
		PaymentId:   paymentId,
		Content:     content,
	}
	result = db.Create(&kajilabpayLog)
	if result.Error != nil {
		fmt.Printf("梶研Pay履歴登録失敗 %v", result.Error)
		return result.Error
	}

	fmt.Println("4 IncreaseKajilabpayDebt")

	// = 現在の総財産情報をDBへ登録 =
	// 現在の財産情報を取得
	asset := model.Asset{}
	result = db.Last(&asset)
	if result.Error != nil {
		fmt.Printf("残高取得失敗 %v", result.Error)
		return result.Error
	}
	newAsset := model.Asset{
		Money: asset.Money,
		Debt:  asset.Debt + debt,
	}
	// 残高情報をDBへ保存
	asset.Debt = asset.Debt + debt
	result = db.Create(&newAsset)
	if result.Error != nil {
		fmt.Printf("残高更新失敗 %v", result.Error)
		return result.Error
	}

	// CloudLoggingにログ出力
	go createCloudLog(user.Name + "(ID: " + strconv.Itoa(int(userId)) + ")" + "の残高\n" + "支払い前：" + strconv.Itoa(int(user.Debt)) + "\n支払い後：" + strconv.Itoa(int(user.Debt+debt)))

	return nil
}
