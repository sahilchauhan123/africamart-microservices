package postgresql

import (
	"auth/types"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgresql struct {
	Db *pgxpool.Pool
}

func NewPostgresql(connString string) *Postgresql {

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	query := `CREATE TABLE IF NOT EXISTS sellers(
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			business_owner_name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`

	_, err = db.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Unable to create table: ", err)
	}

	// only 1 record per email at a time
	query = `CREATE TABLE IF NOT EXISTS otps(
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			otp INT NOT NULL,
			password VARCHAR(255) NOT NULL,
			business_owner_name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP + INTERVAL '5 minutes',
			CONSTRAINT unique_email UNIQUE(email)
		);`

	query = `CREATE TABLE IF NOT EXISTS businesses(
			id SERIAL PRIMARY KEY,
			seller_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(255) NOT NULL,
			address TEXT NOT NULL,
			category VARCHAR(255) NOT NULL,
			country VARCHAR(255),
			city VARCHAR(255),
			state VARCHAR(255),
			logo_url VARCHAR(255),
			phone_number VARCHAR(255),
			email VARCHAR(255) UNIQUE,
			no_of_employees INT,
			year_of_establishment INT,
			documents_url TEXT[],
			verified BOOLEAN DEFAULT false,
			verified_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`

	_, err = db.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Unable to create table: ", err)
	}

	fmt.Println("Database connection successful")
	return &Postgresql{Db: db}

}

func (p *Postgresql) UserExists(email string) (bool, error) {
	var exists bool
	err := p.Db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM sellers WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *Postgresql) StoreOtp(req types.AccountRegistrationReq, otp int) error {
	// return id of inserted row
	var id int
	query := `INSERT INTO otps (email, otp, password, business_owner_name) VALUES ($1, $2, $3, $4) RETURNING id`
	err := p.Db.QueryRow(context.Background(), query, req.Email, otp, req.Password, req.BusinessOwnerName).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("Otp Already Sent Please Wait for 5 minutes")
			}
		}
		return err
	}

	// go routine to delete otp after 5 minutes
	go func() {
		time.Sleep(5 * time.Minute)
		if id == 0 {
			return
		}
		_, err := p.Db.Exec(context.Background(), "DELETE FROM otps WHERE id = $1", id)
		if err != nil {
			log.Println("Unable to delete otp: ", err)
		}
	}()

	return nil
}

func (p *Postgresql) VerifyOtpAndCreateAccount(req types.AccountRegistrationSubmitOTPReq) (bool, error) {
	// check if otp is valid and not expired
	var email string
	var business_owner_name string
	var otp int
	var password string
	var expires_at time.Time
	var created_at time.Time

	err := p.Db.QueryRow(context.Background(), "SELECT email,business_owner_name,otp,password,expires_at,created_at FROM otps WHERE email = $1 AND otp = $2",
		req.Email, req.OTP).Scan(&email, &business_owner_name, &otp, &password, &expires_at, &created_at)

	if err != nil {
		return false, err
	}
	if expires_at.Before(time.Now()) {
		return false, fmt.Errorf("OTP has expired")
	}

	_, err = p.Db.Exec(context.Background(), "INSERT INTO sellers (email, password, business_owner_name) VALUES ($1, $2, $3)", email, password, business_owner_name)
	if err != nil {
		return false, err
	}

	go p.Db.Exec(context.Background(), "DELETE FROM otps WHERE email = $1", email)

	return true, nil
}

func (p *Postgresql) DeleteAllExpiredOtps() {
	_, err := p.Db.Exec(context.Background(), "DELETE FROM otps WHERE expires_at < NOW()")
	if err != nil {
		log.Println("Unable to delete expired otps: ", err)
	}
}

func (p *Postgresql) CheckIdOfEmail(email string) (int64, error) {
	var id int64
	err := p.Db.QueryRow(context.Background(), "SELECT id FROM sellers WHERE email = $1", email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *Postgresql) StoreBusiness(sellerId int64, req types.BusinessRegistrationReq) error {

	query := `INSERT INTO businesses 
	(seller_id, name, type, address, category) 
	VALUES ($1, $2, $3, $4, $5)`

	_, err := p.Db.Exec(context.Background(), query, sellerId, req.Name, req.Type, req.Address, req.Category)
	if err != nil {
		return err
	}
	return nil

}

func (p *Postgresql) StoreAllBusinessDetails(sellerId int64, req types.CompleteBusinessRegistrationReq) error {

	query := `INSERT INTO businesses 
	(seller_id, name, type, address, category, country, city, state, 
	logo_url, phone_number, email, no_of_employees, year_of_establishment,
	documents_url) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := p.Db.Exec(context.Background(), query, sellerId,
		req.Name, req.Type, req.Address, req.Category, req.Country, req.City, req.State,
		req.LogoUrl, req.PhoneNumber, req.Email, req.NoOfEmployees, req.YearOfEstablishment,
		req.DocumentsUrl)
	if err != nil {
		return err
	}
	return nil

}
