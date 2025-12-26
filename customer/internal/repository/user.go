package repository

import (
	"customer/models"
	"customer/pkg"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// работа с базкой
// dto - data transfer object. Объект в который парсится результат запрос SQL и из которого он формируется

type UserRepo interface {
	SaveWithPassword(uuid.UUID, string, string, string, string, []byte) error
	LoadByWalletAddress(walletAddress string) (models.User, error)
}

type userRepo struct { //с маленькой = private; большая - public
	db *sql.DB
}

func NewUser(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) SaveWithPassword(id uuid.UUID, name string, walletAddress string, address string, passwordHash string, passwordSalt []byte) error {
	logger, err := pkg.Logger()
	if err != nil {
		return err
	}

	sqlStatement := `
		INSERT INTO CUSTOMERS (emp_id, name, wallet_address, address, password_hash, password_salt)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
	stmt, err := r.db.Prepare(sqlStatement)
	if err != nil {
		logger.Printf("Failed to prepare statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, walletAddress, address, passwordHash, passwordSalt)
	if err != nil {
		logger.Printf("Failed to execute insert: %v", err)
		return err
	}

	logger.Printf("Successfully saved customer with password - ID: %s", id)
	return nil
}

func (r *userRepo) LoadByWalletAddress(walletAddress string) (models.User, error) {
	logger, err := pkg.Logger()
	if err != nil {
		return models.User{}, err
	}

	sqlStatement := `
		SELECT emp_id, name, wallet_address, address, password_hash, password_salt
		FROM CUSTOMERS
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
		&passwordHash,
		&passwordSalt,
	)

	if err == sql.ErrNoRows {
		logger.Printf("No customer found with wallet address: %s", walletAddress)
		return models.User{}, err
	}

	if err != nil {
		logger.Printf("Failed to load customer: %v", err)
		return models.User{}, err
	}

	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
		logger.Printf("Successfully loaded customer with wallet address: %s", walletAddress)
	} else {
		logger.Printf("Password hash is NULL for wallet address: %s", walletAddress)
		return models.User{}, errors.New("password hash is null")
	}
	user.PasswordSalt = passwordSalt

	return user, nil
}
