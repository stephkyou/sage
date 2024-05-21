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

	r.POST("/add", addHandler)
	r.GET("/log", logHandler)
	r.GET("/summary", summaryHandler)
	r.DELETE("/delete/:id", deleteHandler)

	r.Run(":8080")

	return nil
}
