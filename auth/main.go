package main

import (
	"auth/internal/handler"
	"auth/internal/token"
	"auth/storage/postgresql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	r := gin.Default()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	storage := postgresql.NewPostgresql(os.Getenv("DB_CONN_STRING"))
	go storage.DeleteAllExpiredOtps()

	jwtMaker, err := token.NewJwtMaker(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		panic(err)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	seller := r.Group("/seller")
	{
		seller.POST("/registration/sendotp", handler.UserRegistration(storage))
		seller.POST("/registration/submitotp", handler.UserRegistrationSubmitOTP(storage))
		seller.POST("/login", handler.UserLogin(storage, jwtMaker))
		seller.Use(jwtMaker.AuthMiddleware())
		{
			seller.POST("/register/business", handler.RegisterBusiness(storage))
			seller.POST("/register/business/complete", handler.CompleteBusinessRegistration(storage))
			seller.GET("/refresh", handler.RefreshToken(storage, jwtMaker))
			seller.GET("/logout", handler.Logout())
		}
	}

	r.Run(":4005")

}
