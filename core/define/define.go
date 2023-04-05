package define

import (
	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.StandardClaims
}

var JwtKey = "cloud-disk-key"
var MailPassword = "FXEIYYVYVFFKNRWD"

var CodeLength = 6

// 验证码过期时间

var CodeExpire = 300

// TencentSecretKey 腾讯云对象存储
//var TencentSecretKey = os.Getenv("TencentSecretKey")
var TencentSecretKey = "A8kHl4LWSBj6JK8k6GNGuISVMyNdiYBY"

//var TencentSecretID = os.Getenv("TencentSecretID")
var TencentSecretID = "AKIDHVYLvMl1B9sMbHjx63pXEk00wyrdIm4I"

var CosBucket = "https://1-1307688964.cos.ap-beijing.myqcloud.com"

// 分页的默认参数
var PageSize int = 20

var DateTime = "2006-01-01 15:01:01"

var TokenExpire = 3600
var RefreshTokenExpire = 7200
