package oss

import (
	"bytes"
	"errors"
	"fmt"
	ali_oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// 阿里云OSS, Document: https://github.com/aliyun/aliyun-oss-go-sdk/tree/master/sample

type AliYunOss struct {
	url    string
	client *ali_oss.Client
	bucket *ali_oss.Bucket
}

func NewAliYunOss(c AliYun) *AliYunOss {
	client, err := ali_oss.New(
		c.EndPoint,
		c.AccessId,
		c.AccessSecret,
	)
	if err != nil {
		zap.S().Errorf("Init Oss Failed. Err: %v", err)
		os.Exit(-1)
		return nil
	}
	bucket, err := client.Bucket(c.DefaultBucket)
	if err != nil {
		zap.S().Errorf("Init Bucket Failed. Err: %v, bucket %s, c: %+v", err, c.DefaultBucket, c)
		os.Exit(-1)
		return nil
	}
	return &AliYunOss{
		url:    c.CustomDomain,
		client: client,
		bucket: bucket,
	}
}

// UploadFile 上传本地文件
func (ao *AliYunOss) UploadFile(objectKey, localFilePath string, opt ...ali_oss.Option) (err error) {
	fd, err := os.Open(localFilePath)
	if err != nil {
		zap.S().Errorf("Oss [os.Open] Error: %v [objectKey: %s, localFilePath: %s]", err, objectKey, localFilePath)
		return err
	}
	defer fd.Close()

	if err = ao.bucket.PutObject(objectKey, fd, opt...); err != nil {
		zap.S().Errorf("Oss [ao.bucket.PutObject] Error: %v [objectKey: %s, localFilePath: %s]", err, objectKey, localFilePath)
		return err
	}

	return nil
}

// DownloadFile 下载文件
func (ao *AliYunOss) DownloadFile(objectKey string) (data []byte, err error) {
	body, err := ao.bucket.GetObject(objectKey)
	if err != nil {
		zap.S().Errorf("Oss [ao.bucket.GetObject] Error: %v [objectKey: %s]", err, objectKey)
		return data, err
	}
	if data, err = ioutil.ReadAll(body); err != nil {
		zap.S().Errorf("Oss [ioutil.ReadAll] Error: %v [objectKey: %s, body: %v]", err, objectKey, body)
		return data, err
	}
	body.Close()

	return data, nil
}

// UploadByteFile 上传文件
func (ao *AliYunOss) UploadByteFile(objectKey string, file []byte, opt ...ali_oss.Option) (url string, err error) {
	if err = ao.bucket.PutObject(objectKey, bytes.NewReader(file), opt...); err != nil {
		zap.S().Errorf("upload oss failed, err: %v", err)
		return "", errors.New("upload oss failed")
	}

	return fmt.Sprintf("https://%s/%s", ao.url, objectKey), nil
}

// UploadUrl 上传网络文件
func (ao *AliYunOss) UploadUrl(objectKey string, url string, opt ...ali_oss.Option) (aliUrl string, err error) {
	res, err := http.Get(url)
	if err != nil {
		zap.S().Errorf("Get Image Err, err: %v", err)
		return "", errors.New("get url error")
	}
	if err = ao.bucket.PutObject(objectKey, io.Reader(res.Body), opt...); err != nil {
		zap.S().Errorf("upload oss failed, err: %v", err)
		return "", errors.New("upload oss failed")
	}
	return fmt.Sprintf("https://%s/%s", ao.url, objectKey), nil
}
