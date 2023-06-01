package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func ValidateHeader(requestBody string, c *gin.Context) bool {
	var (
		result gin.H
	)

	Signature := c.GetHeader("Signature")
	ContentType := c.GetHeader("Content-Type")

	enckey := os.Getenv("ENCKEY")
	base64String := base64.StdEncoding.EncodeToString([]byte(requestBody))
	SignatureKey := fmt.Sprintf("%x", md5.Sum([]byte(enckey+base64String)))

	if ContentType == "" {
		errorMessage := "Error, Header - Content-Type is not application/json or empty value "
		result = gin.H{
			"ErrCode":    "1",
			"ErrMessage": errorMessage,
			"Result":     "",
		}
		c.JSON(http.StatusOK, result)
		return false
	}

	if Signature == "" {
		errorMessage := "Header Signature can not null"
		result = gin.H{
			"ErrCode":    "1",
			"ErrMessage": errorMessage,
			"Result":     "",
		}
		c.JSON(http.StatusOK, result)
		return false
	} else {
		if Signature == SignatureKey {
			return true
		} else {
			result = gin.H{
				"ErrCode":    "1",
				"ErrMessage": "Header Signature invalid " + Signature + "==" + SignatureKey,
				"Result":     "",
			}
			c.JSON(http.StatusOK, result)
			return false
		}
	}
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
