package endpoint

import (
	"context"
	customErrors "github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/custom-errors"
	api "github.com/Bernigend/mb-cw3-phll-schedule-service/pkg/schedule-service-api"
	uuid "github.com/satori/go.uuid"
)

// Необходимые для работы API методы сервиса
type IService interface {
	// получение расписания группы
	GetScheduleByGroupName(ctx context.Context, groupName string) (*api.GetSchedule_Response, error)
	GetScheduleByGroupUuid(ctx context.Context, groupUuid uuid.UUID) (*api.GetSchedule_Response, error)

	// добавление занятий
	AddLessons(ctx context.Context, lessonsList []*api.AddLessons_LessonItem) (*api.AddLessons_Response, error)
}

type Endpoint struct {
	service IService
}

func NewEndpoint(service IService) *Endpoint {
	return &Endpoint{service: service}
}

// [Метод API] Возвращает расписание группы
func (e Endpoint) GetSchedule(ctx context.Context, request *api.GetSchedule_Request) (*api.GetSchedule_Response, error) {
	if groupUuid := request.GetGroupUuid(); len(groupUuid) > 0 {
		parsedGroupUuid, err := uuid.FromString(groupUuid)
		if err != nil {
			return nil, customErrors.InvalidArgument.New(ctx, "invalid group UUID")
		}

		result, err := e.service.GetScheduleByGroupUuid(ctx, parsedGroupUuid)
		return result, customErrors.ToGRPC(err)
	}

	if groupName := request.GetGroupName(); len(groupName) > 0 {
		result, err := e.service.GetScheduleByGroupName(ctx, groupName)
		return result, customErrors.ToGRPC(err)
	}

	return nil, customErrors.InvalidArgument.New(ctx, "expected group name or group uuid")
}

// [Метод API] Добавляет занятия в систему
func (e Endpoint) AddLessons(ctx context.Context, request *api.AddLessons_Request) (*api.AddLessons_Response, error) {
	if lessonsList := request.GetLessonsList(); lessonsList != nil {
		result, err := e.service.AddLessons(ctx, lessonsList)
		return result, customErrors.ToGRPC(err)
	}

	return nil, customErrors.InvalidArgument.New(ctx, "expected lessons list")
}
