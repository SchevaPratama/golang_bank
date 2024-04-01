package entity

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func (prod *User) TableName() string {
	return "users"
}
