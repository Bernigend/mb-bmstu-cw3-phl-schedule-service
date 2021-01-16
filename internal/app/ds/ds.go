package ds

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

// Базовая модель включает в себя общие столбцы всех таблиц
type BaseModel struct {
	UUID      uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(_ *gorm.DB) error {
	b.UUID = uuid.NewV4()
	return nil
}

// Список занятий
type LessonsList []*Lesson

const (
	LessonNameMaxLength  = 128
	LessonWhomMaxLength  = 250
	LessonWhereMaxLength = 32

	LessonStartTimeFormat = "15:04"
	LessonEndTimeFormat   = "15:04"
)

// Модель занятия
type Lesson struct {
	BaseModel
	GroupUUID   uuid.UUID `gorm:"type:uuid;INDEX;not null"`
	Name        string    `gorm:"size:128;not null"`
	Type        int32     `gorm:"type:smallint;not null"`
	Where       string    `gorm:"default:'';size:32;not null"`
	Whom        string    `gorm:"size:250;not null"`
	StartAt     time.Time `gorm:"type:time;not null"`
	EndAt       time.Time `gorm:"type:time;not null"`
	Weekday     int32     `gorm:"type:smallint;not null"`
	IsNumerator bool      `gorm:"not null"`
}
