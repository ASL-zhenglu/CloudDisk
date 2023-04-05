package test

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestFileUploadByFilePath(t *testing.T) {
	u, _ := url.Parse("https://1-1307688964.cos.ap-beijing.myqcloud.com") //存储桶的路径
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "cloud-disk/test.jpg"

	_, _, err := client.Object.Upload(
		context.Background(), key, "./img/lol.mp4", nil,
	)
	if err != nil {
		panic(err)
	}
}

func TestFileUploadByReader(t *testing.T) {
	u, _ := url.Parse("https://1-1307688964.cos.ap-beijing.myqcloud.com") //存储桶的路径
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	key := "cloud-disk/test2.mp4"

	f, err := os.ReadFile("./img/lol.mp4")
	if err != nil {
		return
	}
	_, err = client.Object.Put(
		context.Background(), key, bytes.NewReader(f), nil,
	)
	if err != nil {
		panic(err)
	}
}

// 分片上传初始化
func TestInitPartUpload(t *testing.T) {
	u, _ := url.Parse("https://1-1307688964.cos.ap-beijing.myqcloud.com") //要填自己的
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})
	key := "cloud-disk/test2.mp4"
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		t.Fatal(err)
	}
	//UploadID := v.UploadID // 1680686259f21e7ab0bd2f9423272b218df83da6c68e4571e25269465114e1bf460d5c7bb8
	UploadID := v.UploadID // 16806895929da13c4d006ccc38f8a777117dd106da1ff1dbd0a5548b9ec2fd82a63e390fad
	fmt.Println(UploadID)
}

// 分片上传
func TestPartUpload(t *testing.T) {
	u, _ := url.Parse("https://1-1307688964.cos.ap-beijing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})
	key := "cloud-disk/test2.mp4"
	//UploadID := "1680686259f21e7ab0bd2f9423272b218df83da6c68e4571e25269465114e1bf460d5c7bb8"
	UploadID := "16806895929da13c4d006ccc38f8a777117dd106da1ff1dbd0a5548b9ec2fd82a63e390fad"

	//f, err := os.ReadFile("0.chunk") // md5 : 8dcec7e456dc23eaf2fecdb4483e5deb
	//f, err := os.ReadFile("0.chunk") // md5 : c103991d5dd56be442a9881a1d1acf5a
	//f, err := os.ReadFile("1.chunk") // md5 : b23753abaa8bb79475f31d394f70bc57
	f, err := os.ReadFile("2.chunk") // md5 : 364a2b56e9ccb97c5316f7e8622de54d

	if err != nil {
		t.Fatal(err)
	}
	// opt可选
	resp, err := client.Object.UploadPart(
		context.Background(), key, UploadID, 3, bytes.NewReader(f), nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	PartETag := resp.Header.Get("ETag")
	fmt.Println(PartETag)
}

// 分片上传完成
func TestPartUploadComplete(t *testing.T) {
	u, _ := url.Parse("https://1-1307688964.cos.ap-beijing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})
	key := "cloud-disk/test2.mp4"
	//UploadID := "1680686259f21e7ab0bd2f9423272b218df83da6c68e4571e25269465114e1bf460d5c7bb8"
	UploadID := "16806895929da13c4d006ccc38f8a777117dd106da1ff1dbd0a5548b9ec2fd82a63e390fad"

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, cos.Object{
		PartNumber: 1, ETag: "c103991d5dd56be442a9881a1d1acf5a"},
		cos.Object{
			PartNumber: 2, ETag: "b23753abaa8bb79475f31d394f70bc57"},
		cos.Object{
			PartNumber: 3, ETag: "364a2b56e9ccb97c5316f7e8622de54d"},
	)

	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)
	if err != nil {
		t.Fatal(err)
	}
}
