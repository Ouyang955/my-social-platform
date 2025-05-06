package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// FileUploadImageHandler 处理图片上传请求
func FileUploadImageHandler(c *gin.Context) {
	// 1. 获取上传文件
	// 从HTTP请求中获取名为"image"的文件字段
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		// 如果获取失败，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传文件"})
		return
	}
	defer file.Close() // 确保函数结束时关闭文件

	// 2. 校验文件类型
	// 获取文件扩展名
	fileExt := filepath.Ext(header.Filename)
	// 定义允许的图片格式
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[fileExt] {
		// 如果不是允许的格式，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持jpg/jpeg/png/webp格式图片"})
		return
	}

	// 3. 校验文件大小（限制为32MB）
	if header.Size > 32*1024*1024 {
		// 如果超过32MB，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过32MB"})
		return
	}

	// 4. 创建保存目录（如果不存在）
	uploadDir := "uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		// 如果创建目录失败，返回500错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 5. 生成唯一文件名（时间戳+原文件名）
	// 使用纳秒级时间戳确保文件名唯一
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
	filePath := filepath.Join(uploadDir, fileName)

	// 6. 保存文件
	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		// 如果创建文件失败，返回500错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存文件"})
		return
	}
	defer dst.Close() // 确保函数结束时关闭文件

	// 将上传的文件内容复制到目标文件
	if _, err = io.Copy(dst, file); err != nil {
		// 如果复制失败，返回500错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 7. 返回文件URL
	// 构建可访问的URL路径
	fileURL := fmt.Sprintf("/uploads/images/%s", fileName)

	// 构建完整URL（开发环境）
	fullURL := fmt.Sprintf("http://localhost:8080%s", fileURL)

	// 返回成功响应和文件URL
	c.JSON(http.StatusOK, gin.H{
		"url":      fileURL, // 相对URL路径
		"full_url": fullURL, // 完整URL
		"message":  "上传成功",
	})
}

// GetImageHandler 获取图片
func GetImageHandler(c *gin.Context) {
	// 1. 获取文件名
	// 从URL参数中获取文件名
	fileName := c.Param("filename")

	// 2. 安全检查：防止目录遍历攻击
	// 确保文件名不包含路径信息，防止访问上级目录
	if filepath.Base(fileName) != fileName {
		// 如果文件名不安全，返回400错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
		return
	}

	// 3. 拼接文件路径
	// 构建完整的文件路径
	filePath := filepath.Join("uploads/images", fileName)

	// 4. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 如果文件不存在，返回404错误
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 5. 返回文件
	// 直接将文件内容发送给客户端
	c.File(filePath)
}
