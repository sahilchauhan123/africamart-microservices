package service

import (
	"auth/internal/smtp"
	"auth/internal/token"
	"auth/storage"
	"auth/types"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GenerateOtp(min int, max int) int {
	return min + rand.Intn(max-min)
}

func CheckAndSendOTP(db storage.Storage, req types.AccountRegistrationReq) (int, error) {
	exists, err := db.UserExists(req.Email)

	if err != nil {
		return http.StatusInternalServerError, err
	}
	if exists {
		return http.StatusBadRequest, fmt.Errorf("User Already Exists")
	}

	otp := GenerateOtp(1000, 9999)
	err = db.StoreOtp(req, otp)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = smtp.SendEmailOTP(req.Email, strconv.Itoa(otp))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func VerifyAndCreateAccount(db storage.Storage, req types.AccountRegistrationSubmitOTPReq) (int, error) {
	valid, err := db.VerifyOtpAndCreateAccount(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !valid {
		return http.StatusBadRequest, fmt.Errorf("Invalid OTP")
	}
	return http.StatusOK, nil
}

func CreateToken(jwt *token.JWT, userID int64) (string, string, error) {
	accessToken, err := jwt.CreateToken(userID, time.Hour*24*7)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.CreateToken(userID, time.Hour*24*30)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func VerifyUser(db storage.Storage, req types.SellerLoginReq) (int64, error) {
	id, err := db.CheckIdOfEmail(req.Email)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func RegisterBusiness(db storage.Storage, req types.BusinessRegistrationReq, sellerId int64) error {

	err := db.StoreBusiness(sellerId, req)

	if err != nil {
		return err
	}
	return nil
}

func CompleteBusinessRegistration(db storage.Storage, sellerId int64, req types.CompleteBusinessRegistrationReq) error {

	err := db.StoreAllBusinessDetails(sellerId, req)
	if err != nil {
		return err
	}
	return nil
}
