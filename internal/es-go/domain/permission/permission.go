package permission

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type Permission struct {
	Id          ulid.ULID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
}
