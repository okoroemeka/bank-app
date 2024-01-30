package gapi

import (
	"context"
	"fmt"
	"github.com/okoroemeka/simple_bank/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationType   = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context, authorisedRoles []string) (*token.Payload, error) {
	metaData, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("missing metadata from request")
	}
	values := metaData.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHerder := values[0]
	fields := strings.Fields(authHerder)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authTyp := strings.ToLower(fields[0])
	if authTyp != authorizationType {
		return nil, fmt.Errorf("unsupported authorization type %s", authTyp)
	}

	tokenString := fields[1]
	payload, err := server.tokenMaker.VerifyToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %s", err)
	}

	if !hasPermission(authorisedRoles, payload.Role) {
		return nil, fmt.Errorf("unauthorized: %s", err)
	}

	return payload, nil
}

func hasPermission(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
