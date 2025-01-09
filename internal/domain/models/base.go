package models

import "time"

type BaseModel struct {
    Version   int64     `json:"version"`
    UpdatedAt time.Time `json:"updated_at"`
} 
