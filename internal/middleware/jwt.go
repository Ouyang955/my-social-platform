package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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
		// rand.Reader 提供一个安全的随机数生成源
		// 2048 表示RSA密钥的位数,2048位提供足够的安全性
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
	// 在Linux/Unix系统中:
	// - 7 表示 读(4) + 写(2) + 执行(1) 权限
	// - 0 表示 无权限
	// 所以 0700 表示仅文件所有者有完整权限,其他人无权限
	if err := os.MkdirAll("keys", 0700); err != nil {
		return err
	}

	// RSA是一种非对称加密算法,需要一对密钥:
	// - 私钥用于签名和解密
	// - 公钥用于验证签名和加密
	// 这两个密钥必须安全保存,尤其是私钥
	// 	Go原生 RSA 私钥对象 (*rsa.PrivateKey)
	//         ↓（转换为标准二进制）
	// x509.MarshalPKCS1PrivateKey()
	//         ↓
	// PKCS#1 DER 格式（标准二进制）
	//         ↓（转换为可读文本）
	// pem.EncodeToMemory()
	//         ↓
	// PEM 格式（-----BEGIN ...）
	//         ↓
	// 写入 private.pem 文件

	// 将私钥转换为标准的PKCS1格式
	// PKCS1是RSA密钥的标准格式之一
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// PEM是一种文件格式,用于存储密钥、证书等
	// 它会将二进制数据用Base64编码,并加上头尾标记
	// 如 -----BEGIN RSA PRIVATE KEY-----
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// 将私钥写入文件,权限设为600(仅所有者可读写)
	// 私钥必须严格保护,不能让其他用户访问
	if err := ioutil.WriteFile(
		filepath.Join("keys", "private.pem"),
		pem.EncodeToMemory(privateKeyPEM),
		0600,
	); err != nil {
		return err
	}

	// 公钥的处理过程类似
	// 但公钥可以公开,所以权限可以设为644(所有人可读)
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
	// 从项目根目录下的keys/private.pem加载私钥
	// 使用os.ReadFile替代已弃用的ioutil.ReadFile
	privateKeyBytes, err := os.ReadFile(filepath.Join("keys", "private.pem"))
	if err != nil {
		return err // 如果文件不存在或无法读取则返回错误
	}
	// 使用pem.Decode解码PEM格式的私钥数据
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return fmt.Errorf("failed to decode private key PEM block")
	}
	// 将解码后的数据解析为RSA私钥
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// 从项目根目录下的keys/public.pem加载公钥
	publicKeyBytes, err := os.ReadFile(filepath.Join("keys", "public.pem"))
	if err != nil {
		return err // 如果文件不存在或无法读取则返回错误
	}
	// 使用pem.Decode解码PEM格式的公钥数据
	block, _ = pem.Decode(publicKeyBytes)
	if block == nil {
		return fmt.Errorf("failed to decode public key PEM block")
	}
	// 将解码后的数据解析为RSA公钥
	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}

	return nil
}

// GenerateJWT 生成JWT(JSON Web Token)
// JWT包含三部分:
// 1. Header: 包含签名算法和类型
// 2. Payload: 包含用户信息和过期时间
// 3. Signature: 使用私钥对前两部分进行签名
// PayLoad相当于下面的claims

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
		// 这段代码用于验证JWT令牌的签名方法是否为RSA算法
		// token.Method获取令牌使用的签名方法
		// (*jwt.SigningMethodRSA)是类型断言,检查签名方法是否为RSA类型
		// 如果不是RSA类型(即!ok为true),则返回签名无效错误
		// 这是一个安全措施,确保只接受使用预期算法(RSA)签名的令牌
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return publicKey, nil
	})
}

// JWTAuthMiddleware 创建一个Gin中间件用于验证JWT
// 该中间件执行以下操作:
// 1. 从HTTP请求头中获取Authorization字段
// 2. 检查Authorization字段是否存在,不存在则返回401未授权错误
// 3. 从Authorization字段中提取JWT(去除"Bearer "前缀)
// 4. 使用ParseJWT函数验证token的有效性和签名
// 5. 如果token无效或已过期,则返回401未授权错误
// 6. 如果token有效,则允许请求继续处理
//
// 返回:
//   - gin.HandlerFunc: Gin中间件函数,用于集成到路由中
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
