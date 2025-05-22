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
func (s *UserServiceImpl) LoginJwt(req JwtRequest, jwt JwtProvider) (JwtResponse, error) {
	user, err := s.repo.FindByLogin(req.Login)
	if err != nil {
		return JwtResponse{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return JwtResponse{}, errors.New("unauthorized")
	}
	accessToken, err := jwt.GenerateAccessToken(user)
	if err != nil {
		return JwtResponse{}, errors.New("failed to generate access token")
	}
	refreshToken, err := jwt.GenerateRefreshToken(user)
	if err != nil {
		return JwtResponse{}, errors.New("failed to generate access token")
	}

	return JwtResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		RefreshToken: refreshToken,

	}, nil
}

func (s *UserServiceImpl) RefreshAccessToken(req RefreshJwtRequest) (JwtResponse, error){
	var jwt JwtProvider
	id, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return JwtResponse{}, errors.New("invalid refresh token")
	}
	flag, userinfo := s.repo.FindByUUID(id.String())
	if !flag {
		return JwtResponse{}, errors.New("user not found")
	}
	var user User
	user.UUID = id
	user.Login = userinfo[0]
	user.Password = userinfo[1]
	accessToken, err := jwt.GenerateAccessToken(user)
	if err != nil {
		return JwtResponse{}, errors.New("failed to generate access token")
	}

	return JwtResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		RefreshToken: req.RefreshToken,

	}, nil

}
func (s *UserServiceImpl) RefreshRefreshToken(req RefreshJwtRequest, oldAccessToken string) (JwtResponse, error){
	var jwt JwtProvider
	id, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return JwtResponse{}, errors.New("invalid refresh token")
	}
	flag, userinfo := s.repo.FindByUUID(id.String())
	if !flag {
		return JwtResponse{}, errors.New("user not found")
	}
	var user User
	user.UUID = id
	user.Login = userinfo[0]
	user.Password = userinfo[1]
	refreshToken, err := jwt.GenerateRefreshToken(user)
	if err != nil {
		return JwtResponse{}, errors.New("failed to generate access token")
	}

	return JwtResponse{
		Type:        "Bearer",
		AccessToken: oldAccessToken,
		RefreshToken: refreshToken,

	}, nil
}