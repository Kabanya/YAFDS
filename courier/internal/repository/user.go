package repository

import (
	"courier/models"
	"database/sql"
	"errors"

	"customer/pkg"

	"github.com/google/uuid"
)

// работа с базкой
// dto - data transfer object. Объект в который парсится результат запрос SQL и из которого он формируется

type UserRepo interface {
	SaveWithPassword(id uuid.UUID, name string, walletAddress string, transportType string, passwordHash string, passwordSalt []byte) error
	LoadByWalletAddress(walletAddress string) (models.User, error)
}

type userRepo struct { //с маленькой = private; большая - public
	db *sql.DB
}

func NewUser(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) SaveWithPassword(id uuid.UUID, name string, walletAddress string, transportType string, passwordHash string, passwordSalt []byte) error {
	logger, err := pkg.Logger()
	if err != nil {
		return err
	}

	sqlStatement := `
		INSERT INTO COURIERS (empId, name, walletAddress, transport_type, is_active, geolocation, password_hash, password_salt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
	stmt, err := r.db.Prepare(sqlStatement)
	if err != nil {
		logger.Printf("Failed to prepare statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, walletAddress, transportType, true, "0,0", passwordHash, passwordSalt)
	if err != nil {
		logger.Printf("Failed to execute insert: %v", err)
		return err
	}

	logger.Printf("Successfully saved courier with password - ID: %s", id)
	return nil
}

func (r *userRepo) LoadByWalletAddress(walletAddress string) (models.User, error) {
	logger, err := pkg.Logger()
	if err != nil {
		return models.User{}, err
	}

	sqlStatement := `
		SELECT empId, name, walletAddress, transport_type, is_active, geolocation, password_hash, password_salt
		FROM COURIERS
		WHERE walletAddress = $1
		LIMIT 1
	`

	var user models.User
	var passwordHash sql.NullString
	var passwordSalt []byte

	err = r.db.QueryRow(sqlStatement, walletAddress).Scan(
		&user.Id,
		&user.Name,
		&user.WalletAddress,
		&user.TransportType,
		&user.IsActive,
		&user.Geolocation,
		&passwordHash,
		&passwordSalt,
	)

	if err == sql.ErrNoRows {
		logger.Printf("No courier found with wallet address: %s", walletAddress)
		return models.User{}, err
	}

	if err != nil {
		logger.Printf("Failed to load courier: %v", err)
		return models.User{}, err
	}

	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
		logger.Printf("Successfully loaded courier with wallet address: %s", walletAddress)
	} else {
		logger.Printf("Password hash is NULL for wallet address: %s", walletAddress)
		return models.User{}, errors.New("password hash is null")
	}
	user.PasswordSalt = passwordSalt

	return user, nil
}

// CREATE TABLE courierS (
//   empId UUID PRIMARY KEY,
//   name TEXT NOT NULL,
//   walletAddress TEXT NOT NULL,
//   address TEXT NOT NULL
// );
