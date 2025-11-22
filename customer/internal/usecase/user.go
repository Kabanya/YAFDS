package usecase

// import (
// 	"customer/internal/service"
// 	modules "customer/modules/user"

// 	"github.com/google/uuid"
// )

// type User interface {
// 	Save(uuid.UUID, string, string, string) error
// 	Load(walletAddress string) (modules.User, error)
// }

// type user struct {
// 	repo service.User
// }

// func NewService(repo service.User) *user {
// 	return &user{repo: repo}
// }

// func (u /*usecase*/ *user) Save(id uuid.UUID, name string, walletAddress string, address string) error {

// 	return u.repo.Save(id, name, walletAddress, address)
// }
