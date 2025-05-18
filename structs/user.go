package structs

import "time"

type BaseUser struct {
    Name     string `json:"name" binding:"required,min=2,max=100"`
    Username string `json:"username" binding:"required,alphanum,min=3,max=50" gorm:"unique;not null"`
    Email    string `json:"email" binding:"required,email" gorm:"unique;not null"`
}

type UserResponse struct {
    Id        uint      `json:"id"`
    BaseUser
    CreatedAt time.Time `json:"created_at" time_format:"2006-01-02T15:04:05Z"`
    UpdatedAt time.Time `json:"updated_at" time_format:"2006-01-02T15:04:05Z"`
    Token     *string   `json:"token,omitempty"`
}

type UserCreateRequest struct {
    BaseUser
    Password string `json:"password" binding:"required,min=8,containsany=!@#$%^&*()_+"`
}

type UserUpdateRequest struct {
    BaseUser
    Password string `json:"password,omitempty" binding:"omitempty,min=8,containsany=!@#$%^&*()_+"`
}

type UserLoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}