package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	MaxFileSize int64    `json:"max_file_size"`
	AllowedExts []string `json:"allowed_exts"`
	UploadPath  string   `json:"upload_path"`
	HttpUlr     string   `json:"http_url"`
	Token       string   `json:"token"`
}

func loadConfig() (*Config, error) {
	// 从配置文件中读取配置使用VIP库
	v := viper.New()
	v.SetConfigFile("config.yaml") // 设置配置文件的路径和名称

	// 可以根据需要设置其他配置选项
	v.SetConfigType("yaml") // 指定配置文件的类型为YAML格式

	// 加载配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(err) // 处理配置文件读取错误
	}
	size := v.GetInt64("server.MaxFileSize")
	path := v.GetString("server.UploadPath")
	url := v.GetString("server.HttpUlr")
	token := v.GetString("server.Token")
	// 加载配置文件，定义各项参数
	return &Config{
		MaxFileSize: size * 1024 * 1024,                        // 文件最大大小
		AllowedExts: []string{".jpg", ".jpeg", ".png", ".gif"}, // 允许的文件类型
		UploadPath:  path,                                      // 图片上传存放目录
		HttpUlr:     url,                                       // 图片访问URL
		Token:       token,
	}, nil
}

func isFileAllowed(fileName string, allowedExts []string) bool {
	ext := filepath.Ext(fileName)

	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}

	return false
}

func handleUpload(c *gin.Context) {
	// 读取配置文件
	config, err := loadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load config"})
		return
	}
	// 添加header token验证
	token := c.Request.Header.Get("Authorization")
	if token != config.Token {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid token"})
		return
	}
	// 获取上传文件和其他表单数据
	file, header, err := c.Request.FormFile("upload_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to fetch file"})
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// 检查文件大小是否超过限制
	if header.Size > config.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File too large"})
		return
	}

	// 检查文件类型是否符合要求
	if !isFileAllowed(header.Filename, config.AllowedExts) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file type"})
		return
	}

	// 根据上传表单的名称生成存储路径和文件名
	fileName := header.Filename
	// 保存文件时将文件名进行哈希，避免文件名重复
	fileName = fmt.Sprintf("%x", md5.Sum([]byte(fileName))) + filepath.Ext(fileName)
	filePath := filepath.Join(config.UploadPath, fileName)

	// 创建一个新文件并将上传数据写入该文件中
	newFile, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create file"})
		return
	}
	defer func(newFile *os.File) {
		_ = newFile.Close()
	}(newFile)

	_, err = io.Copy(newFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save file"})
		return
	}

	// 响应 JSON 数据给客户端
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "url": config.HttpUlr + fileName})
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // 添加 Recovery 中间件
	// 设置静态文件目录
	r.Static("/images", "./images")
	r.POST("/upload", handleUpload)
	err := r.Run(":16001")
	if err != nil {
		return
	}
}
