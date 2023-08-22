package service

type IUserRepository interface {
	Register(string, string) error
	GetById(int) string
	DeleteById(int) error
	LogIn(string, string) string
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) GetById(id int) string {
	return svc.userRepo.GetById(id)
}

func (svc *UserService) Register(email string, pw string) error {
	return svc.userRepo.Register(email, pw)
}

func (svc *UserService) LogIn(email string, pw string) string {
	return svc.userRepo.LogIn(email, pw)
}

func (svc *UserService) DeleteById(id int) error {
	return svc.userRepo.DeleteById(id)
}
