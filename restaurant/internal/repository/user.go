package repository

import (
	"database/sql"
	"errors"
	"restaurant/models"

	"customer/pkg"

	"github.com/google/uuid"
)

// работа с базкой
// dto - data transfer object. Объект в который парсится результат запрос SQL и из которого он формируется

type UserRepo interface {
	SaveWithPassword(id uuid.UUID, name string, walletAddress string, address string, isActive bool, passwordHash string, passwordSalt []byte) error
	LoadByWalletAddress(walletAddress string) (models.User, error)
}

type userRepo struct { //с маленькой = private; большая - public
	db *sql.DB
}

func NewUser(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) SaveWithPassword(id uuid.UUID, name string, walletAddress string, address string, isActive bool, passwordHash string, passwordSalt []byte) error {
	logger, err := pkg.Logger()
	if err != nil {
		return err
	}

	sqlStatement := `
		INSERT INTO RESTAURANTS (emp_id, name, wallet_address, address, status, password_hash, password_salt)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
	stmt, err := r.db.Prepare(sqlStatement)
	if err != nil {
		logger.Printf("Failed to prepare statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, walletAddress, address, isActive, passwordHash, passwordSalt)
	if err != nil {
		logger.Printf("Failed to execute insert: %v", err)
		return err
	}

	logger.Printf("Successfully saved restaurant with password - ID: %s", id)
	return nil
}

func (r *userRepo) LoadByWalletAddress(walletAddress string) (models.User, error) {
	logger, err := pkg.Logger()
	if err != nil {
		return models.User{}, err
	}

	sqlStatement := `
		SELECT emp_id, name, wallet_address, address, status, password_hash, password_salt
		FROM RESTAURANTS
		WHERE wallet_address = $1
		LIMIT 1
	`

	var user models.User
	var passwordHash sql.NullString
	var passwordSalt []byte

	err = r.db.QueryRow(sqlStatement, walletAddress).Scan(
		&user.Id,
		&user.Name,
		&user.WalletAddress,
		&user.Address,
		&user.IsActive,
		&passwordHash,
		&passwordSalt,
	)

	if err == sql.ErrNoRows {
		logger.Printf("No restaurant found with wallet address: %s", walletAddress)
		return models.User{}, err
	}

	if err != nil {
		logger.Printf("Failed to load restaurant: %v", err)
		return models.User{}, err
	}

	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
		logger.Printf("Successfully loaded restaurant with wallet address: %s", walletAddress)
	} else {
		logger.Printf("Password hash is NULL for wallet address: %s", walletAddress)
		return models.User{}, errors.New("password hash is null")
	}
	user.PasswordSalt = passwordSalt

	return user, nil
}

// CREATE TABLE restaurantS (
//   empId UUID PRIMARY KEY,
//   name TEXT NOT NULL,
//   walletAddress TEXT NOT NULL,
//   address TEXT NOT NULL
// );
