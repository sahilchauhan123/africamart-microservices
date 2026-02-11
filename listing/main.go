package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/test", func(ctx *gin.Context) {
		fmt.Println("its working")
		ctx.JSON(http.StatusOK, gin.H{"message": "its working"})
	})

	r.Run(":4005")
}
