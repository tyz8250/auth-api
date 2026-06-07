package repository

import (
	"auth-api/model"
	"sync"
	"time"
)

type MemoryUserRepository struct {
	mu     sync.Mutex
	users  map[int]model.User
	nextID int
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:  make(map[int]model.User),
		nextID: 1,
	}
}

func (r *MemoryUserRepository) Create(user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	user.ID = r.nextID
	user.CreatedAt = now
	user.UpdatedAt = now

	r.users[user.ID] = user
	r.nextID++

	return user, nil
}

func (r *MemoryUserRepository) FindByID(id int) (model.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	return user, ok
}
