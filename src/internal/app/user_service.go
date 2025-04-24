package app

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo UserRepository;
}

func NewUserServiceImpl(repo UserRepository) *UserServiceImpl{
	return &UserServiceImpl{
		repo:repo,
	}
}

// RegisterUser обрабатывает регистрацию нового пользователя
func (s *UserServiceImpl) RegisterUser(req SignUpRequest) (User, error) {

	_, err := s.repo.FindByLogin(req.Login)
	if err == nil {
		return User{}, errors.New("user already exists")
	}
	
	if req.Login == "" || req.Password == "" {
		return User{}, errors.New("login and password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user := User{
		UUID:     uuid.New(),
		Login:    req.Login,
		Password: string(hashedPassword),
	}

	return user, nil
}

// Authenticate проверяет логин и пароль и возвращает UUID пользователя
func (s *UserServiceImpl) Authenticate(login, password string) (uuid.UUID, error) {
	user, err := s.repo.FindByLogin(login)
	if err != nil {
		return uuid.Nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return uuid.Nil, errors.New("unauthorized")
	}

	return user.UUID, nil
}