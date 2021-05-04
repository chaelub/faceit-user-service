package user

import (
	"errors"
	"github.com/chaelub/faceit-user-service/internal/models"
	"regexp"
	"sync"
	"sync/atomic"
)

var (
	ErrUserNotExists = errors.New("user doesn't exist")
)

type (
	InMemoryUserRepo struct {
		usersMu      sync.RWMutex
		users        map[int64]models.User
		usersCounter *int64
	}
)

func (ur *InMemoryUserRepo) New(user models.User) (models.User, error) {
	newId := atomic.AddInt64(ur.usersCounter, 1)
	user.Id = newId
	ur.usersMu.Lock()
	ur.users[newId] = user
	ur.usersMu.Unlock()
	return user, nil
}
func (ur *InMemoryUserRepo) Get(userId int64) (models.User, error) {
	ur.usersMu.RLock()
	user, got := ur.users[userId]
	ur.usersMu.RUnlock()
	if !got {
		return user, ErrUserNotExists
	}
	return user, nil
}
func (ur *InMemoryUserRepo) Update(user models.User) error {
	ur.usersMu.Lock()
	ur.users[user.Id] = user
	ur.usersMu.Unlock()
	return nil
}
func (ur *InMemoryUserRepo) Delete(userId int64) error {
	ur.usersMu.Lock()
	delete(ur.users, userId)
	ur.usersMu.Unlock()
	return nil
}

func (ur *InMemoryUserRepo) Find(criteria map[string][]string) ([]models.User, error) {
	needMatches := len(criteria)
	result := make([]models.User, 0)
	ur.usersMu.RLock()
	for _, user := range ur.users {
		matchCount := 0
		for c, values := range criteria {
			switch c {
			case "country":
				if isPresent(user.Country, values) {
					matchCount++
				}
			case "email":
				for j := range values {
					matched, err := regexp.MatchString(values[j], user.Email)
					if err != nil {
						return nil, err
					}
					if matched {
						matchCount++
						break
					}
				}
			default:
				continue
			}
		}
		if matchCount == needMatches {
			result = append(result, user)
		}
	}
	ur.usersMu.RUnlock()
	return result, nil
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	usersCounter := int64(0)
	return &InMemoryUserRepo{
		usersMu:      sync.RWMutex{},
		users:        make(map[int64]models.User),
		usersCounter: &usersCounter,
	}
}

func isPresent(el string, elms []string) bool {
	for i := range elms {
		if el == elms[i] {
			return true
		}
	}
	return false
}
