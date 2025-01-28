package utils

import (
	db "main/db/sqlc"
	"main/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user db.User) *pb.User {
	return &pb.User{
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
