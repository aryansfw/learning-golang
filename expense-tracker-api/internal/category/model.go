package category

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"not null" json:"name"`
}
