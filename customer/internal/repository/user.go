package repository

import (
	"customer/models"
	logger "customer/pkg"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// работа с базкой
// dto - data transfer object. Объект в который парсится результат запрос SQL и из которого он формируется

type user struct { //с маленькой = private; большая - public
	db *sql.DB
}

func NewUser(db *sql.DB) *user {
	// user := user{db: db}
	return &user{db: db}
}

func (r *user) Save(id uuid.UUID, name string, walletAddress string, address string) error {

	sqlStatement := `
		INSERT INTO users (empId, name, walletAddress, address)
		VALUES (?, ?, ?, ?)
		`
	// Prepare the statement for execution
	stmt, err := r.db.Prepare(sqlStatement)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close() // Ensure the prepared statement is closed

	_, err = stmt.Exec(id, name, walletAddress, address)
	if err != nil {
		log.Fatal(err)
		return err
	}

	logger.PrintLog(fmt.Sprintf("inserted user with %s", walletAddress))

	return nil
}

func (r *user) Load(walletAddress string) (models.User, error) { // должен быть models

	return models.User{}, nil
}

// CREATE TABLE CUSTOMERS (
//   empId UUID PRIMARY KEY,
//   name TEXT NOT NULL,
//   walletAddress TEXT NOT NULL,
//   address TEXT NOT NULL
// );
