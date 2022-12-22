package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" gorm:"type: varchar(255)"`
	Password string `json:"password" gorm:"type: varchar(255)"`
	FullName string `json:"fullname" gorm:"type: varchar(255)"`
	Phone    string `json:"phone" gorm:"type: varchar(255)"`
	Image    string  `json:"image" gorm:"type: varchar(255)"`
}

type UsersProfileResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"name"`
	Location string `json:"location"`
	Email    string `json:"email"`
	Image    string `json:"image"`
}

func (UsersProfileResponse) TableName() string {
	return "users"
}
