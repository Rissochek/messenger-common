package model

type User struct {
	UserId   int       `db:"id"`
	Username string    `json:"username" db:"username" binding:"required,min=4"`
	Password string    `json:"password" db:"password" binding:"required,min=8"`
	Info     *UserInfo `json:"profile" binding:"omitempty"`
	Groups   []int
}

type UserInfo struct {
	Name    string `json:"name" db:"name" binding:"required,min=2"`
	Surname string `json:"surname" db:"surname" binding:"required,min=2"`
}
