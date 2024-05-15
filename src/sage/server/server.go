package server

import (
	"github.com/gin-gonic/gin"
)

func RunServer() error {
	r := gin.Default()

	r.POST("/add", addHandler)
	r.GET("/log", logHandler)
	r.GET("/summary", summaryHandler)
	r.DELETE("/delete/:id", deleteHandler)

	r.Run(":8080")

	return nil
}
