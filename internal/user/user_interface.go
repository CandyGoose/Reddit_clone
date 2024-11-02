package user

type Service interface {
	Register(username, password string) (User, error)
	Login(username, password string) (User, error)
	GetUserByID(id int) (User, error)
	GetUserByUsername(username string) (User, error)
}
