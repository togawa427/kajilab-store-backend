package controller

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)


func GetStorage(c *gin.Context) {
	bucketName  := "kajilab-store.appspot.com"	// GCSバケット名
    credentialsFile := "kajilab-store-cf40dbeb6615.json"	// サービスアカウント鍵ファイルのパス
		localFileName := "testgo.txt"
		gcsFileName := "path/file.txt"	// GCSバケットのアップロード先のパス

		 // Google Cloud Storageクライアントの作成
		 ctx := context.Background()
		 client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
		 if err != nil {
				 log.Fatal(err)
		 }
 
		 // アップロードするファイルを開く
		 file, err := os.Open(localFileName)
		 if err != nil {
				 log.Fatal(err)
		 }
		 defer file.Close()
 
		 // バケットオブジェクトの作成
		 bucket := client.Bucket(bucketName)
 
		 // バケット内のアップロード先のオブジェクトを作成
		 obj := bucket.Object(gcsFileName)
 
		 // ファイルをアップロード
		 wc := obj.NewWriter(ctx)
		 if _, err := io.Copy(wc, file); err != nil {
				 log.Fatal(err)
		 }
		 if err := wc.Close(); err != nil {
				 log.Fatal(err)
		 }
 
		 fmt.Println("file upload success")

    // // Google Cloud Storageクライアントの作成
    // ctx := context.Background()
    // client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
    // if err != nil {
    //     log.Fatal(err)
    // }

    // // バケットオブジェクトの作成
    // bucket := client.Bucket(bucketName)

    // // バケット内のオブジェクト一覧を取得
    // it := bucket.Objects(ctx, nil)
    // for {
    //     attrs, err := it.Next()
    //     if err == iterator.Done {
    //         break
    //     }
    //     if err != nil {
    //         log.Fatal(err)
    //     }
    //     fmt.Println(attrs.Name)
    // }

	c.JSON(http.StatusOK, "sucess")
}