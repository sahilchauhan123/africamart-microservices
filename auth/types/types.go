package types

type AccountRegistrationReq struct {
	Email             string `json:"email" binding:"required"`
	Password          string `json:"password" binding:"required"`
	BusinessOwnerName string `json:"business_owner_name" binding:"required"`
}

type AccountRegistrationSubmitOTPReq struct {
	Email string `json:"email" binding:"required"`
	OTP   int    `json:"otp" binding:"required"`
}

type SellerLoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SellerLoginRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type BusinessRegistrationReq struct {
	Name     string `json:"business_name" binding:"required"`
	Type     string `json:"business_type" binding:"required"`
	Address  string `json:"business_address" binding:"required"`
	Category string `json:"business_category" binding:"required"`
}

type CompleteBusinessRegistrationReq struct {
	Name                string   `json:"business_name" binding:"required"`
	Type                string   `json:"business_type" binding:"required"`
	Address             string   `json:"business_address" binding:"required"`
	Category            string   `json:"business_category" binding:"required"`
	Country             string   `json:"country"`
	City                string   `json:"city"`
	State               string   `json:"state"`
	LogoUrl             string   `json:"logo_url"`
	PhoneNumber         string   `json:"phone_number"`
	Email               string   `json:"email"`
	NoOfEmployees       int      `json:"no_of_employees"`
	YearOfEstablishment int      `json:"year_of_establishment"`
	DocumentsUrl        []string `json:"documents_url"`
}

// CREATE TABLE IF NOT EXISTS businesses(
// 			id SERIAL PRIMARY KEY,
// 			seller_id INT NOT NULL,
// 			name VARCHAR(255) NOT NULL,
// 			type VARCHAR(255) NOT NULL,
// 			address TEXT NOT NULL,
// 			category VARCHAR(255) NOT NULL,
// 			country VARCHAR(255),
// 			city VARCHAR(255),
// 			state VARCHAR(255),
// 			logo_url VARCHAR(255),
// 			phone_number VARCHAR(255),
// 			email VARCHAR(255) UNIQUE,
// 			no_of_employees INT,
// 			year_of_establishment INT,
// 			documents_url TEXT[],
// 			verified BOOLEAN DEFAULT false,
// 			verified_at TIMESTAMP WITH TIME ZONE,
// 			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
// 		)
