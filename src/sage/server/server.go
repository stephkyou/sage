package server

import (
	"database/sql"
	"sage/src/sage/cmd"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func RunServer() error {
	db, _ = cmd.ConnectDB(cmd.SAGE_DB_NAME)
	r := gin.Default()
	r.Use(corsMiddleware())

	r.POST("/add", addHandler)
	r.GET("/log", logHandler)
	r.GET("/summary", summaryHandler)
	r.DELETE("/delete/:id", deleteHandler)
	r.GET("/count", countHandler)

	r.Run(":8080")

	return nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
