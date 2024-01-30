package gapi

import (
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		IsEmailVerified:   user.IsEmailVerified,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt.Time),
		CreatedAt:         timestamppb.New(user.CreatedAt.Time),
	}

}

func convertVerifyEmail(verifyEmail db.VerifyEmail) *pb.VerifyEmail {
	return &pb.VerifyEmail{
		Id:         verifyEmail.ID,
		Username:   verifyEmail.Username,
		Email:      verifyEmail.Email,
		SecretCode: verifyEmail.SecretCode,
		IsUsed:     verifyEmail.IsUsed,
		ExpiredAt:  timestamppb.New(verifyEmail.ExpiredAt.Time),
		CreatedAt:  timestamppb.New(verifyEmail.CreatedAt.Time),
	}

}
