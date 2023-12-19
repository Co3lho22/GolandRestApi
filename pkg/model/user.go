package model

import "time"

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`        // The '-' tag means this field won't be included in JSON
	HashedPassword string    `json:"-"`               // The '-' tag means this field won't be included in JSON
	Email          string    `json:"email,omitempty"` // omitempty will omit the field if it's empty
	Country        string    `json:"country,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	DateCreated    time.Time `json:"date_created"`
}

// NewUser is a constructor for User struct
func NewUser(username, password, hashedPassword, email, country, phone string) *User {
	return &User{
		Username:       username,
		Password:       password,
		HashedPassword: hashedPassword,
		Email:          email,
		Country:        country,
		Phone:          phone,
		DateCreated:    time.Now(),
	}
}
