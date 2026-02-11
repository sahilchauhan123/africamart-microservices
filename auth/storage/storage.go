package storage

import "auth/types"

type Storage interface {
	UserExists(email string) (bool, error)
	StoreOtp(req types.AccountRegistrationReq, otp int) error
	VerifyOtpAndCreateAccount(req types.AccountRegistrationSubmitOTPReq) (bool, error)
	CheckIdOfEmail(email string) (int64, error)
	StoreBusiness(sellerId int64, req types.BusinessRegistrationReq) error
	StoreAllBusinessDetails(sellerId int64, req types.CompleteBusinessRegistrationReq) error
}
