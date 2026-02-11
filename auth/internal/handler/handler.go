package handler

import (
	"auth/internal/service"
	"auth/internal/token"
	"auth/response"
	"auth/storage"
	"auth/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// check if user exists
// if user exists, return error
// if user does not exist, send otp to mail
// return otp

func UserRegistration(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {

		var request types.AccountRegistrationReq
		err := c.ShouldBindBodyWithJSON(&request)
		if err != nil {
			response.Error(c, err, http.StatusBadRequest)
			return
		}

		statusCode, err := service.CheckAndSendOTP(storage, request)
		if err != nil {
			response.Error(c, err, statusCode)
			return
		}

		response.Success(c, "User registered successfully", statusCode)
	}
}

func UserRegistrationSubmitOTP(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.AccountRegistrationSubmitOTPReq
		err := c.ShouldBindBodyWithJSON(&request)
		if err != nil {
			response.Error(c, err, http.StatusBadRequest)
			return
		}

		statusCode, err := service.VerifyAndCreateAccount(storage, request)
		if err != nil {
			response.Error(c, err, statusCode)
			return
		}

		response.Success(c, "User registered successfully", statusCode)
	}
}

func UserLogin(storage storage.Storage, jwt *token.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.SellerLoginReq
		var resp types.SellerLoginRes
		err := c.ShouldBindBodyWithJSON(&request)
		if err != nil {
			response.Error(c, err, http.StatusBadRequest)
			return
		}

		userID, err := service.VerifyUser(storage, request)
		if err != nil {
			response.Error(c, err, http.StatusBadRequest)
			return
		}
		fmt.Println("User ID: ", userID)
		resp.AccessToken, resp.RefreshToken, err = service.CreateToken(jwt, userID)
		if err != nil {
			response.Error(c, err, http.StatusInternalServerError)
			return
		}
		c.SetCookie("access_token", resp.AccessToken, 60*60*24*7, "/", "", true, true)
		c.SetCookie("refresh_token", resp.RefreshToken, 60*60*24*7, "/", "", true, true)
		response.Success(c, "Login successfull", http.StatusOK)
	}
}
func RefreshToken(storage storage.Storage, jwt *token.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie("access_token", "", -1, "/", "", true, true)
		ctx.SetCookie("refresh_token", "", -1, "/", "", true, true)
		response.Success(ctx, "Logout successfull", http.StatusOK)
	}
}

func RegisterBusiness(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.BusinessRegistrationReq
		err := c.ShouldBindBodyWithJSON(&request)
		if err != nil {
			response.Error(c, err, http.StatusBadRequest)
			return
		}
		sellerId := c.GetInt64("seller_id")
		err = service.RegisterBusiness(storage, request, sellerId)
		if err != nil {
			response.Error(c, err, http.StatusBadRequest)
			return
		}

		response.Success(c, "Business registered successfully", http.StatusOK)
	}
}

func CompleteBusinessRegistration(storage storage.Storage) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var request types.CompleteBusinessRegistrationReq
		err := ctx.ShouldBindBodyWithJSON(&request)
		if err != nil {
			response.Error(ctx, err, http.StatusBadRequest)
			return
		}

		userId := ctx.GetInt64("user_id")
		err = service.CompleteBusinessRegistration(storage, userId, request)
		if err != nil {
			response.Error(ctx, err, http.StatusBadRequest)
			return
		}

		response.Success(ctx, "Business registered successfully", http.StatusOK)
	}
}
