package util

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
)

func UploadFile(key, path string) {
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	putExtra := storage.PutExtra{
		Params: nil,
	}
	ret := storage.PutRet{}
	err := formUploader.PutFile(context.Background(), &ret, MakeToken(), key, path, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(ret.Key, ret.Hash)
}

func DownloadFile() {

}

func MakeToken() string {
	accessKey := ""
	secretKey := ""
	bucket := "store-vinf-top"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}
