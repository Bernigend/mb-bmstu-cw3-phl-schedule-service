package endpoint

import (
	"context"
	api "github.com/Bernigend/mb-cw3-phll-schedule-service/pkg/schedule-service-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IService interface {
	// получение расписания группы
	GetScheduleByGroupName(ctx context.Context, groupName string) (*api.GetSchedule_Response, error)
	GetScheduleByGroupUuid(ctx context.Context, groupUuid string) (*api.GetSchedule_Response, error)

	// добавление занятий
	AddLessons(ctx context.Context, lessonsList []*api.AddLessons_LessonItem) (*api.AddLessons_Response, error)
}

type Endpoint struct {
	service IService
}

func NewEndpoint(service IService) *Endpoint {
	return &Endpoint{service: service}
}

func (e Endpoint) GetSchedule(ctx context.Context, request *api.GetSchedule_Request) (*api.GetSchedule_Response, error) {
	groupUuid := request.GetGroupUuid()
	if len(groupUuid) > 0 {
		return e.service.GetScheduleByGroupUuid(ctx, groupUuid)
	}

	groupName := request.GetGroupName()
	if len(groupName) > 0 {
		return e.service.GetScheduleByGroupName(ctx, groupName)
	}

	return nil, status.Error(codes.InvalidArgument, "expected group name or group uuid")
}

func (e Endpoint) AddLessons(ctx context.Context, request *api.AddLessons_Request) (*api.AddLessons_Response, error) {
	lessonsList := request.GetLessonsList()
	if lessonsList != nil {
		return e.service.AddLessons(ctx, lessonsList)
	}

	return nil, status.Error(codes.InvalidArgument, "expected lessons list")
}
