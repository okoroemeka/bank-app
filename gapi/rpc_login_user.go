package gapi

import (
	"context"
	"database/sql"
	"errors"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/pb"
	"github.com/okoroemeka/simple_bank/util"
	"github.com/okoroemeka/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	violations := validateLoginUserRequest(req)

	if violations != nil {
		return nil, InvalidArgument(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	if err := util.ComparePassword(req.GetPassword(), user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password: %s", err)
	}

	token, payload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	refreshToken, RefreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	mtdt := server.extractMetadata(ctx)
	sessionArg := db.CreateSessionParams{
		ID:           RefreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    RefreshTokenPayload.ExpiresAt,
	}

	session, err := server.store.CreateSession(ctx, sessionArg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		AccessToken:           token,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(RefreshTokenPayload.ExpiresAt),
		AccessTokenExpiresAt:  timestamppb.New(payload.ExpiresAt),
		SessionId:             session.ID.String(),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, FieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, FieldViolation("password", err))
	}

	return
}
