package ds

import (
	"database/sql"
	"time"
)

type SystemLoad struct {
	SystemLoadID uint         `gorm:"primaryKey;column:system_load_id"`
	Status       string       `gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time    `gorm:"not null"`
	CreatorID    uint         `gorm:"not null"`
	FormingDate  *time.Time   `gorm:"column:forming_date"`
	FinishDate   sql.NullTime `gorm:"column:finish_date"`
	ModeratorID  *uint        `gorm:"column:moderator_id"`
	Description  *string      `gorm:"type:varchar(2000)"`

	Creator   Users  `gorm:"foreignKey:CreatorID"`
	Moderator *Users `gorm:"foreignKey:ModeratorID"`
}

func (SystemLoad) TableName() string {
	return "system_loads"
}
