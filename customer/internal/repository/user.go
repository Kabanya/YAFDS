package repository

import (
	"customer/models"
	"customer/pkg"
	"database/sql"

	"github.com/google/uuid"
)

// работа с базкой
// dto - data transfer object. Объект в который парсится результат запрос SQL и из которого он формируется

type UserRepo interface {
	Save(uuid.UUID, string, string, string) error
	Load(walletAddress string) (models.User, error)
}

type userRepo struct { //с маленькой = private; большая - public
	db *sql.DB
}

func NewUser(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Save(id uuid.UUID, name string, walletAddress string, address string) error {
	logger, err := pkg.Logger()
	if err != nil {
		return err
	}

	sqlStatement := `
		INSERT INTO CUSTOMERS (empId, name, walletAddress, address)
		VALUES ($1, $2, $3, $4)
		`
	// Prepare the statement for execution
	stmt, err := r.db.Prepare(sqlStatement)
	if err != nil {
		logger.Printf("Failed to prepare statement: %v", err)
		return err
	}
	defer stmt.Close() // Ensure the prepared statement is closed

	_, err = stmt.Exec(id, name, walletAddress, address)
	if err != nil {
		logger.Printf("Failed to execute insert: %v", err)
		return err
	}

	logger.Printf("Successfully saved customer with ID: %s", id)
	return nil
}

func (r *userRepo) Load(walletAddress string) (models.User, error) {
	logger, err := pkg.Logger()
	if err != nil {
		return models.User{}, err
	}

	sqlStatement := `
		SELECT empId, name, walletAddress, address
		FROM CUSTOMERS
		WHERE walletAddress = $1
		LIMIT 1
	`

	var user models.User
	err = r.db.QueryRow(sqlStatement, walletAddress).Scan(
		&user.Id,
		&user.Name,
		&user.WalletAddress,
		&user.Address,
	)

	if err == sql.ErrNoRows {
		logger.Printf("No customer found with wallet address: %s", walletAddress)
		return models.User{}, err
	}

	if err != nil {
		logger.Printf("Failed to load customer: %v", err)
		return models.User{}, err
	}

	logger.Printf("Successfully loaded customer with wallet address: %s", walletAddress)
	return user, nil
}

// CREATE TABLE CUSTOMERS (
//   empId UUID PRIMARY KEY,
//   name TEXT NOT NULL,
//   walletAddress TEXT NOT NULL,
//   address TEXT NOT NULL
// );
