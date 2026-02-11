package postgresql

import (
	"context"
	"log"

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

	query := `
			CREATE TABLE IF NOT EXISTS categories (
				id SERIAL PRIMARY KEY,
				name TEXT NOT NULL,
				slug TEXT UNIQUE,
				parent_id INT REFERENCES categories(id) ON DELETE CASCADE,
				is_leaf BOOLEAN DEFAULT false
			)

			CREATE TABLE products (
				id SERIAL PRIMARY KEY,
				title TEXT NOT NULL,
				description TEXT,
				category_id INT REFERENCES categories(id),
				supplier_id INT REFERENCES suppliers(id),
				min_price NUMERIC,
				max_price NUMERIC,
				rating NUMERIC(2,1),
				total_reviews INT DEFAULT 0,
				latitude DOUBLE PRECISION,
				longitude DOUBLE PRECISION,
				city TEXT,
				created_at TIMESTAMP DEFAULT now()
			)
			
			CREATE TABLE filters (
				id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
				category_id INTEGER NOT NULL,
				name VARCHAR(255) NOT NULL,
				type VARCHAR(50) NOT NULL,    -- string, number, boolean, range
				is_multi BOOLEAN NOT NULL DEFAULT FALSE,
				unit VARCHAR(20),             -- GB, mAh, etc.
				created_at TIMESTAMPTZ DEFAULT NOW(),

				CONSTRAINT fk_category
					FOREIGN KEY(category_id) 
					REFERENCES categories(id)
					ON DELETE CASCADE,

				CONSTRAINT unique_filter_per_category UNIQUE (category_id, name)
			)

			CREATE TABLE product_attributes (
				id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
				product_id UUID NOT NULL,
				filter_id UUID NOT NULL,
				value_string TEXT,
				value_number NUMERIC,
				value_boolean BOOLEAN,

				CONSTRAINT fk_product 
					FOREIGN KEY (product_id) 
					REFERENCES products(id) 
					ON DELETE CASCADE,
					
				CONSTRAINT fk_filter 
					FOREIGN KEY (filter_id) 
					REFERENCES filters(id) 
					ON DELETE CASCADE,

				CONSTRAINT unique_product_filter UNIQUE (product_id, filter_id)
			)

			CREATE TABLE product_images (
				id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
				product_id UUID NOT NULL,
				image_url TEXT NOT NULL,
				position INTEGER NOT NULL DEFAULT 0,

				CONSTRAINT fk_product
					FOREIGN KEY (product_id) 
					REFERENCES products(id) 
					ON DELETE CASCADE,
					
				CONSTRAINT unique_product_image_position UNIQUE (product_id, position)
			)


			`

	_, err = db.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Unable to create table: ", err)
	}

	return &Postgresql{Db: db}
}
