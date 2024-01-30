package gapi

import (
	"context"
	"errors"
	customerror "github.com/okoroemeka/simple_bank/custom-error"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/pb"
	"github.com/okoroemeka/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (res *pb.VerifyEmailResponse, err error) {
	violations := validateVerifyEmailRequest(req)

	if violations != nil {
		return nil, customerror.InvalidArgument(violations)
	}

	verifyEmailTxResp, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		ID:   req.GetId(),
		Code: req.GetCode(),
	})

	if err != nil {
		if errors.Is(err, customerror.ErrorNoRecordFound) {
			return nil, status.Errorf(codes.NotFound, "verify email record not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "verify email transaction failed: %s", err)
	}

	res = &pb.VerifyEmailResponse{
		User:        convertUser(verifyEmailTxResp.User),
		VerifyEmail: convertVerifyEmail(verifyEmailTxResp.VerifyEmail),
	}

	return res, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violation []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateVerifyEmailId(req.GetId()); err != nil {
		violation = append(violation, customerror.FieldViolation("id", err))
	}

	if err := val.ValidateString(req.GetCode(), 32, 32); err != nil {
		violation = append(violation, customerror.FieldViolation("secret_code", err))
	}

	return violation
}
