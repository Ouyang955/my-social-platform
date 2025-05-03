package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"my-social-platform/internal/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// RSA密钥对,用于JWT的签名和验证
var (
	privateKey *rsa.PrivateKey // 私钥用于签名JWT
	publicKey  *rsa.PublicKey  // 公钥用于验证JWT签名
)

// init 在包初始化时执行,负责RSA密钥对的初始化
// 首先尝试从文件加载已有的密钥对,如果失败则生成新的密钥对并保存
func init() {
	var err error
	// 尝试从文件加载密钥
	if err = loadKeys(); err != nil {
		// 如果加载失败，生成新的密钥对
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic("Failed to generate RSA key pair")
		}
		publicKey = &privateKey.PublicKey

		// 保存新生成的密钥到文件系统
		if err = saveKeys(); err != nil {
			panic("Failed to save RSA key pair")
		}
	}
}

// saveKeys 将RSA密钥对保存到文件系统
// 私钥保存为private.pem,权限设为600(仅当前用户可读写)
// 公钥保存为public.pem,权限设为644(所有用户可读,仅当前用户可写)
func saveKeys() error {
	// 创建keys目录,权限设为700(仅当前用户可读写执行)
	if err := os.MkdirAll("keys", 0700); err != nil {
		return err
	}

	// 将私钥转换为PKCS1格式并用PEM编码
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	if err := ioutil.WriteFile(
		filepath.Join("keys", "private.pem"),
		pem.EncodeToMemory(privateKeyPEM),
		0600,
	); err != nil {
		return err
	}

	// 将公钥转换为PKCS1格式并用PEM编码
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return ioutil.WriteFile(
		filepath.Join("keys", "public.pem"),
		pem.EncodeToMemory(publicKeyPEM),
		0644,
	)
}

// loadKeys 从文件系统加载RSA密钥对
// 从private.pem加载私钥,从public.pem加载公钥
// 如果任一文件不存在或格式错误,返回error
func loadKeys() error {
	// 从private.pem加载并解析私钥
	privateKeyBytes, err := ioutil.ReadFile(filepath.Join("keys", "private.pem"))
	if err != nil {
		return err
	}
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	// 从public.pem加载并解析公钥
	publicKeyBytes, err := ioutil.ReadFile(filepath.Join("keys", "public.pem"))
	if err != nil {
		return err
	}
	block, _ = pem.Decode(publicKeyBytes)
	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return err
	}

	return nil
}

// GenerateJWT 生成JWT(JSON Web Token)
// 使用RS256算法和私钥对token进行签名
// 参数:
//   - user: 用户模型对象,包含用户信息
//
// 返回:
//   - string: 生成的JWT字符串
//   - error: 如果生成过程中出现错误则返回error
func GenerateJWT(user model.User) (string, error) {
	// 创建JWT的claims(声明)
	claims := jwt.MapClaims{
		"username": user.Username,                         // 用户名,用于标识token所属用户
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 过期时间,24小时后
	}

	// 使用RS256算法创建token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 使用RSA私钥对token进行签名
	return token.SignedString(privateKey)
}

// ParseJWT 解析JWT并验证其签名
// 使用公钥验证token的签名,并检查签名算法是否为RSA
// 参数:
//   - tokenStr: JWT字符串
//
// 返回:
//   - *jwt.Token: 解析后的token对象
//   - error: 如果解析失败或签名无效则返回error
func ParseJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return publicKey, nil
	})
}

// JWTAuthMiddleware 创建一个Gin中间件用于验证JWT
// 检查请求头中的Authorization字段是否包含有效的JWT
// 如果token无效或已过期,则返回401未授权错误
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// 去除Bearer前缀并验证token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ParseJWT(tokenStr)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// token验证通过,继续处理请求
		c.Next()
	}
}
