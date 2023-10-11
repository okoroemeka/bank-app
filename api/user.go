package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password"  binding:"required,min=6"`
	FullName string `json:"full_name"  binding:"required"`
	Email    string `json:"email"  binding:"required,email"`
}

type createUserResponse struct {
	Username          string
	FullName          string
	Email             string
	IsEmailVerified   bool
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPass, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPass,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		IsEmailVerified:   user.IsEmailVerified,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	})
}
