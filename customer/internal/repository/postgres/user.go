package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Kabanya/YAFDS/customer/internal/domain"
	"github.com/Kabanya/YAFDS/pkg/utils"

	"github.com/google/uuid"
)

// работа с базкой
// dto - data transfer object. Объект в который парсится результат запрос SQL и из которого он формируется

type UserRepo interface {
	SaveWithPassword(uuid.UUID, string, string, string, string, []byte) error
	LoadByWalletAddress(walletAddress string) (domain.User, error)
	GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error)
}

type userRepo struct { //с маленькой = private; большая - public
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) SaveWithPassword(id uuid.UUID, name string, walletAddress string, address string, passwordHash string, passwordSalt []byte) error {
	logger, err := utils.Logger()
	if err != nil {
		return err
	}

	sqlStatement := `
		INSERT INTO CUSTOMERS (id, name, wallet_address, address, password_hash, password_salt)
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

func (r *userRepo) LoadByWalletAddress(walletAddress string) (domain.User, error) {
	logger, err := utils.Logger()
	if err != nil {
		return domain.User{}, err
	}

	sqlStatement := `
		SELECT id, name, wallet_address, address, password_hash, password_salt
		FROM CUSTOMERS
		WHERE wallet_address = $1
		LIMIT 1
	`

	var user domain.User
	var passwordHash sql.NullString
	var passwordSalt []byte

	err = r.db.QueryRow(sqlStatement, walletAddress).Scan(
		&user.ID,
		&user.Name,
		&user.WalletAddress,
		&user.Address,
		&passwordHash,
		&passwordSalt,
	)

	if err == sql.ErrNoRows {
		logger.Printf("No customer found with wallet address: %s", walletAddress)
		return domain.User{}, err
	}

	if err != nil {
		logger.Printf("Failed to load customer: %v", err)
		return domain.User{}, err
	}

	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
		logger.Printf("Successfully loaded customer with wallet address: %s", walletAddress)
	} else {
		logger.Printf("Password hash is NULL for wallet address: %s", walletAddress)
		return domain.User{}, errors.New("password hash is null")
	}
	user.PasswordSalt = passwordSalt

	return user, nil
}

func (r *userRepo) GetCustomerWalletAddress(ctx context.Context, customerID uuid.UUID) (string, error) {
	if customerID == uuid.Nil {
		return "", errors.New("customer_id must be a valid UUID")
	}

	var walletAddress string
	err := r.db.QueryRowContext(ctx, "SELECT wallet_address FROM CUSTOMERS WHERE id = $1", customerID).Scan(&walletAddress)
	if err == sql.ErrNoRows {
		return "", errors.New("customer not found")
	}
	if err != nil {
		return "", err
	}
	return walletAddress, nil
}
