package service

import (
	"context"
	"fmt"
	group_service_api "github.com/Bernigend/mb-cw3-phll-group-service/pkg/group-service-api"
	customErrors "github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/custom-errors"
	"github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/ds"
	api "github.com/Bernigend/mb-cw3-phll-schedule-service/pkg/schedule-service-api"
	uuid "github.com/satori/go.uuid"
	"time"
)

type IRepository interface {
	GetLessonsByGroupUuid(ctx context.Context, groupUuid uuid.UUID) (ds.LessonsList, error)
	AddLesson(ctx context.Context, lesson *ds.Lesson) error
}

type Service struct {
	repository      IRepository
	groupServiceApi group_service_api.GroupServiceClient
}

func NewService(repository IRepository, groupServiceApi group_service_api.GroupServiceClient) *Service {
	return &Service{repository: repository, groupServiceApi: groupServiceApi}
}

// Получение расписания по имени группы
func (s Service) GetScheduleByGroupName(ctx context.Context, groupName string) (*api.GetSchedule_Response, error) {
	group, err := s.groupServiceApi.GetGroup(ctx, &group_service_api.GetGroup_Request{
		GroupName: groupName,
	})
	if err != nil {
		return nil, err
	}

	return s.GetScheduleByGroup(ctx, &api.GetSchedule_GroupItem{
		Uuid:                 group.GetGroupUuid(),
		Name:                 group.GetGroupName(),
		SemesterStart:        group.GetSemesterStartAt(),
		SemesterEnd:          group.GetSemesterEndAt(),
		IsFirstWeekNumerator: group.GetIsFirstWeekNumerator(),
	})
}

// Получение расписания по UUID группы
func (s Service) GetScheduleByGroupUuid(ctx context.Context, groupUuid uuid.UUID) (*api.GetSchedule_Response, error) {
	group, err := s.groupServiceApi.GetGroup(ctx, &group_service_api.GetGroup_Request{
		GroupUuid: groupUuid.String(),
	})
	if err != nil {
		return nil, err
	}

	return s.GetScheduleByGroup(ctx, &api.GetSchedule_GroupItem{
		Uuid:                 group.GetGroupUuid(),
		Name:                 group.GetGroupName(),
		SemesterStart:        group.GetSemesterStartAt(),
		SemesterEnd:          group.GetSemesterEndAt(),
		IsFirstWeekNumerator: group.GetIsFirstWeekNumerator(),
	})
}

// Добавляет занятия в сервис
func (s Service) AddLessons(ctx context.Context, lessonsList []*api.AddLessons_LessonItem) (*api.AddLessons_Response, error) {
	var resultError string

	if len(lessonsList) == 0 {
		return nil, customErrors.InvalidArgument.New(ctx, "список занятий пуст")
	}

	// подготавливаем результат
	results := make([]*api.AddLessons_ResultItem, 0, len(lessonsList))

	for _, lesson := range lessonsList {
		groupUUID := uuid.FromStringOrNil(lesson.GroupUuid)

		if groupUUID == uuid.Nil {
			results = append(results, &api.AddLessons_ResultItem{
				Result: false,
				Error:  "неверный UUID группы",
			})
			continue
		}

		if len(lesson.GetName()) > ds.LessonNameMaxLength {
			results = append(results, &api.AddLessons_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("максимальная длина названия занятия = %d", ds.LessonNameMaxLength),
			})
			continue
		}

		if len(lesson.GetWhere()) > ds.LessonWhereMaxLength {
			results = append(results, &api.AddLessons_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("максимальная длина места проведения занятия = %d", ds.LessonWhereMaxLength),
			})
			continue
		}

		if len(lesson.GetWhom()) > ds.LessonWhomMaxLength {
			results = append(results, &api.AddLessons_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("максимальная длина ФИО преподавателя = %d", ds.LessonWhomMaxLength),
			})
			continue
		}

		startTime, err := time.Parse(ds.LessonStartTimeFormat, lesson.GetStartTime())
		if err != nil {
			results = append(results, &api.AddLessons_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("неверный формат времени начала занятия, требуется: %s и т.п.", ds.LessonStartTimeFormat),
			})
			continue
		}
		startTime = startTime.UTC()

		endTime, err := time.Parse(ds.LessonEndTimeFormat, lesson.GetEndTime())
		if err != nil {
			results = append(results, &api.AddLessons_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("неверный формат времени окончания занятия, требуется: %s и т.п.", ds.LessonEndTimeFormat),
			})
			continue
		}
		endTime = endTime.UTC()

		err = s.repository.AddLesson(ctx, &ds.Lesson{
			GroupUUID:   groupUUID,
			Name:        lesson.GetName(),
			Type:        int32(lesson.GetType()),
			Where:       lesson.GetWhere(),
			Whom:        lesson.GetWhom(),
			StartAt:     startTime,
			EndAt:       endTime,
			Weekday:     int32(lesson.GetWeekday()),
			IsNumerator: lesson.GetIsNumerator(),
		})

		if err != nil {
			resultError = customErrors.ToGRPC(err).Error()
		} else {
			resultError = ""
		}

		results = append(results, &api.AddLessons_ResultItem{
			Result: len(resultError) == 0,
			Error:  resultError,
		})
	}

	return &api.AddLessons_Response{ResultsList: results}, nil
}

// Возвращает расписание группы
func (s Service) GetScheduleByGroup(ctx context.Context, group *api.GetSchedule_GroupItem) (*api.GetSchedule_Response, error) {
	if group == nil {
		return nil, customErrors.InvalidArgument.New(ctx, "группа не найдена")
	}

	// проверяем group UUID на корректность
	groupUUID := uuid.FromStringOrNil(group.Uuid)
	if groupUUID == uuid.Nil {
		return nil, customErrors.InvalidArgument.New(ctx, "неверный UUID группы")
	}

	// получаем список занятий из базы данных
	lessonsList, err := s.repository.GetLessonsByGroupUuid(ctx, groupUUID)
	if err != nil {
		return nil, err
	}

	// подготовленные данные для возврата клиенту
	resultLessonsList := make([]*api.GetSchedule_LessonItem, 0, len(lessonsList))

	for _, lesson := range lessonsList {
		resultLessonsList = append(resultLessonsList, &api.GetSchedule_LessonItem{
			Uuid:        lesson.UUID.String(),
			Name:        lesson.Name,
			Type:        api.LessonType(lesson.Type),
			Where:       lesson.Where,
			Whom:        lesson.Whom,
			StartTime:   lesson.StartAt.UTC().Format(ds.LessonStartTimeFormat),
			EndTime:     lesson.EndAt.UTC().Format(ds.LessonEndTimeFormat),
			Weekday:     api.Weekday(lesson.Weekday),
			IsNumerator: lesson.IsNumerator,
		})
	}

	return &api.GetSchedule_Response{
		LessonsList: resultLessonsList,
		Group:       group,
	}, nil
}
