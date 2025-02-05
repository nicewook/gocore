package domain

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserRepository interface {
	Save(user User) error
	FindByID(id int) (User, error)
}

type UserUseCase interface {
	CreateUser(user User) error
	GetUser(id int) (User, error)
}
