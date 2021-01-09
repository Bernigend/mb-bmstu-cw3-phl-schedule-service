package main

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func main() {
	sts, ok := status.FromError(errors.New("test error"))
	log.Println(sts.Code(), sts.Details(), sts.Message(), ok)

	sts, ok = status.FromError(status.Error(codes.InvalidArgument, "test error invalid arg"))
	log.Println(sts.Code(), sts.Details(), sts.Message(), ok)
}
