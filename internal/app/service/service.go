package service

import (
	"context"
	customErrors "github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/custom-errors"
	"github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/ds"
	api "github.com/Bernigend/mb-cw3-phll-schedule-service/pkg/schedule-service-api"
	uuid "github.com/satori/go.uuid"
	"time"
)

type IRepository interface {
	GetLessonsByGroupUuid(groupUuid uuid.UUID) (ds.LessonsList, error)
	AddLesson(lesson *ds.Lesson) error
}

type Service struct {
	repository      IRepository
	groupServiceApi interface{} // TODO: подключение group service API
}

// Получение расписания по имени группы
func (s Service) GetScheduleByGroupName(ctx context.Context, groupName string) (*api.GetSchedule_Response, error) {
	// TODO: получение данных группы из group service
	group := &api.GetSchedule_GroupItem{
		Uuid:                 "",
		Name:                 groupName,
		SemesterStart:        nil,
		SemesterEnd:          nil,
		IsFirstWeekNumerator: false,
	}

	return s.GetScheduleByGroup(ctx, group)
}

// Получение расписания по UUID группы
func (s Service) GetScheduleByGroupUuid(ctx context.Context, groupUuid uuid.UUID) (*api.GetSchedule_Response, error) {
	// TODO: получение данных группы из group service
	group := &api.GetSchedule_GroupItem{
		Uuid:                 groupUuid.String(),
		Name:                 "",
		SemesterStart:        nil,
		SemesterEnd:          nil,
		IsFirstWeekNumerator: false,
	}

	return s.GetScheduleByGroup(ctx, group)
}

// Добавляет занятия в сервис
func (s Service) AddLessons(ctx context.Context, lessonsList []*api.AddLessons_LessonItem) (*api.AddLessons_Response, error) {
	if lessonsList == nil {
		return nil, customErrors.InvalidArgument.New(ctx, "lessons list is empty")
	}

	// подготавливаем результат
	results := make([]*api.AddLessons_ResultItem, len(lessonsList))

	for _, lesson := range lessonsList {
		groupUUID := uuid.FromStringOrNil(lesson.GroupUuid)

		if groupUUID == uuid.Nil {
			results = append(results, &api.AddLessons_ResultItem{
				Result: false,
				Error:  "group UUID is invalid",
			})
			continue
		}

		err := s.repository.AddLesson(&ds.Lesson{
			GroupUUID:   groupUUID,
			Name:        "",
			Type:        0,
			Where:       "",
			Whom:        "",
			StartAt:     time.Time{},
			EndAt:       time.Time{},
			Weekday:     0,
			IsNumerator: false,
		})

		if err != nil {
			results = append(results, &api.AddLessons_ResultItem{
				Result: false,
				Error:  customErrors.ToGRPC(err).Error(),
			})
			continue
		}
	}

	return &api.AddLessons_Response{ResultsList: results}, nil
}

// Возвращает расписание группы
func (s Service) GetScheduleByGroup(ctx context.Context, group *api.GetSchedule_GroupItem) (*api.GetSchedule_Response, error) {
	if group == nil {
		return nil, customErrors.InvalidArgument.New(ctx, "group not found")
	}

	// проверяем group UUID на корректность
	groupUUID := uuid.FromStringOrNil(group.Uuid)
	if groupUUID == uuid.Nil {
		return nil, customErrors.InvalidArgument.New(ctx, "invalid group uuid")
	}

	// получаем список занятий из базы данных
	lessonsList, err := s.repository.GetLessonsByGroupUuid(groupUUID)
	if err != nil {
		return nil, err
	}

	// подготовленные данные для возврата клиенту
	resultLessonsList := make([]*api.GetSchedule_LessonItem, len(lessonsList))

	for _, lesson := range lessonsList {
		resultLessonsList = append(resultLessonsList, &api.GetSchedule_LessonItem{
			Uuid:        lesson.UUID.String(),
			Name:        lesson.Name,
			Type:        api.LessonType(lesson.Type),
			Where:       lesson.Where,
			Whom:        lesson.Whom,
			StartTime:   lesson.StartAt.String(),
			EndTime:     lesson.EndAt.String(),
			Weekday:     api.Weekday(lesson.Weekday),
			IsNumerator: lesson.IsNumerator,
		})
	}

	return &api.GetSchedule_Response{
		LessonsList: resultLessonsList,
		Group:       group,
	}, nil
}
