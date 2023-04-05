package helper

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	"github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, identity, name string, second int) (string, error) {
	// id
	// identity
	// name
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(second)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	fmt.Println(tokenString)
	return tokenString, nil
}

// 邮箱验证码的发送
// 这里的目的邮箱是你要填的

func MailSendCode(mail, code string) error {
	e := email.NewEmail()
	e.From = "Get <z15874610586@163.com>"
	e.To = []string{mail}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("你的验证码为：<h1>" + code + "</h1>")
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "z15874610586@163.com", define.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})

	if err != nil {
		return err
	}
	return nil
}

// 生成一个随机code

func RandCode() string {
	s := "123456789"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < define.CodeLength; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

func UUID() string {
	return uuid.NewV4().String()
}

//上传到腾讯云

func CosUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(define.CosBucket) //存储桶的路径
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	//key := "cloud-disk/test2.jpg" //地址
	file, fileHeader, err := r.FormFile("file")

	key := "cloud-disk/" + UUID() + path.Ext(fileHeader.Filename)

	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		panic(err)
	}
	return define.CosBucket + "/" + key, nil
}

// AnalyzeToken
// Token 解析
func AnalyzeToken(token string) (*define.UserClaim, error) {
	uc := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token is invalid")
	}
	return uc, err
}

// 分片上传初始化
func CosInitPart(ext string) (string, string, error) {
	u, _ := url.Parse(define.CosBucket) //要填自己的
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})
	key := "cloud-disk/" + UUID() + ext
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		return "", "", err
	}
	//UploadID := v.UploadID // 1680686259f21e7ab0bd2f9423272b218df83da6c68e4571e25269465114e1bf460d5c7bb8
	UploadID := v.UploadID // 16806895929da13c4d006ccc38f8a777117dd106da1ff1dbd0a5548b9ec2fd82a63e390fad
	fmt.Println(UploadID)
	return key, v.UploadID, nil
}

// 分片上传
func CosPartUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})
	key := r.PostForm.Get("key")

	UploadID := r.PostForm.Get("upload_id")

	partNumber, err := strconv.Atoi(r.PostForm.Get("part_number"))

	f, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	// opt可选
	resp, err := client.Object.UploadPart(
		context.Background(), key, UploadID, partNumber, bytes.NewReader(buf.Bytes()), nil,
	)
	if err != nil {
		return "", err
	}
	return strings.Trim(resp.Header.Get("ETag"), "\""), nil
}

// CosPartUploadComplete 分片上传完成
func CosPartUploadComplete(key, uploadId string, co []cos.Object) error {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, co...)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, uploadId, opt,
	)
	return err
}
