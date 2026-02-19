package internal

import "context"

type UserService struct {
	Repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (service *UserService) RegisterUser(c context.Context, user UserRegister) (*User, error) {
	return service.Repo.registerUser(c, user)
}

func (service *UserService) LoginUser(c context.Context, email string, password string) (string, error) {
	return service.Repo.loginUser(c, email, password)
}

func (service *UserService) FindMe(ctx context.Context, userId int32) (*User, error) {
	return service.Repo.findMe(ctx, userId)
}
