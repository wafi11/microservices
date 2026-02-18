package internal

type UserService struct {
	Repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (service *UserService) RegisterUser(user UserRegister) (User, error) {
	return service.Repo.registerUser(user)
}
