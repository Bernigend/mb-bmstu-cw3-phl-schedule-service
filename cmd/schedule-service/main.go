package main

import (
	"fmt"
	group_service_api "github.com/Bernigend/mb-cw3-phll-group-service/pkg/group-service-api"
	"github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/endpoint"
	"github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/repository"
	"github.com/Bernigend/mb-cw3-phll-schedule-service/internal/app/service"
	api "github.com/Bernigend/mb-cw3-phll-schedule-service/pkg/schedule-service-api"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	grpcServerPort = 8992
)

const (
	dbHost = "localhost"
	dbPort = "5433"
	dbUser = "user"
	dbPass = "password"
	dbName = "db"
)

const (
	GroupServiceAddr = ""
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)
	db, err := repository.NewRepository(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	log.Println("db ok")

	err = db.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("migrations ok")

	grpcClientConfig, err := grpc.Dial(GroupServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	groupServiceClient := group_service_api.NewGroupServiceClient(grpcClientConfig)
	log.Println("group service ok")

	srv := service.NewService(db, groupServiceClient)
	log.Println("service ok")

	endpoints := endpoint.NewEndpoint(srv)
	log.Println("endpoint ok")

	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServerPort))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen ok")

	grpcServer := grpc.NewServer()
	api.RegisterScheduleServiceServer(grpcServer, endpoints)

	log.Println("starting grpc server...")
	err = grpcServer.Serve(listenPort)
	if err != nil {
		log.Fatal(err)
	}
}
