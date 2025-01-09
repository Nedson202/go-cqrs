package models

// Query models - used for read operations
type UserQuery struct {
    ID string `json:"id"`
}

type UserDTO struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Version  int64  `json:"version"`
} 
