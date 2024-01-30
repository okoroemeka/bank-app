package gapi

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	customerror "github.com/okoroemeka/simple_bank/custom-error"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/pb"
	"github.com/okoroemeka/simple_bank/util"
	"github.com/okoroemeka/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx, []string{util.DepositorRole, util.BankerRole})
	if err != nil {
		return nil, customerror.UnAuthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)

	if authPayload.Role == util.DepositorRole && authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's information")
	}

	if violations != nil {
		return nil, customerror.InvalidArgument(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{String: req.GetEmail(), Valid: req.Email != nil},
	}

	if req.GetPassword() != "" {

		hashedPass, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot hash password: %s", err)
		}

		arg.HashedPassword = pgtype.Text{String: hashedPass, Valid: true}
		arg.PasswordChangedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	}

	user, err := server.store.UpdateUser(ctx, arg)

	if err != nil {
		if errors.Is(err, customerror.ErrorNoRecordFound) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, customerror.FieldViolation("username", err))
	}

	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, customerror.FieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, customerror.FieldViolation("email", err))
		}
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, customerror.FieldViolation("password", err))
		}
	}

	return
}
