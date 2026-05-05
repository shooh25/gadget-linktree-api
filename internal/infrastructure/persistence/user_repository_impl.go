package persistence

import (
	"gadget-linktree-api/internal/domain/model/user"
	"sync"
)

type UserRepositoryImpl struct {
	mu    sync.RWMutex
	users map[string]user.User // key: UserID
}

func NewUserRepositoryImpl() *UserRepositoryImpl {
	return &UserRepositoryImpl{
		users: make(map[string]user.User),
	}
}

func (r *UserRepositoryImpl) ExistsByGoogleId(googleId string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.GoogleId() == googleId {
			return true, nil
		}
	}
	return false, nil
}

func (r *UserRepositoryImpl) Save(u user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.UserId().Value()] = u
	return nil
}

func (r *UserRepositoryImpl) FindByUserId(userId user.UserId) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[userId.Value()]
	if !ok {
		return nil, nil
	}
	return &u, nil
}

func (r *UserRepositoryImpl) FindByGoogleId(googleId string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.GoogleId() == googleId {
			return &u, nil
		}
	}
	return nil, nil
}
