package repository

import (
	"context"
	"errors"
	customErrors "github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/custom-errors"
	"github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/ds"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// Создаёт новое подключение к базе данных
func NewRepository(dsn string) (*Repository, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Repository{db: conn}, nil
}

// Закрывает соединение с базой данных, если оно было установлено
func (r Repository) Close() error {
	if r.db == nil {
		return nil
	}

	db, err := r.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

// Выполняет автоматические миграции в базе данных
func (r Repository) AutoMigrate() error {
	if r.db == nil {
		return nil
	}

	return r.db.AutoMigrate(&ds.Lesson{})
}

// Возвращает список занятий по UUID группы
func (r Repository) GetLessonsByGroupUuid(ctx context.Context, groupUuid uuid.UUID) (ds.LessonsList, error) {
	var lessons []*ds.Lesson

	resultError := r.db.Where(&ds.Lesson{
		GroupUUID: groupUuid,
	}).Find(&lessons).Error

	if resultError != nil {
		if errors.Is(resultError, gorm.ErrRecordNotFound) {
			return nil, customErrors.NotFound.New(ctx, "занятия данной группы не найдены")
		} else {
			return nil, customErrors.Internal.NewWrap(ctx, "произошла непредвиденная ошибка", resultError)
		}
	}

	return lessons, nil
}

func (r Repository) AddLesson(_ context.Context, lesson *ds.Lesson) error {
	return r.db.Create(lesson).Error
}
