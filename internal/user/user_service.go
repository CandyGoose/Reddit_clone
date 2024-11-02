package user

import (
	"errors"
	"log"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	mu     sync.Mutex
	users  map[int]User
	nextID int
	logger *log.Logger
}

func NewUserService(logger *log.Logger) Service {
	return &userService{
		users:  make(map[int]User),
		nextID: 1,
		logger: logger,
	}
}

func (s *userService) Register(username, password string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range s.users {
		if u.Username == username {
			return User{}, errors.New("username already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:       s.nextID,
		Username: username,
		Password: string(hashedPassword),
	}
	s.nextID++
	s.users[user.ID] = user

	return user, nil
}

func (s *userService) Login(username, password string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var user User
	found := false
	for _, u := range s.users {
		if u.Username == username {
			user = u
			found = true
			break
		}
	}
	if !found {
		return User{}, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *userService) GetUserByID(id int) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) GetUserByUsername(username string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}
