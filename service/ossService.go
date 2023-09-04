package service

import (
	"douyin/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"path/filepath"
)

var bucket *oss.Bucket

func InitOSS() error {
	client, err := oss.New(config.GlobalConfig.Database.OSS.Endpoint,
		config.GlobalConfig.Database.OSS.AccessKeyID,
		config.GlobalConfig.Database.OSS.AccessKeySecret)
	if err != nil {
		log.Printf("oss init failed: %v", err)
		return err
	}

	bucket, err = client.Bucket(config.GlobalConfig.Database.OSS.BucketName)
	if err != nil {
		fmt.Printf("oss get bucket failed: %v", err)
		return err
	}
	return nil
}

func generateFileName(filename string) string {
	return uuid.New().String() + filepath.Ext(filename)
}

func UploadVideoToOSS(videoReader *multipart.FileHeader) (string, string, error) {

	filename := filepath.Base(videoReader.Filename)
	videoName := generateFileName(filename)
	fmt.Printf("videoName:%v\n", videoName)
	video, err := videoReader.Open()
	if err != nil {
		log.Printf("oss open video fileReader failed: %v", err)
		return "", "", err
	}

	err = bucket.PutObject(videoName, video)
	if err != nil {
		log.Printf("oss upload video failed: %v", err)
		return "", "", err
	}

	videoURL, imgURL := getURLFromOSS(videoName)
	return videoURL, imgURL, nil
}

func getURLFromOSS(filename string) (string, string) {
	videoURL := fmt.Sprintf("https://%s.%s/%s",
		config.GlobalConfig.Database.OSS.BucketName,
		config.GlobalConfig.Database.OSS.Endpoint,
		filename)
	coverURL := videoURL + "?x-oss-process=video/snapshot,t_0,f_jpg,w_0,h_0,m_fast,ar_auto"
	return videoURL, coverURL
}
