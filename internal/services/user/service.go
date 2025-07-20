package user

import "github.com/andrew-womeldorf/connect-boilerplate/internal/services/user/store"

// Service handles the business logic
type Service struct {
	store store.Store
}

func NewService(store store.Store) *Service {
	return &Service{store}
}
