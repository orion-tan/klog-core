package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"

	"golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash 生成密码哈希
// @password 密码
// @return 密码哈希, 错误
func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePasswordHash 比较密码哈希
// @password 密码
// @hash 密码哈希
// @return 是否匹配, 错误
func ComparePasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// CalculateFileHash 计算文件哈希
// @filePath 文件路径
// @return 文件哈希, 错误
func CalculateFileHash(file *multipart.FileHeader) (string, error) {
	hash := md5.New()
	openFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openFile.Close()
	io.Copy(hash, openFile)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
