module github.com/Bernigend/mb-cw3-phll-schedule-service

go 1.15

require (
	github.com/Bernigend/mb-cw3-phll-group-service/pkg/group-service-api v1.0.0
	github.com/Bernigend/mb-cw3-phll-schedule-service/pkg/schedule-service-api v1.0.3
	github.com/satori/go.uuid v1.2.0
	google.golang.org/grpc v1.35.0
	gorm.io/driver/postgres v1.0.6
	gorm.io/gorm v1.20.10
)

replace (
	github.com/Bernigend/mb-cw3-phll-schedule-service/pkg/schedule-service-api => ./pkg/schedule-service-api
)
