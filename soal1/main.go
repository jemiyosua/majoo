package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	apiVersion := "/api/v1/"

	router.POST(apiVersion+"Login", Login)
	router.POST(apiVersion+"LaporanD", LaporanD)
	router.POST(apiVersion+"LaporanC", LaporanC)

	router.Run(":9000")
}
