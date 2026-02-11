package main

import (
	"api-gateway/internal/proxy"
	"api-gateway/internal/ratelimiter"
	"api-gateway/internal/token"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("server started")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	jwtMaker := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))
	go ratelimiter.CleanupVisitors()

	r := gin.Default()
	r.Use(ratelimiter.RateLimiter())

	// non-protected services
	nonProtected := r.Group("/api/v1")
	{
		nonProtected.Any("/auth/*path", proxy.ProxyHandler("/api/v1/auth", os.Getenv("AUTH_URL"), jwtMaker))
	}

	//protected services
	api := r.Group("/api/v1")
	api.Use(jwtMaker.JWTAuthMiddleware(jwtMaker))
	{
		api.Any("/product/*path", proxy.ProxyHandler("/api/v1/product", os.Getenv("PRODUCT_URL"), jwtMaker))
		api.Any("/category/*path", proxy.ProxyHandler("/api/v1/category", os.Getenv("CATEGORY_URL"), jwtMaker))
	}

	r.Run(":4000")

}
