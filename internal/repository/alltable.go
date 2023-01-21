package repository

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const sellerTable = `CREATE TABLE IF NOT EXISTS sellers (
	id INTEGER PRIMARY KEY AUTO_INCREMENT UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password TEXT NOT NULL
	)`

const clientTable = `CREATE TABLE IF NOT EXISTS clients (
	id INTEGER PRIMARY KEY AUTO_INCREMENT UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password TEXT NOT NULL
);`

const productTable = `CREATE TABLE IF NOT EXISTS product (
	id INTEGER PRIMARY KEY AUTO_INCREMENT UNIQUE NOT NULL, 
	seller_id INTEGER NOT NULL,
	company TEXT NOT NULL,
	description TEXT NOT NULL, 
	price FLOAT NOT NULL
);`

// const JWTTable = `CREATE TABLE IF NOT EXISTS "tokens"  (
// "id" INTEGER PRIMARY KEY AUTO_INCREMENT UNIQUE NOT NULL,
// "seller_id" INTEGER NOT NULL,
// "signingkey" TEXT NOT NULL,
// "date" DATETIME DEFAULT CURRENT_TIMESTAMP
// );
// `

var tables = []string{sellerTable, clientTable, productTable}

func Init() (*sql.DB, error) {
	var err error
	db, err := sql.Open("mysql", "root:root@tcp(database:3306)/shopdb")
	if err != nil {
		log.Println("❌ error | can't open DB")
		return nil, err
	}
	return db, nil
}

func CreateDatabase(db *sql.DB) error {
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}
